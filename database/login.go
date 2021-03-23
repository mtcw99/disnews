package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mtcw99/disnews/core"
)

func (d *Database) Signup(login core.Login) error {
	profileid, err := d.CreateProfile(login.Username)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
	INSERT INTO Users(name, pass, profile_id)
	values(?, ?, ?)
	`, login.Username, login.Hashpass, profileid)
	if err != nil {
		errp := d.DeleteProfile(profileid)
		if errp != nil {
			return errp
		} else {
			return err
		}
	}

	return nil
}

//func (d *Database) DeleteUser(userid int64) error {
//}

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

func (d *Database) GetLoginId(username string) (int, error) {
	row := d.db.QueryRow(`
	SELECT id FROM Users
	WHERE name=?
	`, username)

	if row == nil {
		return 0, fmt.Errorf("ERROR: Database.GetLoginId: username %s not found.", username)
	}

	var id int
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	} else {
		return id, nil
	}
}

func (d *Database) GetLoginUserFromId(id int) (string, error) {
	row := d.db.QueryRow(`
	SELECT name FROM Users
	WHERE id=?
	`, id)

	if row == nil {
		return "", fmt.Errorf("ERROR: Database.GetLoginUserFromId: id %d not found.", id)
	}

	var username string
	err := row.Scan(&username)
	if err != nil {
		return "", err
	} else {
		return username, nil
	}
}
