package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/handlers"
)

func main() {
	core.Info.Path = "./templates/"
	fmt.Println("Serving at: http://localhost:8080/")
	http.HandleFunc("/", handlers.Root)
	http.HandleFunc("/new/", handlers.NewPost)
	http.HandleFunc("/post/", handlers.PostView)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
