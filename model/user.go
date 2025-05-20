package model

import (
	"context"
	"fmt"
	"time"

	"phpbb-golang/internal/helper"
)

type User struct {
	UserId             int    `json:"user_id"`
	UserName           string `json:"user_name"`
	UserPasswordHashed string `json:"user_password_hashed"`
	UserSig            string `json:"user_sig"`
	UserRegTime        int64  `json:"user_reg_time"`
}

func InitUsers(ctx context.Context) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	// On PostgreSQL:  user_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		user_name VARCHAR(255) NOT NULL DEFAULT '',
		user_password_hashed VARCHAR(255) NOT NULL DEFAULT '',
		user_sig MEDIUMTEXT NOT NULL DEFAULT '',
		user_reg_time INT(11) NOT NULL DEFAULT '0'
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating users table: %s", err)
	}
	return nil
}

func InsertUser(ctx context.Context, userName string, userPassword string, userSig string) (int, error) {
	db := OpenDb(ctx, "users")
	defer db.Close()
	salt := helper.GenerateRandomSalt(4)
	hashedPasswordWithSalt := helper.HashPassword(userPassword, salt)
	now := time.Now().UTC()
	userRegTime := now.Unix()
	res, err := db.Exec(`INSERT INTO users (user_name, user_password_hashed, user_sig, user_reg_time) VALUES ($1, $2, $3, $4)`, userName, hashedPasswordWithSalt, userSig, userRegTime)
	if err != nil {
		return -1, fmt.Errorf("Error while inserting user '%s' into users table: %s", userName, err)
	}
	userId, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("Error while retrieving last insert id for user '%s': %s", userName, err)
	}
	return int(userId), nil
}

func ListUsers(ctx context.Context, topicId int) ([]User, error) {
	// Warning: Issue on Golang Template may reveal sensitive information of users. So avoid reading sensitive information here.
	db := OpenDb(ctx, "users")
	defer db.Close()
	// ChatGPT: SQL Database with "users" and "posts" table. A user may post multiple things. Now generate SQL SELECT statement to list unique users given a post id.
	rows, err := db.Query("SELECT DISTINCT users.user_id, users.user_name, users.user_sig, users.user_reg_time FROM users JOIN posts ON posts.user_id = users.user_id WHERE posts.topic_id = $1 ORDER BY users.user_id", topicId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying users table: %s", err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserId, &user.UserName, &user.UserSig, &user.UserRegTime); err != nil {
			return nil, fmt.Errorf("Error while scanning rows on users table: %s", err)
		}
		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error on rows on users table: %s", err)
	}
	return users, nil
}
