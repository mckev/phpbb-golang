package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

const (
	SESSION_KEY = "session"
)

func SessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		queryParams := r.URL.Query()
		sessionId := queryParams.Get("sid")
		session := model.Session{}
		if sessionId != "" {
			// Resume user session (for between pages)
			ip, browser, forwardedFor := helper.ExtractUserFingerprint(r)
			var err error
			session, err = resumeSession(ctx, sessionId, ip, browser, forwardedFor)
			if err != nil {
				logger.Debugf(ctx, "Error while resuming user session for session id '%s': %s", sessionId, err)
				session = model.Session{}
			}
		}
		ctx = context.WithValue(ctx, SESSION_KEY, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getSession(r *http.Request) model.Session {
	ctx := r.Context()
	session, ok := ctx.Value(SESSION_KEY).(model.Session)
	if !ok {
		return model.Session{}
	}
	return session
}

func createSession(ctx context.Context, userId int, userName string, ip string, browser string, forwardedFor string) (model.Session, error) {
	// Create user session (for user registration and user login)
	err := model.CheckIfGuestUser(ctx, userId, userName)
	if err != nil {
		return model.Session{}, fmt.Errorf("Error while creating user session: %s", err)
	}
	sessionId, err := helper.GenerateSessionId()
	if err != nil {
		return model.Session{}, fmt.Errorf("Error while generating a random SID: %s", err)
	}
	return model.CreateSession(ctx, sessionId, userId, userName, ip, browser, forwardedFor)
}

func resumeSession(ctx context.Context, sessionId string, ip string, browser string, forwardedFor string) (model.Session, error) {
	// Existing session is valid if:
	//   - Session id string is valid
	//   - Session exists on sessions table
	//   - 0 <= CurrentTime-SessionTimeLast < SESSION_TIMEOUT_IN_SECONDS
	//   - IP, Browser and ForwardedFor matched.
	if !helper.IsSessionIdValid(sessionId) {
		return model.Session{}, fmt.Errorf("Error while resuming user session: Session id '%s' is not valid", sessionId)
	}
	session, err := model.GetSession(ctx, sessionId)
	if err != nil {
		return model.Session{}, fmt.Errorf("Error while resuming user session: %s", err)
	}
	now := time.Now().UTC()
	currentTime := now.Unix()
	if currentTime-session.SessionTimeLast < 0 || currentTime-session.SessionTimeLast >= model.SESSION_TIMEOUT_IN_SECONDS {
		return model.Session{}, fmt.Errorf("Error while resuming user session: Session has timed out with delta time %d seconds", currentTime-session.SessionTimeLast)
	}
	if session.SessionIp != ip || session.SessionBrowser != browser || session.SessionForwardedFor != forwardedFor {
		return model.Session{}, fmt.Errorf("Error while resuming user session: User fingerprint does not match: IP %s, Browser %s, ForwardedFor %s do not equal values in database IP %s, Browser %s, ForwardedFor %s", ip, browser, forwardedFor, session.SessionIp, session.SessionBrowser, session.SessionForwardedFor)
	}
	err = model.UpdateSessionTimeLast(ctx, sessionId)
	if err != nil {
		return model.Session{}, fmt.Errorf("Error while resuming user session: %s", err)
	}
	return session, nil
}
