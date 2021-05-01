package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/justinas/nosurf"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/database"
	"github.com/mtcw99/disnews/sessions"
)

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
					nil, "ERROR: Invalid Login",
					nosurf.Token(r))
				return
			}

			if !login.Validate(password) {
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Invalid Login",
					nosurf.Token(r))
				return
			}

			uuid, err := sessions.GSession.NewSession(username)
			if err != nil {
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Session Error",
					nosurf.Token(r))
				return
			}
			sessionInfo := sessions.GSession.Keys[uuid]

			cookie := http.Cookie{
				Name:    "session_id",
				Value:   uuid,
				Expires: sessionInfo.Expire,
				Path:    "/",
			}
			http.SetCookie(w, &cookie)
		case "Create Account":
			login, err := core.LoginCreate(username, password)
			if err != nil {
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Cannot create account",
					nosurf.Token(r))
				return
			}

			err = database.DBase.Signup(login)
			if err != nil {
				fmt.Println(err)
				core.RenderTemplate(w, "login.html",
					nil, "ERROR: Cannot create account",
					nosurf.Token(r))
				return
			}
		}
		http.Redirect(w, r, "/", http.StatusFound)
	default:
		// User entering the login page
		core.RenderTemplate(w, "login.html", fetchSession(r), nil, nosurf.Token(r))
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Path[len("/profile/"):]
	profile, err := database.DBase.GetProfile(username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		core.RenderTemplate(w, "profile.html", fetchSession(r), profile, nosurf.Token(r))
	}
}

func ProfileUpdate(w http.ResponseWriter, r *http.Request) {
	session := fetchSession(r)
	if session == nil {
		fmt.Println("Failed to get session")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	r.ParseForm()
	username := r.FormValue("username")
	if username != session.Name {
		fmt.Println("Mis-match session")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	displayName := r.FormValue("display_name")
	link := r.FormValue("link")
	info := r.FormValue("info")
	profileid, err := database.DBase.GetProfileId(username)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	err = database.DBase.UpdateProfile(profileid, displayName, info, link)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/profile/"+username, http.StatusFound)
}
