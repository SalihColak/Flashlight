package controller

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"web-ss20/flashlight/app/model"
)

var tmpl *template.Template

func init() {
	rand.Seed(time.Now().UnixNano())
	tmpl = template.Must(template.ParseGlob("template/*.tmpl"))
}

func Index(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	postings, _ := model.GetAllPostings()

	if session.IsNew {
		dataNoLog := struct {
			Authenticated bool
			Postings      *[]map[string]interface{}
		}{
			false,
			&postings,
		}
		tmpl.ExecuteTemplate(w, "index.tmpl", dataNoLog)
	} else {

		authenticated := session.Values["authenticated"].(bool)
		username := session.Values["username"].(string)
		/*data := struct {
			Authenticated bool
			Username      string
			Postings      *[]map[string]interface{}
		}{
			authenticated,
			username,
			&postings,
		}*/

		var postingsStruct []model.Posting
		var postingsStructNew []model.Posting

		for _, element := range postings {
			posting, _ := map2Posting(element)
			postingsStruct = append(postingsStruct, posting)

		}

		for _, post := range postingsStruct {
			likes := post.Likes
			for _, like := range likes {
				if like.Username == username {
					post.UserLiked = true
				}
			}
			postingsStructNew = append(postingsStructNew, post)
		}

		//println(postingsStructNew[0].UserLiked)

		var postingmap []map[string]interface{}

		for _, posting := range postingsStructNew {
			postIf, _ := posting2Map(posting)
			postingmap = append(postingmap, postIf)
		}

		data2 := struct {
			Authenticated bool
			Username      string
			Postings      *[]map[string]interface{}
		}{
			authenticated,
			username,
			&postingmap,
		}

		tmpl.ExecuteTemplate(w, "index.tmpl", data2)
	}

}

func Upload(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)

	data := struct {
		Username string
	}{
		username,
	}

	tmpl.ExecuteTemplate(w, "upload.tmpl", data)
}

func UploadPosting(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "session")

	r.ParseMultipartForm(32 << 20)
	// "img" ist das Attribut name des Html-Input-Tags
	file, metadata, err := r.FormFile("img")
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	defer file.Close()

	str := strings.SplitAfter(metadata.Filename, ".")

	name_of_image := strings.Trim(str[0], ".") + RandStringBytes(8) + "." + str[1]
	f, err := os.OpenFile("./postings/"+name_of_image, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	posting := model.Posting{}

	posting.Type = "Posting"
	posting.UserID = session.Values["userid"].(string)
	posting.Username = session.Values["username"].(string)
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		fmt.Println(err)
	}
	t := time.Now()
	date := fmt.Sprintf("%02d.%02d.%d - %02d:%02d Uhr", t.In(location).Day(), t.In(location).Month(), t.In(location).Year(), t.In(location).Hour(), t.In(location).Minute())
	posting.Date = date
	posting.Description = r.FormValue("description")
	posting.Path = "/postings/" + name_of_image
	posting.Comments = []model.Comment{}
	posting.LikeNumber = 0
	posting.CommentNumber = 0
	posting.Zero = 0
	posting.One = 1

	posting.Add()

	http.Redirect(w, r, "/", http.StatusFound)
}

func DeletePosting(w http.ResponseWriter, r *http.Request) {
	_id := r.FormValue("_id")
	posting, _ := model.GetPosting(_id)
	posting.Delete()
	err := os.Remove("." + posting.Path)
	if err != nil {
		fmt.Println(err)
	}
	http.Redirect(w, r, "bilder", http.StatusFound)
}

func Bilder(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)
	userid := session.Values["userid"].(string)

	postings, _ := model.GetAllPostingsByUserID(userid)

	data := struct {
		Username string
		Postings *[]map[string]interface{}
	}{
		username,
		&postings,
	}

	tmpl.ExecuteTemplate(w, "bilder.tmpl", data)
}

func AddComment(w http.ResponseWriter, r *http.Request) {
	postingid := r.FormValue("postingid")

	posting, _ := model.GetPosting(postingid)

	text := r.FormValue("newcomment")
	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)

	comment := model.Comment{}
	comment.Text = text
	comment.Username = username

	Comments := posting.Comments

	Comments = append(Comments, comment)

	posting.Comments = Comments
	posting.CommentNumber = posting.CommentNumber + 1

	err := posting.UpdatePosting()
	if err != nil {
		println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func AddLike(w http.ResponseWriter, r *http.Request) {
	postingid := r.FormValue("postingid")

	posting, _ := model.GetPosting(postingid)

	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)

	like := model.Like{}
	like.Username = username

	Likes := posting.Likes

	Likes = append(Likes, like)

	posting.Likes = Likes

	posting.LikeNumber = posting.LikeNumber + 1

	err := posting.UpdatePosting()
	if err != nil {
		println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func DeleteLike(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("postingid")
	posting, _ := model.GetPosting(id)

	session, _ := store.Get(r, "session")
	username := session.Values["username"].(string)

	Likes := posting.Likes
	var newLikes []model.Like

	for _, like := range Likes {
		if like.Username != username {
			newLikes = append(newLikes, like)
		}
	}
	posting.Likes = newLikes
	posting.LikeNumber = posting.LikeNumber - 1

	err := posting.UpdatePosting()
	if err != nil {
		println(err)
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

//Creates Random String for Filenames with length n using the letterBytes
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Convert from User struct to map[string]interface{} as required by golang-couchdb methods
func posting2Map(p model.Posting) (posting map[string]interface{}, err error) {
	uJSON, err := json.Marshal(p)
	json.Unmarshal(uJSON, &posting)

	return posting, err
}

// Convert from map[string]interface{} to User struct as required by golang-couchdb methods
func map2Posting(posting map[string]interface{}) (p model.Posting, err error) {
	uJSON, err := json.Marshal(posting)
	json.Unmarshal(uJSON, &p)

	return p, err
}
