// Package handlers provides HTTP handling functions for http.HandleFunc to use.
package handlers

import (
	"log"
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

	core.RenderTemplate(w, "index.html", fetchSession(r), posts, "")
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
		core.RenderTemplate(w, errorTmplFile, fetchSession(r), nil, "")
	}
}
