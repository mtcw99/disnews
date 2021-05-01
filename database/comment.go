package database

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/mtcw99/disnews/core"
)

func (d *Database) CommentCreate(comment core.Comment) (int64, error) {
	res, err := d.db.Exec(`
	INSERT INTO Comments(comment, user_id, post_id, creation_date)
	values(?, ?, ?, datetime('now'))
	`, comment.Comment, comment.UserId, comment.PostId)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (d *Database) GetComments(postId int64) ([]core.Comment, error) {
	rows, err := d.db.Query(`
	SELECT	Profiles.display_name,
		Users.name,
		Comments.id, Comments.post_id, Comments.creation_date, Comments.comment
	FROM Comments LEFT JOIN Users
	ON Comments.user_id = Users.id
	LEFT JOIN Profiles
	ON Users.profile_id = Profiles.id
	WHERE post_id=?
	`, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []core.Comment

	for rows.Next() {
		var comment core.Comment

		err = rows.Scan(&comment.DisplayName,
			&comment.Username,
			&comment.Id,
			&comment.PostId,
			&comment.DateTime,
			&comment.Comment)

		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return comments, nil
}
