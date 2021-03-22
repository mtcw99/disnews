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
	INSERT INTO Posts(title, link, comment, user_id)
	values(?, ?, ?, ?)
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

func (d *Database) GetPost(indexStr string) (core.Post, error) {
	row := d.db.QueryRow(`
	SELECT title, link, comment, user_id FROM Posts
	WHERE id=?
	`, indexStr)

	if row == nil {
		return core.Post{}, fmt.Errorf("ERROR: Database.GetPost: id %s not found.", indexStr)
	}

	var post core.Post
	var userid int
	err := row.Scan(&post.Title, &post.Link, &post.Comment, &userid)
	if err != nil {
		return core.Post{}, err
	}

	post.User, err = d.GetLoginUserFromId(userid)
	if err != nil {
		return core.Post{}, err
	} else {
		return post, nil
	}
}

// Gets the newests posts from the database
func (d *Database) GetNewestPosts() ([]core.Post, error) {
	rows, err := d.db.Query(`
	SELECT id, title, link, comment, user_id FROM Posts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []core.Post

	for rows.Next() {
		var post core.Post
		var userid int
		err = rows.Scan(&post.Id, &post.Title, &post.Link, &post.Comment, &userid)
		if err != nil {
			return nil, err
		}

		post.User, err = d.GetLoginUserFromId(userid)
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
