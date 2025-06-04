package helper

import "net/url"

func UrlWithSID(rawUrl string, sessionId string) string {
	u, _ := url.Parse(rawUrl)
	q := u.Query()
	q.Set("sid", sessionId)
	u.RawQuery = q.Encode()
	return u.String()
}
