package controller

import (
	"context"
	"net/http"

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
			session, err = model.ResumeSession(ctx, sessionId, ip, browser, forwardedFor)
			if err != nil {
				logger.Debugf(ctx, "Error while resuming user session for session id '%s': %s", sessionId, err)
				session = model.Session{}
			}
		}
		ctx = context.WithValue(ctx, SESSION_KEY, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetSession(r *http.Request) model.Session {
	ctx := r.Context()
	session, ok := ctx.Value(SESSION_KEY).(model.Session)
	if !ok {
		return model.Session{}
	}
	return session
}
