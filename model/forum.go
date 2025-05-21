package model

import (
	"context"
	"database/sql"
	"fmt"
)

type Forum struct {
	ForumId   int    `json:"forum_id"`
	ParentId  int    `json:"parent_id"`
	ForumName string `json:"forum_name"`
	ForumDesc string `json:"forum_desc"`
}

func InitForums(ctx context.Context) error {
	// Ref: https://www.erdcloud.com/d/23zvQbme2zHiLtYmf
	db := OpenDb(ctx, "forums")
	defer db.Close()
	// On PostgreSQL:  forum_id MEDIUMINT(8) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE forums (
		forum_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		parent_id MEDIUMINT(8) NOT NULL DEFAULT '0',
		forum_name VARCHAR(255) NOT NULL DEFAULT '',
		forum_desc TEXT NOT NULL DEFAULT '',
		FOREIGN KEY (parent_id) REFERENCES forums(forum_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating forums table: %s", err)
	}
	_, err = db.Exec(`INSERT INTO forums (forum_id, parent_id, forum_name, forum_desc) VALUES ($1, $2, $3, $4)`, 0, 0, "Root Forum", "")
	if err != nil {
		return fmt.Errorf("Error while inserting root forum into forums table: %s", err)
	}
	return nil
}

func InsertForum(ctx context.Context, parentId int, forumName string, forumDesc string) (int, error) {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	res, err := db.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES ($1, $2, $3)`, parentId, forumName, forumDesc)
	if err != nil {
		return -1, fmt.Errorf("Error while inserting forum name '%s' with forum description '%s' and parent forum %d into forums table: %s", forumName, forumDesc, parentId, err)
	}
	forumId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Error while retrieving last insert id for forum name '%s': %s", forumName, err)
	}
	return int(forumId), nil
}

func ListForums(ctx context.Context) ([]Forum, error) {
	// Ref: https://go.dev/doc/tutorial/database-access
	db := OpenDb(ctx, "forums")
	defer db.Close()
	rows, err := db.Query("SELECT forum_id, parent_id, forum_name, forum_desc FROM forums ORDER BY forum_id")
	if err != nil {
		return nil, fmt.Errorf("Error while querying forums table: %s", err)
	}
	defer rows.Close()
	var forums []Forum
	for rows.Next() {
		var forum Forum
		if err := rows.Scan(&forum.ForumId, &forum.ParentId, &forum.ForumName, &forum.ForumDesc); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on forums table: %s", err)
		}
		forums = append(forums, forum)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on forums table: %s", err)
	}
	return forums, nil
}

func GetForum(ctx context.Context, forumId int) (Forum, error) {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	row := db.QueryRow("SELECT forum_id, parent_id, forum_name, forum_desc FROM forums WHERE forum_id = $1", forumId)
	var forum Forum
	if err := row.Scan(&forum.ForumId, &forum.ParentId, &forum.ForumName, &forum.ForumDesc); err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return Forum{}, nil
		}
		return Forum{}, fmt.Errorf("Error while scanning row on forums table for forum id %d: %s", forumId, err)
	}
	return forum, nil
}
