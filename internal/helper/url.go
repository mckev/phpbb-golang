package helper

import "net/url"

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
