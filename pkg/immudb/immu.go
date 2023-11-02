package immudb

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/DenisPalnitsky/immu-svn/pkg/data"
	resty "github.com/go-resty/resty/v2"
)

type ImmudbClient struct {
	baseUrl    string
	restClient *resty.Client
}

var RepositoryNotFound = fmt.Errorf("collection not found")

func NewImmudbClient(apiKey string) *ImmudbClient {
	var client = resty.New().SetDoNotParseResponse(true)
	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"X-API-Key":    apiKey,
	})

	return &ImmudbClient{
		baseUrl:    "https://vault.immudb.io/ics/api/v1/ledger/default/collection/",
		restClient: client,
	}
}

func (c *ImmudbClient) CreateRepo(collectionName string) error {
	url := c.baseUrl + collectionName

	payload := `{
		"fields": [
			{
				"name": "path",
				"type": "STRING"
			},
			{
				"name": "content",
				"type": "STRING"
			}
		],
		"indexes": [
			{
				"fields": [
					"path"
				],
				"isUnique": true
			}
		]
	}`

	resp, err := c.restClient.R().SetBody(payload).Put(url)
	if err != nil {
		return fmt.Errorf("error creating collection %w", err)
	}
	if resp.StatusCode() == http.StatusConflict {
		return fmt.Errorf("collection already exists")
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("error creating collection", "body", string(resp.Body()))
		return fmt.Errorf("error creating collection")
	}

	return nil
}

// AddOrUpdateFiles adds or updates files in the collection. The map contains path as key and content as value
func (c *ImmudbClient) AddOrUpdateFiles(collectionName string, files map[string]string) (added int, updated int, err error) {
	existingDocs, err := c.getDocumentsList(collectionName)
	if err != nil {
		return 0, 0, fmt.Errorf("error getting documents list. Error: %w", err)
	}

	// update those that exist removing files from the list
	for _, doc := range existingDocs {
		if cont, ok := files[doc.Path]; ok {
			err = c.updateDocument(collectionName, doc.Path, cont)
			if err != nil {
				slog.Error("error updating document", "path", doc.Path, "error", err)
				return added, updated, fmt.Errorf("error updating document")
			}
			updated++
			delete(files, doc.Path)
		}
	}
	added = len(files)
	if added > 0 {
		err = c.addDocuments(collectionName, files)
		if err != nil {
			slog.Error("error adding documents", "error", err)
			return added, updated, fmt.Errorf("error adding documents")
		}
	}

	return added, updated, nil
}

func (c *ImmudbClient) addDocuments(collectionName string, files map[string]string) error {

	// TODO: in case we have big files or too many of them, we need to split them into chunks
	url := c.baseUrl + collectionName + "/documents"

	docs := make([]Document, 0)
	for path, content := range files {
		docs = append(docs, Document{Path: path, Content: content})
	}

	resp, err := c.restClient.R().SetBody(CreateDocumentsRequest{Documents: docs}).Put(url)

	if err != nil {
		slog.Error("error adding document", "error", err)
		return fmt.Errorf("error adding document")
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("error adding document", "body", string(resp.Body()))
		return fmt.Errorf("error adding document")
	}
	return nil
}

func (c *ImmudbClient) updateDocument(collectionName string, path string, content string) error {
	url := c.baseUrl + collectionName + "/document"

	req := UpdateRequest{
		Document: Document{
			Path:    path,
			Content: content,
		},
		Query: Query{
			Expressions: []Expression{
				{
					FieldComparisons: []FieldComparison{
						{
							Field:    "path",
							Operator: "EQ",
							Value:    path,
						},
					},
				},
			}}}

	marshReq, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshalling request %w", err)
	}
	slog.Debug("update request", "request", string(marshReq))

	resp, err := c.restClient.R().SetBody(marshReq).Post(url)
	if err != nil {
		return fmt.Errorf("error updating document. Error:t %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("error updating document", "body", string(resp.Body()), "status", resp.StatusCode())
		return fmt.Errorf("error updating document")
	}

	return nil
}

func (c *ImmudbClient) Diff(collectionName string, path string) ([]data.DiffLogItem, error) {

	// it would be better to get ID by a path filed but for the sake of simplicity I'll get all documents
	documents, err := c.getDocumentsList(collectionName)
	if err != nil {
		slog.Error("error getting documents list", "error", err)
		return nil, fmt.Errorf("error getting documents list")
	}

	d, ok := documents[path]
	if !ok {
		return nil, fmt.Errorf("document not found")
	}

	url, err := url.JoinPath(c.baseUrl, collectionName, "/document/", d.Id, "/audit")
	if err != nil {
		slog.Error("error creating url", "error", err)
		return nil, fmt.Errorf("error creating url")
	}

	resp, err := c.restClient.R().SetBody(`{ "page": 1,	"perPage": 10, "desc": true }`).Put(url)
	if err != nil {
		slog.Error("error getting diff", "error", err)
		return nil, fmt.Errorf("error getting diff")
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("error getting diff", "body", string(resp.Body()), "status", resp.StatusCode())
		return nil, fmt.Errorf("error getting diff")
	}

	decoder := json.NewDecoder(resp.RawBody())
	var result DiffResponse

	err = decoder.Decode(&result)
	if err != nil {
		slog.Error("error decoding diff", "error", err)
		return nil, fmt.Errorf("error decoding diff")
	}

	res := []data.DiffLogItem{}
	for _, diff := range result.Revision {
		res = append(res, data.DiffLogItem{
			Revision:  diff.Revision,
			Content:   diff.Document.Content,
			Timestamp: time.Unix(diff.Document.VaultMD.Ts, 0),
		})
	}
	return res, nil
}

func (c *ImmudbClient) getDocumentsList(collectionName string) (map[string]Document, error) {
	// TODO: implement fetching all documents. Right now we support max 100 files
	url := c.baseUrl + collectionName + "/documents/search"
	payload := `{   
		"query": {
		  "orderBy": [
			{
			  "field": "path",
			  "desc": true
			}
		  ],
		  "limit": 0
		},
		"page": 1,
		"perPage": 100
	  }`

	resp, err := c.restClient.R().SetBody(payload).Post(url)

	if err != nil {
		return nil, fmt.Errorf("error getting documents list. Error: %w", err)
	}

	if resp.StatusCode() != http.StatusOK {
		slog.Error("error getting documents list", "body", string(resp.Body()))
		return nil, fmt.Errorf("error getting documents list")
	}

	j := &SearchResponse{}
	err = json.NewDecoder(resp.RawBody()).Decode(j)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling documents list. Error: %w", err)
	}

	res := make(map[string]Document, 0)
	for _, r := range j.Revisions {
		res[r.Document.Path] = r.Document
	}

	return res, nil
}

func (c *ImmudbClient) GetCollectionInfo(collectionName string) error {
	url := c.baseUrl + collectionName
	resp, err := c.restClient.R().Get(url)
	if err != nil {
		return fmt.Errorf("error getting collection info %w", err)
	}
	if resp.StatusCode() == http.StatusNotFound {
		return RepositoryNotFound
	}
	if resp.StatusCode() != http.StatusOK {
		slog.Error("error getting collection info", "body", string(resp.Body()), "status", resp.StatusCode())
		return fmt.Errorf("error getting collection info")
	}

	return nil
}
