package main

import (
	"fmt"
	"net/http"
	"time"
	"web-ss20/flashlight/app/controller"
)

func main() {

	http.HandleFunc("/", controller.Index)
	http.HandleFunc("/authenticate", controller.AuthenticateUser)
	http.HandleFunc("/adduser", controller.AddUser)
	http.HandleFunc("/bilder", controller.Auth(controller.Bilder))
	http.HandleFunc("/login", controller.Login)
	http.HandleFunc("/register", controller.Register)
	http.HandleFunc("/upload", controller.Auth(controller.Upload))
	http.HandleFunc("/upload-posting", controller.Auth(controller.UploadPosting))
	http.HandleFunc("/delete-posting", controller.Auth(controller.DeletePosting))
	http.HandleFunc("/add-comment", controller.Auth(controller.AddComment))
	http.HandleFunc("/add-like", controller.Auth(controller.AddLike))
	http.HandleFunc("/del-like", controller.Auth(controller.DeleteLike))
	http.HandleFunc("/logout", controller.Logout)

	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir("."))))

	http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("./static/css"))))
	http.Handle("/fonts/", http.StripPrefix("/fonts", http.FileServer(http.Dir("./static/fonts"))))
	http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("./static/js"))))
	http.Handle("/images/", http.StripPrefix("/images", http.FileServer(http.Dir("./static/images"))))
	http.Handle("/postings/", http.StripPrefix("/postings", http.FileServer(http.Dir("./postings"))))

	fmt.Println("Starting Server")
	t := time.Now()
	fmt.Printf("%02d.%02d.%d - %02d:%02d Uhr\n", t.Day(), t.Month(), t.Year(), t.Hour(), t.Minute())
	server := http.Server{
		Addr: ":80",
	}
	server.ListenAndServe()
}
