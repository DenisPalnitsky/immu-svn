package immudb

import (
	"encoding/json"
	"fmt"
	"testing"
)

// func TestCreateCollection(t *testing.T) {
// 	c := NewImmudbClient("default.mTVqmj_u1Ov4Js0b4mRYyg.jgB_i0o28E0ZKfEBeg3LWuVzgFP0uxd4J1dXjmOic4nMMBU0")
// 	c.CreateCollection("unit_test_1")

// }

func TestParsing(t *testing.T) {
	var resp = `{
		"page": 1,
		"perPage": 100,
		"revisions": [
		  {
			"document": {
			  "_id": "6542b39b0000000000000037350e2a7a",
			  "_vault_md": {
				"creator": "a:ea113865-ae22-4baa-94b0-7bce29cb8a85",
				"ts": 1698870171
			  },
			  "content": "test",
			  "path": "main1.go"
			},
			"revision": "",
			"transactionId": ""
		  },
		  {
			"document": {
			  "_id": "6542b39b0000000000000037350e2a7b",
			  "_vault_md": {
				"creator": "a:ea113865-ae22-4baa-94b0-7bce29cb8a85",
				"ts": 1698870171
			  },
			  "content": "test",
			  "path": "main.go"
			},
			"revision": "",
			"transactionId": ""
		  },
		  {
			"document": {
			  "_id": "6542aa500000000000000030350e2a79",
			  "_vault_md": {
				"creator": "a:ea113865-ae22-4baa-94b0-7bce29cb8a85",
				"ts": 1698867792
			  },
			  "content": "Hello world!!!!!!!!!",
			  "path": "/main2.go"
			},
			"revision": "",
			"transactionId": ""
		  },
		  {
			"document": {
			  "_id": "6542aa39000000000000002f350e2a77",
			  "_vault_md": {
				"creator": "a:ea113865-ae22-4baa-94b0-7bce29cb8a85",
				"ts": 1698870085
			  },
			  "content": "Helloooo Woooorld!!"
			},
			"revision": "",
			"transactionId": ""
		  }
		],
		"searchId": ""
	  }`

	j := &SearchResponse{}
	json.Unmarshal([]byte(resp), j)

	fmt.Printf("%+v", j.Revisions[0])

}
