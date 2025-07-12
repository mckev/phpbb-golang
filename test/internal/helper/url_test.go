package helper

import (
	"bytes"
	"html/template"
	"net/url"
	"testing"

	"phpbb-golang/internal/helper"
	"phpbb-golang/model"
)

func TestUrlWithSID(t *testing.T) {
	{
		actual := helper.UrlWithSID("./forums", "123456")
		expected := "./forums?sid=123456"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./forums?f=1", "123456")
		expected := "./forums?f=1&sid=123456"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./posts?p=101#p101", "123456")
		expected := "./posts?p=101&sid=123456#p101"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&start=100", "123456")
		expected := "./topics?f=1&sid=123456&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&sid=old&start=100", "123456")
		expected := "./topics?f=1&sid=123456&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestUrlWithSID_NoSID(t *testing.T) {
	{
		actual := helper.UrlWithSID("./forums", helper.NO_SID)
		expected := "./forums"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./forums?f=1", helper.NO_SID)
		expected := "./forums?f=1"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./posts?p=101#p101", helper.NO_SID)
		expected := "./posts?p=101#p101"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&start=100", helper.NO_SID)
		expected := "./topics?f=1&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./topics?f=1&sid=old&start=100", helper.NO_SID)
		expected := "./topics?f=1&start=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestUrlWithSID_Xss(t *testing.T) {
	{
		actual := helper.UrlWithSID("./top<script>ics?f=<script>1&start<script>=100", "123456")
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&sid=123456&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./top%3Cscript%3Eics?f=%3Cscript%3E1&start%3Cscript%3E=100", "123456")
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&sid=123456&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
	{
		actual := helper.UrlWithSID("./top%3Cscript%3Eics?sid=old&f=%3Cscript%3E1&start%3Cscript%3E=100", helper.NO_SID)
		expected := "./top%3Cscript%3Eics?f=%3Cscript%3E1&start%3Cscript%3E=100"
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}

func TestUrlWithSID_SimulateRedirectionForLoginPage(t *testing.T) {
	uri := "/post_write?mode=quote&p=6"
	type TestPageData struct {
		RedirectURIForLoginPage string
		Session                 model.Session
	}
	testPageData := TestPageData{
		RedirectURIForLoginPage: url.QueryEscape(helper.UrlWithSID(uri, "")),
		Session:                 model.Session{},
	}
	// Use Go HTML Template to simulate rendered HTML
	var funcMap = template.FuncMap{
		"fnUrlWithSID": func(rawUrl string, sessionId string) string {
			return helper.UrlWithSID(rawUrl, sessionId)
		},
	}
	const templateString = `<a href='{{ fnUrlWithSID (printf "./user_login?redirect=%s" .RedirectURIForLoginPage) .Session.SessionId }}' title="Login" accesskey="x" role="menuitem">`
	templateOutput, err := template.New("").Funcs(funcMap).Parse(templateString)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	var buffer bytes.Buffer
	err = templateOutput.Execute(&buffer, testPageData)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	htmlStr := buffer.String()
	{
		actual := htmlStr
		expected := `<a href='./user_login?redirect=%2Fpost_write%3Fmode%3Dquote%26p%3D6' title="Login" accesskey="x" role="menuitem">`
		if actual != expected {
			t.Errorf("Got %s, wanted %s", actual, expected)
			return
		}
	}
}
