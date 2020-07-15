package model

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	couchdb "github.com/leesper/couchdb-golang"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       string `json:"_id"`
	Rev      string `json:"_rev"`
	Type     string `json:"type"`
	Password string `json:"password"`
	Username string `json:"username"`
	couchdb.Document
}

// Db handle
var userDB *couchdb.Database

func init() {
	var err error
	userDB, err = couchdb.NewDatabase("http://localhost:5984/user")
	if err != nil {
		panic(err)
	}
}

func (user User) Exists() (exists bool) {
	query := `
	{
		"selector": {
			 "type": "User",
			 "username": "%s"
		}
	}`
	uExist, _ := userDB.QueryJSON(fmt.Sprintf(query, user.Username))

	if len(uExist) != 0 {
		return true
	}
	return false
}

func (user User) Add() (err error) {

	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	b64HashedPwd := base64.StdEncoding.EncodeToString(hashedPwd)

	user.Password = b64HashedPwd
	user.Type = "User"

	u, err := user2Map(user)

	delete(u, "_id")
	delete(u, "_rev")

	_, _, err = userDB.Save(u, nil)

	if err != nil {
		fmt.Printf("[Add] error: %s", err)
	}
	return err
}

// GetUserByUsername retrieve User by username
func GetUserByUsername(username string) (user User, err error) {
	user = User{}

	query := `
	{
		"selector": {
			 "type": "User",
			 "username": "%s"
		}
	}`
	u, err := userDB.QueryJSON(fmt.Sprintf(query, username))
	if err != nil {
		return user, err
	}
	if len(u) == 0 {
		return User{}, err
	}
	user, err = map2User(u[0])
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func user2Map(u User) (user map[string]interface{}, err error) {
	uJSON, err := json.Marshal(u)
	json.Unmarshal(uJSON, &user)

	return user, err
}

// Convert from map[string]interface{} to User struct as required by golang-couchdb methods
func map2User(user map[string]interface{}) (u User, err error) {
	uJSON, err := json.Marshal(user)
	json.Unmarshal(uJSON, &u)

	return u, err
}
