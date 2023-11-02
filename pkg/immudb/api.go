package immudb

// I could have used autogenerated client from swagger, but as we use only part of AI API, I decided to write it manually
type CreateDocumentsRequest struct {
	Documents []Document `json:"documents"`
}

type UpdateDocumentResponse struct {
}

type SearchResponse struct {
	Revisions []struct {
		Document Document `json:"document"`
	} `json:"revisions"`
}

type Document struct {
	Id      string `json:"_id,omitempty"`
	Content string `json:"content"`
	Path    string `json:"path"`
}

type FieldComparison struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Expression struct {
	FieldComparisons []FieldComparison `json:"fieldComparisons"`
}

type Query struct {
	Expressions []Expression `json:"expressions"`
}

type UpdateRequest struct {
	Document Document `json:"document"`
	Query    Query    `json:"query"`
}

// revisions
type DiffResponse struct {
	Revision []Revision `json:"revisions"`
}

type VaultMD struct {
	Creator string `json:"creator"`
	Ts      int64  `json:"ts"`
}

type RevisionDocument struct {
	Document
	VaultMD VaultMD `json:"_vault_md"`
}

type Revision struct {
	Document      RevisionDocument `json:"document"`
	Revision      string           `json:"revision"`
	TransactionID string           `json:"transactionId"`
}
