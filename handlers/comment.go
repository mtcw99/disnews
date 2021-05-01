package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mtcw99/disnews/core"
	"github.com/mtcw99/disnews/database"
)

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
