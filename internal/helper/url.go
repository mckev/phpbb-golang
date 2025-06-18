package helper

import (
	"net"
	"net/http"
	"net/url"
)

const (
	NO_SID = ""
)

func UrlWithSID(rawUrl string, sessionId string) string {
	u, _ := url.Parse(rawUrl)
	q := u.Query()
	if sessionId == NO_SID {
		q.Del("sid")
	} else {
		q.Set("sid", sessionId)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func ExtractUserFingerprint(r *http.Request) (string, string, string) {
	// To prevent hijacking of another user's session
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		ip = r.RemoteAddr
	}
	browser := r.Header.Get("User-Agent")
	forwardedFor := r.Header.Get("X-Forwarded-For")
	return ip, browser, forwardedFor
}
