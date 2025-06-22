package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

const (
	SESSION_TIMEOUT_IN_SECONDS = 4 * 3600
)

type Session struct {
	SessionId           string `json:"session_id"`
	SessionUserId       int    `json:"session_user_id"`
	SessionTimeStart    int64  `json:"session_time_start"`
	SessionTimeLast     int64  `json:"session_time_last"`
	SessionIp           string `json:"session_ip"`
	SessionBrowser      string `json:"session_browser"`
	SessionForwardedFor string `json:"session_forwarded_for"`
	// Derived properties to speed up
	SessionUserName string `json:"session_user_name"`
}

func InitSessions(ctx context.Context) error {
	db := OpenDb(ctx, "sessions")
	defer db.Close()
	// On PostgreSQL:  session_id char(32) PRIMARY KEY NOT NULL,
	// Notes:
	//    - session_forwarded_for is used to store Proxy Server info
	sql := `CREATE TABLE sessions (
		session_id char(32) PRIMARY KEY NOT NULL,
		session_user_id INT(10) NOT NULL DEFAULT '0',
		session_user_name VARCHAR(255) NOT NULL DEFAULT '',
		session_time_start INT(11) NOT NULL DEFAULT '0',
		session_time_last INT(11) NOT NULL DEFAULT '0',
		session_ip VARCHAR(40) NOT NULL DEFAULT '',
		session_browser VARCHAR(150) NOT NULL DEFAULT '',
		session_forwarded_for VARCHAR(255) NOT NULL DEFAULT '',
		FOREIGN KEY (session_user_id) REFERENCES users(user_id)
	)`
	_, err := db.Exec(sql)
	if err != nil {
		return fmt.Errorf("Error while creating sessions table: %s", err)
	}
	return nil
}

func CreateSession(ctx context.Context, sessionId string, userId int, userName string, ip string, browser string, forwardedFor string) (Session, error) {
	db := OpenDb(ctx, "sessions")
	defer db.Close()
	now := time.Now().UTC()
	sessionTimeStart := now.Unix()
	_, err := db.Exec("INSERT INTO sessions (session_id, session_user_id, session_user_name, session_time_start, session_time_last, session_ip, session_browser, session_forwarded_for) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)", sessionId, userId, userName, sessionTimeStart, sessionTimeStart, ip, browser, forwardedFor)
	if err != nil {
		return Session{}, fmt.Errorf("Error while inserting session id '%s' for user id %d: %s", sessionId, userId, err)
	}
	session := Session{
		SessionId:           sessionId,
		SessionUserId:       userId,
		SessionUserName:     userName,
		SessionTimeStart:    sessionTimeStart,
		SessionTimeLast:     sessionTimeStart,
		SessionIp:           ip,
		SessionBrowser:      browser,
		SessionForwardedFor: forwardedFor,
	}
	return session, nil
}

func UpdateSessionTimeLast(ctx context.Context, sessionId string) error {
	db := OpenDb(ctx, "sessions")
	defer db.Close()
	now := time.Now().UTC()
	sessionTimeLast := now.Unix()
	_, err := db.Exec("UPDATE sessions SET session_time_last = $1 WHERE session_id = $2", sessionTimeLast, sessionId)
	if err != nil {
		return fmt.Errorf("Error while updating Session Time Last for session id '%s': %s", sessionId, err)
	}
	return nil
}

func GetSession(ctx context.Context, sessionId string) (Session, error) {
	db := OpenDb(ctx, "sessions")
	defer db.Close()
	var session Session
	err := db.
		QueryRow("SELECT session_id, session_user_id, session_user_name, session_time_start, session_time_last, session_ip, session_browser, session_forwarded_for FROM sessions WHERE session_id = $1", sessionId).
		Scan(&session.SessionId, &session.SessionUserId, &session.SessionUserName, &session.SessionTimeStart, &session.SessionTimeLast, &session.SessionIp, &session.SessionBrowser, &session.SessionForwardedFor)
	if err != nil {
		if err == sql.ErrNoRows {
			// No result found
			return Session{}, fmt.Errorf("Error while retrieving session id '%s' on sessions table: %s: No result found", sessionId, DB_ERROR_NO_RESULT)
		}
		return Session{}, fmt.Errorf("Error while scanning row on sessions table for session id '%s': %s", sessionId, err)
	}
	return session, nil
}
