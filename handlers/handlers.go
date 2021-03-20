// Package handlers provides HTTP handling functions for http.HandleFunc to use.
package handlers

import (
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/database"
)

// '/' Root/main page handler, also catches 404 errors
func Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorHandler(w, r, http.StatusNotFound, "404.html")
		return
	}

	posts, err := database.DBase.GetNewestPosts()
	if err != nil {
		log.Fatal(err)
		return
	}
	core.RenderTemplate(w, "index.html", posts)
}

// New post handler
func NewPost(w http.ResponseWriter, r *http.Request) {
	core.RenderTemplate(w, "new.html", nil)
}

// Post submission handler | Requires a POST request
func SubmitPost(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("title")
	link := r.FormValue("link")
	comment := r.FormValue("comment")

	// Check Link
	switch {
	case link[:len("https://")] != "https://",
		link[:len("http://")] != "http://":

		link = "https://" + link
	}

	id, err := database.DBase.SubmitPost(core.Post{
		Title:   title,
		Link:    link,
		Comment: comment})
	if err != nil {
		core.RenderTemplate(w, "index.html", nil)
	} else {
		core.RenderTemplate(w, "submitted.html", id)
	}
}

// View the requested post
func PostView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/post/"):]
	post, err := database.DBase.GetPost(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		core.RenderTemplate(w, "post.html", post)
	}
}

// CSS Handler serves static CSS Files
func Css(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, core.Info.PathStatic+"css/"+r.URL.Path[len("/css/"):])
}

// JS Handler serves static JavaScript Files
func Js(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, core.Info.PathStatic+"js/"+r.URL.Path[len("/js/"):])
}

// Login Handler | Handles signups and logins requests
func Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	actionType := r.FormValue("action")
	switch actionType {
	case "Login", "Create Account":
		username := r.FormValue("username")
		password := r.FormValue("password")
		switch actionType {
		case "Login":
			login, err := database.DBase.Login(username)
			if err != nil {
				fmt.Println(err)
				core.RenderTemplate(w, "login.html",
					"ERROR: Invalid Login")
				return
			}

			if !login.Validate(password) {
				core.RenderTemplate(w, "login.html",
					"ERROR: Invalid Login")
				return
			}

			// TODO: Set session cookie
			expire := time.Now().Add(24 * time.Hour)
			cookie := http.Cookie{
				Name: "session_id",
				Value: "testing",
				Expires: expire,
			}
			http.SetCookie(w, &cookie)
		case "Create Account":
			login, err := core.LoginCreate(username, password)
			if err != nil {
				core.RenderTemplate(w, "login.html",
					"ERROR: Cannot create account")
				return
			}

			err = database.DBase.Signup(login)
			if err != nil {
				fmt.Println(err)
				core.RenderTemplate(w, "login.html",
					"ERROR: Cannot create account")
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		// User entering the login page
		core.RenderTemplate(w, "login.html", nil)
	}
}

// Error handler gets status and file rendered for the status
func errorHandler(w http.ResponseWriter, r *http.Request, status int, errorTmplFile string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		core.RenderTemplate(w, errorTmplFile, nil)
	}
}
