package database

import (
	"fmt"

	"github.com/mtcw99/disnews/core"
	_ "github.com/mattn/go-sqlite3"
)

func (d *Database) Signup(login core.Login) error {
	_, err := d.db.Exec(`
	INSERT INTO Users(name, pass)
	values(?, ?)
	`, login.Username, login.Hashpass)

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (d *Database) Login(username string) (core.Login, error) {
	row := d.db.QueryRow(`
	SELECT pass FROM Users
	WHERE name=?
	`, username)

	if row == nil {
		return core.Login{}, fmt.Errorf("ERROR: Database.Login: username %s not found.", username)
	}

	var login core.Login
	login.Username = username
	err := row.Scan(&login.Hashpass)
	if err != nil {
		return core.Login{}, err
	} else {
		return login, nil
	}
}

