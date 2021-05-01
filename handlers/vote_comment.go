package handlers

/*
import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/mtcw99/disnews/database"
)

func voteComment(w http.ResponseWriter, r *http.Request, urlPath string, voteType VoteType) {
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

	// Get comment ID from URL Path
	commentIdStr := r.URL.Path[len(urlPath):]
	commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
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
func VoteCommentUp(w http.ResponseWriter, r *http.Request) {
	voteComment(w, r, "/vote_up/", VOTETYPE_UP)
}

// Vote Handler | Handles voting down
func VoteCommentDown(w http.ResponseWriter, r *http.Request) {
	voteComment(w, r, "/vote_down/", VOTETYPE_DOWN)
}
*/
