package core

import (
	"database/sql"
	"log"
	"os"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Holds the database session
type Database struct {
	db *sql.DB
}

// Global database session
var DBase Database

// Checks if the file exists
func (d *Database) Check(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// Create a new sqlite3 database session
func (d *Database) New(path string) error {
	var err error
	d.db, err = sql.Open("sqlite3", path)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

// Setup the database
func (d *Database) Setup() error {
	sqlst := `
	CREATE TABLE Posts (
		  id INTEGER NOT NULL PRIMARY KEY
		, title TEXT
		, link TEXT
		, comment TEXT
	);
	`

	_, err := d.db.Exec(sqlst)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlst)
		return err
	}

	return nil
}

// Submit a new post into the database
func (d *Database) SubmitPost(post Post) (int64, error) {
	tx, err := d.db.Begin()
	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	statement, err := tx.Prepare(`
	INSERT INTO Posts(title, link, comment)
	values(?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	defer statement.Close()

	res, err := statement.Exec(post.Title,
		post.Link,
		post.Comment)

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		log.Fatal(err)
		return 0, err
	}

	tx.Commit()
	return id, nil
}

func (d *Database) GetPost(indexStr string) (Post, error) {
	row := d.db.QueryRow(`
	SELECT title, link, comment FROM Posts
	WHERE id=?
	`, indexStr)

	if row == nil {
		return Post{}, fmt.Errorf("ERROR: core: Database.GetPost: id %s not found.", indexStr)
	}

	var post Post
	err := row.Scan(&post.Title, &post.Link, &post.Comment)
	if err != nil {
		return Post{}, err
	} else {
		return post, nil
	}
}

// Gets the newests posts from the database
func (d *Database) GetNewestPosts() ([]Post, error) {
	rows, err := d.db.Query(`
	SELECT title, link, comment FROM Posts
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Title, &post.Link, &post.Comment)
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

// Closes the database
func (d *Database) Close() {
	d.db.Close()
}
