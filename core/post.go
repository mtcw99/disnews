package core

import (
	"fmt"
)

// Information about the post
type Post struct {
	Id      int64
	User    string
	Title   string
	Link    string
	Comment string
	Date    string
}

type PostComments struct {
	Post     Post
	Comments []Comment
}

// Returns the string for the Post struct
func (p *Post) String() string {
	return fmt.Sprintf("%d: %s (%s) | %s", p.Id, p.Title, p.Link, p.Comment)
}

func LinkFix(link string) string {
	if len(link) <= 8 {
		return "https://" + link
	} else if link[:len("https://")] != "https://" &&
		link[:len("http://")] != "http://" {
		return "https://" + link
	} else {
		return link
	}
}
