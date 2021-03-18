package handlers

import (
	"github.com/mtcw99/disnews/core"
	"net/http"
)

func Root(w http.ResponseWriter, r *http.Request) {
	core.RenderTemplate(w, "index.html")
}

func NewPost(w http.ResponseWriter, r *http.Request) {

}

func PostView(w http.ResponseWriter, r *http.Request) {

}
