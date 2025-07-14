package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
)

const (
	INVALID_USER_ID = -1
	ADMIN_USER_ID   = 1
	ADMIN_USER_NAME = "admin"
	GUEST_USER_ID   = 1000
	GUEST_USER_NAME = "guest"
)

type User struct {
	UserId             int      `json:"user_id"`
	UserType           UserType `json:"user_type"`
	UserName           string   `json:"user_name"`
	UserPasswordHashed string   `json:"user_password_hashed"`
	UserEmail          string   `json:"user_email"`
	UserSig            string   `json:"user_sig"`
	UserRegTime        int64    `json:"user_reg_time"`
	UserLastVisitTime  int64    `json:"user_last_visit_time"`

	// Derived properties
	UserNumPosts int    `json:"user_num_posts"`
	UserTypeName string `json:"user_type_name"`
	UserTypeImg  string `json:"user_type_img"`
}

// https://www.phpbb.com/community/viewtopic.php?t=1760075: "phpBB has a couple of special user types, those types are stored in "user_type" field. The definitions of those values are set in the "includes/constants.php" file. define('USER_NORMAL', 0); define('USER_INACTIVE', 1); define('USER_IGNORE', 2); define('USER_FOUNDER', 3);"
type UserType int

const (
	USER_TYPE_NORMAL   UserType = iota
	USER_TYPE_INACTIVE          // May be pending activation
	USER_TYPE_GUEST             // Guest users (Anonymous)
	USER_TYPE_FOUNDER  = 99     // Board founders (Super admins)
)

func InitUsers(ctx context.Context) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	// On PostgreSQL:  user_id INT(10) PRIMARY KEY AUTOINCREMENT NOT NULL,
	sql := `CREATE TABLE users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
		user_type TINYINT(2) NOT NULL DEFAULT '0',
		user_name VARCHAR(255) UNIQUE NOT NULL DEFAULT '',
		user_password_hashed VARCHAR(255) NOT NULL DEFAULT '',
		user_email VARCHAR(100) NOT NULL DEFAULT '',
		user_sig MEDIUMTEXT NOT NULL DEFAULT '',
		user_reg_time INT(11) NOT NULL DEFAULT '0',
		user_last_visit_time INT(11) NOT NULL DEFAULT '0',
		user_num_posts MEDIUMINT(8) NOT NULL DEFAULT '0'
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating users table: %s", err)
	}
	now := time.Now().UTC()
	userRegTime := now.Unix()

	// Admin user
	adminPassword, err := helper.GenerateRandomAlphanumeric(16)
	if err != nil {
		return fmt.Errorf("Error while generating password for Admin user: %s", err)
	}
	logger.Infof(ctx, "Username for Admin user: %s", ADMIN_USER_NAME)
	logger.Infof(ctx, "Password for Admin user: %s", adminPassword)
	adminSalt, err := helper.GenerateRandomBytesInHex(8)
	if err != nil {
		return fmt.Errorf("Error while generating random salt for Admin user: %s", err)
	}
	adminHashedPasswordWithSaltAndHeader := helper.HashPassword(adminPassword, adminSalt)
	_, err = db.Exec("INSERT INTO users (user_id, user_type, user_name, user_password_hashed, user_reg_time, user_last_visit_time) VALUES ($1, $2, $3, $4, $5, $6)", ADMIN_USER_ID, USER_TYPE_FOUNDER, ADMIN_USER_NAME, adminHashedPasswordWithSaltAndHeader, userRegTime, userRegTime)
	if err != nil {
		return fmt.Errorf("Error while inserting Admin user into users table: %s", err)
	}
	logger.Infof(ctx, "")

	// Guest user
	guestPassword, err := helper.GenerateRandomAlphanumeric(16)
	if err != nil {
		return fmt.Errorf("Error while generating password for Guest user: %s", err)
	}
	guestSalt, err := helper.GenerateRandomBytesInHex(8)
	if err != nil {
		return fmt.Errorf("Error while generating random salt for Guest user: %s", err)
	}
	guestHashedPasswordWithSaltAndHeader := helper.HashPassword(guestPassword, guestSalt)
	_, err = db.Exec("INSERT INTO users (user_id, user_type, user_name, user_password_hashed, user_reg_time, user_last_visit_time) VALUES ($1, $2, $3, $4, $5, $6)", GUEST_USER_ID, USER_TYPE_GUEST, GUEST_USER_NAME, guestHashedPasswordWithSaltAndHeader, userRegTime, userRegTime)
	if err != nil {
		return fmt.Errorf("Error while inserting Guest user into users table: %s", err)
	}

	return nil
}

func InsertUser(ctx context.Context, userName string, userPassword string, userEmail string, userSig string) (int, error) {
	// Check that user (lowercase) does not exist on database
	isExists, err := CheckIfUserExists(ctx, userName)
	if err != nil {
		return INVALID_USER_ID, fmt.Errorf("Error while inserting user name '%s' into users table: %s", userName, err)
	}
	if isExists {
		return INVALID_USER_ID, fmt.Errorf("Error while inserting user name '%s' into users table: %s: User already exists", userName, DB_ERROR_UNIQUE_CONSTRAINT)
	}

	// Proceed to insert user on database
	db := OpenDb(ctx, "users")
	defer db.Close()
	salt, err := helper.GenerateRandomBytesInHex(8)
	if err != nil {
		return INVALID_USER_ID, fmt.Errorf("Error while generating random salt for user name '%s': %s", userName, err)
	}
	hashedPasswordWithSaltAndHeader := helper.HashPassword(userPassword, salt)
	now := time.Now().UTC()
	userRegTime := now.Unix()
	result, err := db.Exec("INSERT INTO users (user_name, user_password_hashed, user_email, user_sig, user_reg_time, user_last_visit_time) VALUES ($1, $2, $3, $4, $5, $6)", userName, hashedPasswordWithSaltAndHeader, userEmail, userSig, userRegTime, userRegTime)
	if IsUniqueViolation(err) {
		return INVALID_USER_ID, fmt.Errorf("Error while inserting user name '%s' into users table: %s: %s", userName, DB_ERROR_UNIQUE_CONSTRAINT, err)
	}
	if err != nil {
		return INVALID_USER_ID, fmt.Errorf("Error while inserting user name '%s' into users table: %s", userName, err)
	}
	userId, err := result.LastInsertId()
	if err != nil {
		return INVALID_USER_ID, fmt.Errorf("Error while retrieving last insert id while inserting user name '%s' into users table: %s", userName, err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return INVALID_USER_ID, fmt.Errorf("Error while retrieving rows affected while inserting user name '%s' into users table: %s", userName, err)
	}
	if rowsAffected == 0 {
		return INVALID_USER_ID, fmt.Errorf("No rows were updated while inserting user name '%s' into users table", userName)
	}
	return int(userId), nil
}

func SetUserType(ctx context.Context, userId int, userType UserType) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	_, err := db.Exec("UPDATE users SET user_type = $1 WHERE user_id = $2", userType, userId)
	if err != nil {
		return fmt.Errorf("Error while setting user type %d for user id %d: %s", userType, userId, err)
	}
	return nil
}

func IncreaseNumPostsForUser(ctx context.Context, userId int) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	_, err := db.Exec("UPDATE users SET user_num_posts = user_num_posts + 1 WHERE user_id = $1", userId)
	if err != nil {
		return fmt.Errorf("Error while increasing num posts for user id %d: %s", userId, err)
	}
	return nil
}

func UpdateLastVisitTimeForUser(ctx context.Context, userId int) error {
	db := OpenDb(ctx, "users")
	defer db.Close()
	now := time.Now().UTC()
	userLastVisitTime := now.Unix()
	_, err := db.Exec("UPDATE users SET user_last_visit_time = $1 WHERE user_id = $2", userLastVisitTime, userId)
	if err != nil {
		return fmt.Errorf("Error while updating the last visit time for user id %d: %s", userId, err)
	}
	return nil
}

func CheckIfUserExists(ctx context.Context, userName string) (bool, error) {
	db := OpenDb(ctx, "users")
	defer db.Close()
	var isExists bool
	err := db.
		QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE LOWER(user_name) = $1)", strings.ToLower(userName)).
		Scan(&isExists)
	if err != nil {
		return false, fmt.Errorf("Error while checking if user name '%s' exists on users table: %s", userName, err)
	}
	return isExists, nil
}

func CheckIfGuestUser(ctx context.Context, userId int, userName string) bool {
	if userId == GUEST_USER_ID {
		return true
	}
	if userId == INVALID_USER_ID {
		return true
	}
	if userId == 0 {
		return true
	}
	if userName == "" {
		return true
	}
	return false
}

func GetUserForLogin(ctx context.Context, userName string) (User, error) {
	// WARNING: As this function returns sensitive information such as hashed password of user, use this function for login validation only
	db := OpenDb(ctx, "users")
	defer db.Close()
	var user User
	err := db.
		QueryRow("SELECT user_id, user_name, user_password_hashed FROM users WHERE LOWER(user_name) = $1", strings.ToLower(userName)).
		Scan(&user.UserId, &user.UserName, &user.UserPasswordHashed)
	if err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return User{}, fmt.Errorf("Error while retrieving user name '%s' on users table: %s: No result found", userName, DB_ERROR_NO_RESULT)
		}
		return User{}, fmt.Errorf("Error while scanning row on users table for user name '%s': %s", userName, err)
	}
	return user, nil
}

func ListUsersOfTopic(ctx context.Context, topicId int) ([]User, error) {
	// WARNING: Issue on Golang Template may reveal sensitive information of users. So avoid reading sensitive information here.
	db := OpenDb(ctx, "users")
	defer db.Close()
	rows, err := db.Query("SELECT DISTINCT users.user_id, users.user_type, users.user_name, users.user_sig, users.user_reg_time, users.user_last_visit_time, users.user_num_posts FROM users JOIN posts ON posts.post_user_id = users.user_id WHERE posts.topic_id = $1 ORDER BY users.user_id", topicId)
	if err != nil {
		return nil, fmt.Errorf("Error while querying users table for topic id %d: %s", topicId, err)
	}
	defer rows.Close()
	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.UserId, &user.UserType, &user.UserName, &user.UserSig, &user.UserRegTime, &user.UserLastVisitTime, &user.UserNumPosts); err != nil {
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
		} else if users[i].UserType == USER_TYPE_GUEST {
			users[i].UserTypeName = "Guest"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/guest.png"
		} else if users[i].UserType == USER_TYPE_FOUNDER {
			users[i].UserTypeName = "Administrator"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/administrator.png"
		} else {
			users[i].UserTypeName = "Unknown"
			users[i].UserTypeImg = "/images/ranks/modern-ranks/guest.png"
		}
	}
	return users, nil
}
