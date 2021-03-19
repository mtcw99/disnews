// Package handlers provides HTTP handling functions for http.HandleFunc to use.
package handlers

import (
	"net/http"

	"github.com/mtcw99/disnews/core"
)

// '/' Root/main page handler, also catches 404 errors
func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, "404.html")
		return
	}

	core.RenderTemplate(w, "index.html")
}

// New post handler
func NewPost(w http.ResponseWriter, r *http.Request) {
	core.RenderTemplate(w, "new.html")
}

// Post submission handler | Requires a POST request
func SubmitPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	link := r.FormValue("link")
	comment := r.FormValue("comment")
	core.DBase.SubmitPost(core.Post{
		Title:   title,
		Link:    link,
		Comment: comment})
}

// View the requested post
func PostView(w http.ResponseWriter, r *http.Request) {
	core.RenderTemplate(w, "view.html")
}

// CSS Handler serves static CSS Files
func Css(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, core.Info.PathStatic+"css/"+r.URL.Path[len("/css/"):])
}

// JS Handler serves static JavaScript Files
func Js(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, core.Info.PathStatic+"js/"+r.URL.Path[len("/js/"):])
}

// Error handler gets status and file rendered for the status
func errorHandler(w http.ResponseWriter, r *http.Request, status int, errorTmplFile string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		core.RenderTemplate(w, errorTmplFile)
	}
}
