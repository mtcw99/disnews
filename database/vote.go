package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mtcw99/disnews/core"
)

func (d *Database) VotePost(username string, postId int64) error {
	userId, err := d.GetLoginId(username)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
	INSERT INTO VotesPosts(user_id, post_id)
	values(?, ?)
	`, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DelVotePost(username string, postId int64) error {
	userId, err := d.GetLoginId(username)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
	DELETE FROM VotesPosts
	WHERE user_id=? AND post_id=?
	`, userId, postId)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetVotes(postId int64) (int64, error) {
	var count int64
	err := d.db.QueryRow(`
	SELECT COUNT(post_id)
	FROM VotesPosts
	WHERE post_id=?`, postId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *Database) GetVote(username string, postId int64) (core.UserVoteType, error) {
	return core.USERVOTETYPE_NONE, nil
}
