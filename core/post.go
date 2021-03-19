package core

import (
	"fmt"
)

// Information about the post
type Post struct {
	Title   string
	Link    string
	Comment string
}

// Returns the string for the Post struct
func (p *Post) String() string {
	return fmt.Sprintf("%s (%s) | %s", p.Title, p.Link, p.Comment)
}
