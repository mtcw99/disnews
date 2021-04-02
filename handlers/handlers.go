// Package handlers provides HTTP handling functions for http.HandleFunc to use.
package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/justinas/nosurf"

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
	id := r.URL.Path[len("/post/"):]
	postComments, err := database.DBase.GetPostAndComments(id)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusNotFound)
	} else {
		core.RenderTemplate(w, "post.html", fetchSession(r), postComments, nosurf.Token(r))
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

type VoteType int

const (
	VOTETYPE_UP VoteType = 0
	VOTETYPE_DOWN VoteType = 1
)

func votePost(w http.ResponseWriter, r *http.Request, urlPath string, voteType VoteType) {
	if len(r.URL.Path) <= len(urlPath) {
		// ERROR: URL Path too short
		fmt.Println("URL Path too short")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	// Get session and username from it
	session := fetchSession(r)
	if session == nil {
		// ERROR: Session not found
		fmt.Println("Session not found")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	// Get post ID from URL Path
	postIdStr := r.URL.Path[len(urlPath):]
	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		// ERROR: Cannot convert to int
		fmt.Println("Cannot conver to int")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	// Update the database
	switch voteType {
	case VOTETYPE_UP:
		err = database.DBase.VotePost(session.Name, postId)
	case VOTETYPE_DOWN:
		err = database.DBase.DelVotePost(session.Name, postId)
	}
	if err != nil {
		// ERROR: Database
		fmt.Println("Database error")
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/post/"+postIdStr, http.StatusFound)
}

// Vote Handler | Handles voting up
func VotePostUp(w http.ResponseWriter, r *http.Request) {
	votePost(w, r, "/vote_up/", VOTETYPE_UP)
}

// Vote Handler | Handles voting down
func VotePostDown(w http.ResponseWriter, r *http.Request) {
	votePost(w, r, "/vote_down/", VOTETYPE_DOWN)
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

func Comment(w http.ResponseWriter, r *http.Request) {
	session := fetchSession(r)
	if session == nil {
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

	userId, err := database.DBase.GetLoginId(username)
	if err != nil {
		fmt.Println("Cannot get loginId")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	postIdStr := r.FormValue("post_id")

	postId, err := strconv.ParseInt(postIdStr, 10, 64)
	if err != nil {
		fmt.Println("Cannot get PostId")
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	comment := r.FormValue("comment")
	_, err = database.DBase.CommentCreate(core.Comment{
		UserId:  userId,
		PostId:  postId,
		Comment: comment})

	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, "/post/"+postIdStr, http.StatusFound)
}

// Error handler gets status and file rendered for the status
func errorHandler(w http.ResponseWriter, r *http.Request, status int, errorTmplFile string) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		core.RenderTemplate(w, errorTmplFile, fetchSession(r), nil, "")
	}
}
