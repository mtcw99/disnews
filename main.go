package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/database"
	"github.com/mtcw99/disnews/handlers"
)

func main() {
	core.Info.PathTemplates = "./templates/"
	err := database.DBase.New("./test.db")
	database.DBase.Setup()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer database.DBase.Close()

	fmt.Println("Serving at: http://localhost:8080/")

	http.HandleFunc("/", handlers.Root)
	http.HandleFunc("/new/", handlers.NewPost)
	http.HandleFunc("/submit/", handlers.SubmitPost)
	http.HandleFunc("/post/", handlers.PostView)
	http.HandleFunc("/css/", handlers.Css)
	http.HandleFunc("/js/", handlers.Js)
	http.HandleFunc("/login/", handlers.Login)
	http.HandleFunc("/logout/", handlers.Logout)
	http.HandleFunc("/profile/", handlers.Profile)
	http.HandleFunc("/profile_update/", handlers.ProfileUpdate)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
