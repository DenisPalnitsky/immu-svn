package data

import "fmt"

var RepositoryNotFound = fmt.Errorf("collection not found")
var CollectionAlreadyExist = fmt.Errorf("collection already exists")
