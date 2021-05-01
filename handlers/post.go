package handlers

import (
	"fmt"
	"net/http"

	"github.com/justinas/nosurf"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/database"
)

// New post handler
func NewPost(w http.ResponseWriter, r *http.Request) {
	session := fetchSession(r)
	if session == nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		core.RenderTemplate(w, "new.html", session, nil, nosurf.Token(r))
	}
}

// Post submission handler | Requires a POST request
func SubmitPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	link := r.FormValue("link")
	comment := r.FormValue("comment")

	session := fetchSession(r)
	if session == nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	link = core.LinkFix(link)

	id, err := database.DBase.SubmitPost(core.Post{
		User:    session.Name,
		Title:   title,
		Link:    link,
		Comment: comment})
	if err != nil {
		// TODO: Error message
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	core.RenderTemplate(w, "submitted.html", session, id, nosurf.Token(r))
}

// View the requested post
func PostView(w http.ResponseWriter, r *http.Request) {
	session := fetchSession(r)
	id := r.URL.Path[len("/post/"):]
	postComments, err := database.DBase.GetPostAndComments(id)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		if session != nil {
			postComments.Post.ParseVoted, err = database.DBase.GetVote(
				session.Name, postComments.Post.Id)
		}
		core.RenderTemplate(w, "post.html", session, postComments, nosurf.Token(r))
	}
}
