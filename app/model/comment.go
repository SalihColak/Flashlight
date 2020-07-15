package model

import (
	couchdb "github.com/leesper/couchdb-golang"
)

type Comment struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	couchdb.Document
}
