package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/nosurf"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/database"
	"github.com/mtcw99/disnews/handlers"
)

func main() {
	core.Info.PathTemplates = "./templates/"
	err := database.DBase.New("./db/test.db")
	database.DBase.Setup()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DBase.Close()

	fmt.Println("Serving at: http://localhost:8080/")

	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.Root)
	mux.HandleFunc("/new/", handlers.NewPost)
	mux.HandleFunc("/submit/", handlers.SubmitPost)
	mux.HandleFunc("/post/", handlers.PostView)
	mux.HandleFunc("/css/", handlers.Css)
	mux.HandleFunc("/js/", handlers.Js)
	mux.HandleFunc("/login/", handlers.Login)
	mux.HandleFunc("/logout/", handlers.Logout)
	mux.HandleFunc("/profile/", handlers.Profile)
	mux.HandleFunc("/profile_update/", handlers.ProfileUpdate)
	mux.HandleFunc("/comment/", handlers.Comment)
	mux.HandleFunc("/vote_up/", handlers.VotePostUp)
	mux.HandleFunc("/vote_down/", handlers.VotePostDown)

	http.ListenAndServe(":8080", nosurf.New(mux))
}
