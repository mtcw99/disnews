package database

import (
	"fmt"

	"github.com/mtcw99/disnews/core"
)

func (d *Database) CreateProfile(displayName string) (int64, error) {
	res, err := d.db.Exec(`
	INSERT INTO Profiles(display_name, info, link, creation_date)
	values(?, ?, ?, datetime('now'))
	`, displayName, "", "")
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) DeleteProfile(profileid int64) error {
	_, err := d.db.Exec(`
	DELETE FROM Profiles
	WHERE id=?
	`, profileid)
	if err != nil {
		return err
	}
	return nil
}

func (d *Database) GetProfile(username string) (core.Profile, error) {
	row := d.db.QueryRow(`
	SELECT Profiles.display_name, Profiles.info, Profiles.link, Profiles.creation_date
	FROM Profiles LEFT JOIN Users
	ON Profiles.id = Users.profile_id
	WHERE Users.name=?
	`, username)

	if row == nil {
		return core.Profile{}, fmt.Errorf("ERROR: Database.GetProfile: name %s not found.", username)
	}

	var profile core.Profile
	profile.Username = username
	err := row.Scan(&profile.DisplayName, &profile.Info, &profile.Link, &profile.CreationDate)
	if err != nil {
		return core.Profile{}, err
	}

	return profile, nil
}

func (d *Database) GetProfileId(username string) (int64, error) {
	row := d.db.QueryRow(`
	SELECT Profiles.id
	FROM Profiles LEFT JOIN Users
	ON Profiles.id = Users.profile_id
	WHERE Users.name=?
	`, username)

	if row == nil {
		return 0, fmt.Errorf("ERROR: Database.GetProfileId: name %s not found.", username)
	}

	var profileid int64
	err := row.Scan(&profileid)
	if err != nil {
		return 0, err
	}

	return profileid, nil
}

func (d *Database) UpdateProfile(profileid int64,
	displayName string,
	info string,
	link string) error {
	_, err := d.db.Exec(`
	UPDATE Profiles
	SET display_name=?, info=?, link=?
	WHERE id=?
	`, displayName, info, link, profileid)
	if err != nil {
		return err
	}
	return nil
}
