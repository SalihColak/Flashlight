package model

import (
	couchdb "github.com/leesper/couchdb-golang"
)

type Like struct {
	Username string `json:"username"`
	couchdb.Document
}
