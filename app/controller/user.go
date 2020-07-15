package controller

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"web-ss20/flashlight/app/model"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

var store *sessions.CookieStore

func init() {
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key := make([]byte, 32)
	rand.Read(key)
	store = sessions.NewCookieStore(key)
}

// Register controller
func Register(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "register.tmpl", nil)
}

// AddUser controller
func AddUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	passwordRe := r.FormValue("passwordRe")

	if password != passwordRe {
		data := struct {
			AlertText string
			HasAlert  bool
		}{
			"Passwörter stimmen nicht überein!",
			true,
		}
		tmpl.ExecuteTemplate(w, "register.tmpl", data)
	} else {

		user := model.User{}
		user.Username = username

		if user.Exists() {
			data := struct {
				AlertText string
				HasAlert  bool
			}{
				"Benutzername existiert bereits!",
				true,
			}
			tmpl.ExecuteTemplate(w, "register.tmpl", data)
		} else {
			user.Password = password

			user.Add()

			data := struct {
				Username string
				Success  bool
			}{
				username,
				true,
			}

			tmpl.ExecuteTemplate(w, "login.tmpl", data)
		}
	}

}

// Login controller
func Login(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "login.tmpl", nil)
}

// AuthenticateUser controller
func AuthenticateUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Authentication
	user, _ := model.GetUserByUsername(username)

	if (model.User{}) == user {
		mydata := struct {
			Error    string
			HasError bool
			Success  bool
			Username string
		}{
			"Benutzername und Passwort stimmen nicht überein!",
			true,
			false,
			"",
		}
		tmpl.ExecuteTemplate(w, "login.tmpl", mydata)
	} else {
		// decode base64 String to []byte
		passwordDB, _ := base64.StdEncoding.DecodeString(user.Password)
		err := bcrypt.CompareHashAndPassword(passwordDB, []byte(password))

		if err == nil {
			session, _ := store.Get(r, "session")

			// Set user as authenticated
			session.Values["authenticated"] = true
			session.Values["username"] = username
			session.Values["userid"] = user.Id
			session.Save(r, w)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			mydata := struct {
				Error    string
				HasError bool
				Success  bool
				Username string
			}{
				"Benutzername und Passwort stimmen nicht überein!",
				true,
				false,
				"",
			}
			tmpl.ExecuteTemplate(w, "login.tmpl", mydata)
		}
	}

}

// Logout controller
func Logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	session.Values["authenticated"] = false
	session.Values["username"] = ""
	session.Values["userid"] = ""
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Authentication handler
func Auth(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		// Check if user is authenticated
		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Redirect(w, r, "/login", http.StatusFound)
		} else {
			h(w, r)
		}
	}
}
