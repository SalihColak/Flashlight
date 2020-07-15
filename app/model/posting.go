package model

import (
	"encoding/json"
	"fmt"

	couchdb "github.com/leesper/couchdb-golang"
)

type Posting struct {
	Id            string    `json:"_id"`
	Rev           string    `json:"_rev"`
	Type          string    `json:"type"`
	Description   string    `json:"description"`
	UserID        string    `json:"userid"`
	Username      string    `json:"username"`
	Path          string    `json:"path"`
	Date          string    `json:"date"`
	Comments      []Comment `json:"comments"`
	CommentNumber int       `json:"commentnumber"`
	Likes         []Like    `json:"likes"`
	LikeNumber    int       `json:"likenumber"`
	UserLiked     bool      `json:"userliked"` //zur Anzeige des Likestatus vom aktuellen User
	Zero          int       `json:"zero"`      //zum Vergleich bei der Like Anzeige im Tmpl für den anzuzeigenden String
	One           int       `json:"one"`       //zum Vergleich bei der Like Anzeige im Tmpl für den anzuzeigenden String
	couchdb.Document
}

// Db handle
var postingDB *couchdb.Database

func init() {
	var err error
	postingDB, err = couchdb.NewDatabase("http://localhost:5984/posting")
	if err != nil {
		panic(err)
	}
	//postingDB.PutIndex([]string{"args(date)"}, "", "date")
}

//Add Posting
func (posting Posting) Add() (err error) {

	p, err := posting2Map(posting)

	delete(p, "_id")
	delete(p, "_rev")

	_, _, err = postingDB.Save(p, nil)

	if err != nil {
		fmt.Printf("[Add] error: %s", err)
	}

	return err
}

func GetPosting(id string) (Posting, error) {
	p, err := postingDB.Get(id, nil)
	if err != nil {
		return Posting{}, err
	}

	posting, _ := map2Posting(p)

	/*posting := Posting{
		Id:          p["_id"].(string),
		Rev:         p["_rev"].(string),
		Type:        p["type"].(string),
		Date:        p["date"].(string),
		Description: p["description"].(string),
		Path:        p["path"].(string),
		UserID:      p["userid"].(string),
		Username:    p["username"].(string),
	}

	var comments []Comment

	jsonStr, err := json.Marshal(p["comments"])

	json.Unmarshal(jsonStr, &comments)

	println(comments)
	posting.Comments = comments*/

	return posting, nil
}

func GetAllPostings() ([]map[string]interface{}, error) {
	query := `
	{
		"selector": {
		   "type": "Posting",
		   "date": {
			  "$gt": null
		   }
		},
		"sort": [
		   {
			  "date": "desc"
		   }
		]
	}`
	allPostings, err := postingDB.QueryJSON(fmt.Sprintf(query))
	if err != nil {

		return nil, err
	} else {

		return allPostings, nil
	}
}

func GetAllPostingsByUserID(userid string) ([]map[string]interface{}, error) {
	query := `
	{
		"selector": {
		   "type": "Posting",
		   "userid": "%s",
		   "date": {
			  "$gt": null
		   }
		},
		"sort": [
		   {
			  "date": "desc"
		   }
		]
	}`
	allPostingsOfUser, err := postingDB.QueryJSON(fmt.Sprintf(query, userid))
	if err != nil {
		return nil, err
	} else {
		return allPostingsOfUser, nil
	}
}

func (p Posting) Delete() error {
	err := postingDB.Delete(p.Id)
	return err
}

func (posting Posting) UpdatePosting() error {
	p, err := posting2Map(posting)

	if err != nil {
		fmt.Printf("[Posting2Map] error: %s", err)
	}

	err = postingDB.Set(posting.Id, p)
	if err != nil {
		fmt.Printf("[UpdatePosting] error: %s", err)
	}
	return err
}

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func posting2Map(p Posting) (posting map[string]interface{}, err error) {
	uJSON, err := json.Marshal(p)
	json.Unmarshal(uJSON, &posting)

	return posting, err
}

// Convert from map[string]interface{} to User struct as required by golang-couchdb methods
func map2Posting(posting map[string]interface{}) (p Posting, err error) {
	uJSON, err := json.Marshal(posting)
	json.Unmarshal(uJSON, &p)

	return p, err
}
