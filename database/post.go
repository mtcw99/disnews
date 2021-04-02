package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mtcw99/disnews/core"
)

// Submit a new post into the database
func (d *Database) SubmitPost(post core.Post) (int64, error) {
	userid, err := d.GetLoginId(post.User)
	if err != nil {
		return 0, err
	}

	res, err := d.db.Exec(`
	INSERT INTO Posts(title, link, comment, creation_date, user_id)
	values(?, ?, ?, datetime('now'), ?)
	`, post.Title, post.Link, post.Comment, userid)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) GetPostAndComments(indexStr string) (core.PostComments, error) {
	row := d.db.QueryRow(`
	SELECT id, title, link, comment, creation_date, user_id FROM Posts
	WHERE id=?
	`, indexStr)

	if row == nil {
		return core.PostComments{}, fmt.Errorf("ERROR: Database.GetPostAndComments: id %s not found.", indexStr)
	}

	var postComments core.PostComments
	var userid int64
	err := row.Scan(&postComments.Post.Id,
		&postComments.Post.Title,
		&postComments.Post.Link,
		&postComments.Post.Comment,
		&postComments.Post.Date,
		&userid)
	if err != nil {
		return core.PostComments{}, err
	}

	// Gets votes of this post
	postComments.Post.Votes, err = d.GetVotes(postComments.Post.Id)
	if err != nil {
		return core.PostComments{}, err
	}

	postComments.Post.User, err = d.GetLoginUserFromId(userid)
	if err != nil {
		return core.PostComments{}, err
	}

	postComments.Comments, err = d.GetComments(postComments.Post.Id)
	if err != nil {
		return core.PostComments{}, err
	}

	return postComments, nil
}

// Gets the newests posts from the database
func (d *Database) GetNewestPosts() ([]core.Post, error) {
	rows, err := d.db.Query(`
	SELECT id, title, link, comment, creation_date, user_id FROM Posts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post

	for rows.Next() {
		var post core.Post
		var userid int64
		err = rows.Scan(&post.Id, &post.Title, &post.Link,
			&post.Comment, &post.Date, &userid)
		if err != nil {
			return nil, err
		}

		post.User, err = d.GetLoginUserFromId(userid)
		if err != nil {
			return nil, err
		}

		post.Votes, err = d.GetVotes(post.Id)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return posts, nil
}
