package model

import (
	"context"
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
	_, err = db.Exec(`INSERT INTO forums (forum_id, parent_id, forum_name, forum_desc) VALUES ($1, $2, $3, $4)`, 0, 0, "root", "")
	if err != nil {
		return fmt.Errorf("Error while inserting root forums table: %s", err)
	}
	return nil
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

func InsertForum(ctx context.Context, parentId int, forumName string, forumDesc string) (int, error) {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	res, err := db.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES ($1, $2, $3)`, parentId, forumName, forumDesc)
	if err != nil {
		return -1, fmt.Errorf("Error while inserting forum %s with description %s and parent %d into forums table: %s", forumName, forumDesc, parentId, err)
	}
	forumId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Error while retrieving last insert id: %s", err)
	}
	return int(forumId), nil
}
