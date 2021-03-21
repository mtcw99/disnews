package core

import (
	"fmt"
)

// Information about the post
type Post struct {
	Id	int64
	User	string
	Title   string
	Link    string
	Comment string
}

// Returns the string for the Post struct
func (p *Post) String() string {
	return fmt.Sprintf("%d: %s (%s) | %s", p.Id, p.Title, p.Link, p.Comment)
}
