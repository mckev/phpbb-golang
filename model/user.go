package model

import (
	"context"
	"fmt"
	"time"

	"phpbb-golang/internal/helper"
)

const (
	INVALID_USER_ID     = -1
	FIRST_ADMIN_USER_ID = 1
)

type User struct {
	UserId             int      `json:"user_id"`
	UserType           UserType `json:"user_type"`
	UserName           string   `json:"user_name"`
	UserPasswordHashed string   `json:"user_password_hashed"`
	UserSig            string   `json:"user_sig"`
	UserRegTime        int64    `json:"user_reg_time"`
	// Derived properties
	UserNumPosts int `json:"user_num_posts"`
	UserTypeName string
	UserTypeImg  string
}

// https://www.phpbb.com/community/viewtopic.php?t=1760075: "phpBB has a couple of special user types, those types are stored in "user_type" field. The definitions of those values are set in the "includes/constants.php" file. define('USER_NORMAL', 0); define('USER_INACTIVE', 1); define('USER_IGNORE', 2); define('USER_FOUNDER', 3);"
type UserType int

const (
	USER_TYPE_NORMAL UserType = iota
	USER_TYPE_INACTIVE
	USER_TYPE_IGNORE
	USER_TYPE_FOUNDER = 99
)

func InitUsers(ctx context.Context) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	// On PostgreSQL:  user_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		user_type TINYINT(2) NOT NULL DEFAULT '0',
		user_name VARCHAR(255) NOT NULL DEFAULT '',
		user_password_hashed VARCHAR(255) NOT NULL DEFAULT '',
		user_sig MEDIUMTEXT NOT NULL DEFAULT '',
		user_reg_time INT(11) NOT NULL DEFAULT '0',
		user_num_posts MEDIUMINT(8) NOT NULL DEFAULT '0'
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating users table: %s", err)
	}
	// First Admin user
	salt := helper.GenerateRandomSalt(4)
	userPassword := helper.GenerateRandomSalt(16)
	hashedPasswordWithSaltAndHeader := helper.HashPassword(userPassword, salt)
	now := time.Now().UTC()
	userRegTime := now.Unix()
	_, err = db.Exec(`INSERT INTO users (user_id, user_type, user_name, user_password_hashed, user_sig, user_reg_time) VALUES ($1, $2, $3, $4, $5, $6)`, FIRST_ADMIN_USER_ID, USER_TYPE_FOUNDER, "The Admin", hashedPasswordWithSaltAndHeader, "", userRegTime)
	if err != nil {
		return fmt.Errorf("Error while inserting First Admin user into users table: %s", err)
	}
	return nil
}

func InsertUser(ctx context.Context, userName string, userPassword string, userSig string) (int, error) {
	db := OpenDb(ctx, "users")
	defer db.Close()
	salt := helper.GenerateRandomSalt(4)
	hashedPasswordWithSaltAndHeader := helper.HashPassword(userPassword, salt)
	now := time.Now().UTC()
	userRegTime := now.Unix()
	res, err := db.Exec(`INSERT INTO users (user_name, user_password_hashed, user_sig, user_reg_time) VALUES ($1, $2, $3, $4)`, userName, hashedPasswordWithSaltAndHeader, userSig, userRegTime)
	if err != nil {
		return INVALID_USER_ID, fmt.Errorf("Error while inserting user name '%s' into users table: %s", userName, err)
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return INVALID_USER_ID, fmt.Errorf("Error while retrieving last insert id for user name '%s': %s", userName, err)
	}
	return int(userId), nil
}

func SetUserType(ctx context.Context, userId int, userType UserType) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	result, err := db.Exec(`UPDATE users SET user_type = $1 WHERE user_id = $2`, userType, userId)
	if err != nil {
		return fmt.Errorf("Error while setting user type %d for user id %d: %s", userType, userId, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("Error while retrieving rows affected while setting user type %d for user id %d: %s", userType, userId, err)
	}
	// TODO: Behavior may not be correct if there is no change to user type
	if rowsAffected == 0 {
		return fmt.Errorf("No rows were updated while setting user type %d for user id %d", userType, userId)
	}
	return nil
}

func IncreaseNumPostsForUser(ctx context.Context, userId int) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	_, err := db.Exec(`UPDATE users SET user_num_posts = user_num_posts + 1 WHERE user_id = $1`, userId)
	if err != nil {
		return fmt.Errorf("Error while increasing num posts for user id %d: %s", userId, err)
	}
	return nil
}

func ListUsers(ctx context.Context, topicId int) ([]User, error) {
	// WARNING: Issue on Golang Template may reveal sensitive information of users. So avoid reading sensitive information here.
	db := OpenDb(ctx, "users")
	defer db.Close()
	// ChatGPT: SQL Database with "users" and "posts" table. A user may post multiple things. Now generate SQL SELECT statement to list unique users given a post id.
	rows, err := db.Query("SELECT DISTINCT users.user_id, users.user_type, users.user_name, users.user_sig, users.user_reg_time, user_num_posts FROM users JOIN posts ON posts.post_user_id = users.user_id WHERE posts.topic_id = $1 ORDER BY users.user_id", topicId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying users table for topic id %d: %s", topicId, err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserId, &user.UserType, &user.UserName, &user.UserSig, &user.UserRegTime, &user.UserNumPosts); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on users table for topic id %d: %s", topicId, err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on users table for topic id %d: %s", topicId, err)
	}
	// Derived properties
	for i := range users {
		if users[i].UserType == USER_TYPE_NORMAL {
			users[i].UserTypeName = "Member"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/member.png"
		} else if users[i].UserType == USER_TYPE_INACTIVE {
			users[i].UserTypeName = "Member (inactive)"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/member.png"
		} else if users[i].UserType == USER_TYPE_IGNORE {
			users[i].UserTypeName = "Banned"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/banned.png"
		} else if users[i].UserType == USER_TYPE_FOUNDER {
			users[i].UserTypeName = "Administrator"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/administrator.png"
		} else {
			users[i].UserTypeName = "Guest"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/guest.png"
		}
	}
	return users, nil
}
