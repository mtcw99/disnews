package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mtcw99/disnews/core"
)

func (d *Database) VoteComment(username string, commentId int64) error {
	userId, err := d.GetLoginId(username)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
	INSERT INTO VotesComments(user_id, comment_id)
	values(?, ?)
	`, userId, commentId)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) DelVoteComment(username string, commentId int64) error {
	userId, err := d.GetLoginId(username)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
	DELETE FROM VotesComments
	WHERE user_id=? AND comment_id=?
	`, userId, commentId)
	if err != nil {
		return err
	}

	return nil
}

func (d *Database) GetCommentVotes(commentId int64) (int64, error) {
	var count int64
	err := d.db.QueryRow(`
	SELECT COUNT(comment_id)
	FROM VotesComments
	WHERE comment_id=?`, commentId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (d *Database) GetCommentVote(username string, commentId int64) (core.UserVoteType, error) {
	userId, err := d.GetLoginId(username)
	if err != nil {
		return core.USERVOTETYPE_NONE, err
	}

	var count int64
	err = d.db.QueryRow(`
	SELECT COUNT(comment_id)
	FROM VotesComments
	WHERE comment_id=? AND user_id=?
	`, commentId, userId).Scan(&count)

	if count == 1 {
		return core.USERVOTETYPE_UP, nil
	} else {
		return core.USERVOTETYPE_NONE, nil
	}
}
