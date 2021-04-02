// Package database provides functionalities to interact and set/get information to/from the database.
// Includes login and post functionalities
package database

import (
	"database/sql"
	"log"
	"os"

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
	CREATE TABLE IF NOT EXISTS Profiles (
		  id INTEGER NOT NULL PRIMARY KEY
		, display_name TEXT NOT NULL
		, info TEXT NOT NULL
		, link TEXT NOT NULL
		, creation_date TEXT NOT NULL
		, login_date TEXT
	);
	CREATE TABLE IF NOT EXISTS Users (
		  id INTEGER NOT NULL PRIMARY KEY
		, name TEXT NOT NULL UNIQUE
		, pass TEXT NOT NULL UNIQUE
		, profile_id INTEGER NOT NULL UNIQUE
		, CONSTRAINT fk_profile_id
			FOREIGN KEY(profile_id)
			REFERENCES Profiles(id)
			ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Posts (
		  id INTEGER NOT NULL PRIMARY KEY
		, title TEXT
		, link TEXT
		, comment TEXT
		, creation_date TEXT NOT NULL
		, user_id INTEGER
		, CONSTRAINT fk_user_id
			FOREIGN KEY(user_id)
			REFERENCES Users(id)
			ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS Comments (
		  id INTEGER NOT NULL PRIMARY KEY
		, comment TEXT
		, user_id INTEGER
		, post_id INTEGER
		, creation_date TEXT
		, CONSTRAINT fk_user_id
			FOREIGN KEY(user_id)
			REFERENCES Users(id)
			ON DELETE CASCADE
		, CONSTRAINT fk_post_id
			FOREIGN KEY(post_id)
			REFERENCES Posts(id)
			ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS VotesPosts (
		  id INTEGER NOT NULL PRIMARY KEY
		, user_id INTEGER
		, post_id INTEGER
		, CONSTRAINT fk_user_id
			FOREIGN KEY(user_id)
			REFERENCES Users(id)
			ON DELETE CASCADE
		, CONSTRAINT fk_post_id
			FOREIGN KEY(post_id)
			REFERENCES Posts(id)
			ON DELETE CASCADE
		, UNIQUE(user_id, post_id)
	);
	`

	_, err := d.db.Exec(sqlst)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlst)
		return err
	}

	return nil
}

// Closes the database
func (d *Database) Close() {
	d.db.Close()
}
