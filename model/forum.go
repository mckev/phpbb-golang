package model

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
)

const (
	ROOT_FORUM_ID              = 0
	INVALID_FORUM_ID           = -1
	MAX_FORUM_NAV_TRAILS_DEPTH = 7
	MAX_FORUM_NODES_DEPTH      = 3
)

type Forum struct {
	ForumId   int    `json:"forum_id"`
	ParentId  int    `json:"parent_id"`
	ForumName string `json:"forum_name"`
	ForumDesc string `json:"forum_desc"`
	// Derived properties to speed up
	ForumNumTopics int `json:"forum_num_topics"`
	ForumNumPosts  int `json:"forum_num_posts"`
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
		forum_num_topics MEDIUMINT(8) NOT NULL DEFAULT '0',
		forum_num_posts MEDIUMINT(8) NOT NULL DEFAULT '0',
		FOREIGN KEY (parent_id) REFERENCES forums(forum_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating forums table: %s", err)
	}
	_, err = db.Exec(`INSERT INTO forums (forum_id, parent_id, forum_name, forum_desc) VALUES ($1, $2, $3, $4)`, ROOT_FORUM_ID, ROOT_FORUM_ID, "Root Forum", "")
	if err != nil {
		return fmt.Errorf("Error while inserting Root Forum into forums table: %s", err)
	}
	return nil
}

func InsertForum(ctx context.Context, parentId int, forumName string, forumDesc string) (int, error) {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	res, err := db.Exec(`INSERT INTO forums (parent_id, forum_name, forum_desc) VALUES ($1, $2, $3)`, parentId, forumName, forumDesc)
	if err != nil {
		return INVALID_FORUM_ID, fmt.Errorf("Error while inserting forum name '%s' with forum description '%s' and parent forum %d into forums table: %s", forumName, forumDesc, parentId, err)
	}
	forumId, err := res.LastInsertId()
	if err != nil {
		return INVALID_FORUM_ID, fmt.Errorf("Error while retrieving last insert id for forum name '%s': %s", forumName, err)
	}
	return int(forumId), nil
}

func IncreaseNumTopicsForForum(ctx context.Context, forumId int) error {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	_, err := db.Exec(`UPDATE forums SET forum_num_topics = forum_num_topics + 1 WHERE forum_id = $1`, forumId)
	if err != nil {
		return fmt.Errorf("Error while increasing num topics for forum id %d: %s", forumId, err)
	}
	return nil
}

func IncreaseNumPostsForForum(ctx context.Context, forumId int) error {
	db := OpenDb(ctx, "forums")
	defer db.Close()
	_, err := db.Exec(`UPDATE forums SET forum_num_posts = forum_num_posts + 1 WHERE forum_id = $1`, forumId)
	if err != nil {
		return fmt.Errorf("Error while increasing num posts for forum id %d: %s", forumId, err)
	}
	return nil
}

func ListForums(ctx context.Context) ([]Forum, error) {
	// Ref: https://go.dev/doc/tutorial/database-access
	db := OpenDb(ctx, "forums")
	defer db.Close()
	rows, err := db.Query("SELECT forum_id, parent_id, forum_name, forum_desc, forum_num_topics, forum_num_posts FROM forums ORDER BY forum_id")
	if err != nil {
		return nil, fmt.Errorf("Error while querying forums table: %s", err)
	}
	defer rows.Close()
	var forums []Forum
	for rows.Next() {
		var forum Forum
		if err := rows.Scan(&forum.ForumId, &forum.ParentId, &forum.ForumName, &forum.ForumDesc, &forum.ForumNumTopics, &forum.ForumNumPosts); err != nil {
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
	row := db.QueryRow("SELECT forum_id, parent_id, forum_name, forum_desc, forum_num_topics, forum_num_posts FROM forums WHERE forum_id = $1", forumId)
	var forum Forum
	if err := row.Scan(&forum.ForumId, &forum.ParentId, &forum.ForumName, &forum.ForumDesc, &forum.ForumNumTopics, &forum.ForumNumPosts); err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return Forum{}, nil
		}
		return Forum{}, fmt.Errorf("Error while scanning row on forums table for forum id %d: %s", forumId, err)
	}
	return forum, nil
}

type ForumNavTrail struct {
	Forum  Forum
	IsLeaf bool
}

func ComputeForumNavTrails(ctx context.Context, forumId int) ([]ForumNavTrail, error) {
	// Given a Forum Id, find its parents until Root Forum
	forums, err := ListForums(ctx)
	if err != nil {
		return []ForumNavTrail{}, fmt.Errorf("Error while listing forums upon computing Forum Nav Trails: %s", err)
	}
	// Convert users from a list into a map
	forumsMap := map[int]Forum{}
	for _, forum := range forums {
		forumsMap[forum.ForumId] = forum
	}
	forumNavTrails := []ForumNavTrail{}
	depth := 0
	id := forumId
	for true {
		if id == ROOT_FORUM_ID {
			break
		}
		forum := forumsMap[id]
		forumNavTrails = append(forumNavTrails, ForumNavTrail{
			Forum:  forum,
			IsLeaf: false,
		})
		depth++
		if depth > MAX_FORUM_NAV_TRAILS_DEPTH {
			return []ForumNavTrail{}, fmt.Errorf("Error while computing Forum Nav Trails for forum id %d: Path too deep", forumId)
		}
		id = forum.ParentId
	}
	if len(forumNavTrails) > 0 {
		slices.Reverse(forumNavTrails)
		forumNavTrails[len(forumNavTrails)-1].IsLeaf = true
	}
	return forumNavTrails, nil
}

type ForumNode struct {
	Forum           Forum
	IsLeaf          bool
	ForumChildNodes []ForumNode
}

func ComputeForumChildNodes(ctx context.Context, forums []Forum, forumId int, depth int) []ForumNode {
	// Construct child nodes whose parent is Forum Id
	if depth >= MAX_FORUM_NODES_DEPTH {
		return []ForumNode{}
	}
	forumChildNodes := []ForumNode{}
	for _, forum := range forums {
		if forum.ForumId == ROOT_FORUM_ID {
			continue
		}
		if forum.ParentId == forumId {
			forumGrandchildNodes := ComputeForumChildNodes(ctx, forums, forum.ForumId, depth+1)
			isLeaf := false
			if len(forumGrandchildNodes) == 0 {
				isLeaf = true
			}
			forumChildNodes = append(forumChildNodes, ForumNode{
				Forum:           forum,
				IsLeaf:          isLeaf,
				ForumChildNodes: forumGrandchildNodes,
			})
		}
	}
	return forumChildNodes
}
