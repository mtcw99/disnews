// Package handlers provides HTTP handling functions for http.HandleFunc to use.
package handlers

import (
	"fmt"
	"log"
	"time"
	"net/http"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/database"
	"github.com/mtcw99/disnews/sessions"
)

func fetchSession(r *http.Request) *sessions.SessionInfo {
	hasSession := false
	cookie, err := r.Cookie("session_id")
	var sessionInfo sessions.SessionInfo
	if err == nil {
		uuid := cookie.Value
		hasSession = sessions.GSession.ValidateSession(uuid)
		if hasSession {
			var ok bool
			sessionInfo, ok = sessions.GSession.Get(uuid)
			if !ok {
				hasSession = false
			}
		}
	}

	if hasSession {
		return &sessionInfo
	} else {
		return nil
	}
}

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

	core.RenderTemplate(w, "index.html", fetchSession(r), posts)
}

// New post handler
func NewPost(w http.ResponseWriter, r *http.Request) {
	session := fetchSession(r)
	if session == nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		core.RenderTemplate(w, "new.html", session, nil)
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

	// Check Link
	if link[:len("https://")] != "https://" && link[:len("http://")] != "http://" {
		link = "https://" + link
	}

	id, err := database.DBase.SubmitPost(core.Post{
		User:	 session.Name,
		Title:   title,
		Link:    link,
		Comment: comment})
	if err != nil {
		// TODO: Error message
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		core.RenderTemplate(w, "submitted.html", session, id)
	}
}

// View the requested post
func PostView(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/post/"):]
	post, err := database.DBase.GetPost(id)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		core.RenderTemplate(w, "post.html", fetchSession(r), post)
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
					nil, "ERROR: Invalid Login")
				return
			}

			if !login.Validate(password) {
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Invalid Login")
				return
			}

			uuid, err := sessions.GSession.NewSession(username)
			if err != nil {
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Session Error")
				return
			}
			sessionInfo := sessions.GSession.Keys[uuid]

			cookie := http.Cookie{
				Name: "session_id",
				Value: uuid,
				Expires: sessionInfo.Expire,
				Path: "/",
			}
			http.SetCookie(w, &cookie)
		case "Create Account":
			login, err := core.LoginCreate(username, password)
			if err != nil {
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Cannot create account")
				return
			}

			err = database.DBase.Signup(login)
			if err != nil {
				fmt.Println(err)
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Cannot create account")
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		// User entering the login page
		core.RenderTemplate(w, "login.html", fetchSession(r), nil)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name: "session_id",
		Value: "",
		Expires: time.Unix(0, 0),
		Path: "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

// Error handler gets status and file rendered for the status
func errorHandler(w http.ResponseWriter, r *http.Request, status int, errorTmplFile string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		core.RenderTemplate(w, errorTmplFile, fetchSession(r), nil)
	}
}
