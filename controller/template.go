package controller

import (
	"html/template"

	"phpbb-golang/internal/bbcode"
	"phpbb-golang/internal/helper"
)

var funcMap = template.FuncMap{
	"fnAdd": func(x, y int) int {
		return x + y
	},
	"fnMod": func(x, y int) bool {
		return x%y == 0
	},
	"fnUnixTimeToStr": func(unixTime int64) string {
		return helper.UnixTimeToStr(unixTime)
	},
	"fnUrlWithSID": func(rawUrl string, sessionId string) string {
		return helper.UrlWithSID(rawUrl, sessionId)
	},
	"fnBbcodeToHtml": func(bbcodeStr string) template.HTML {
		// To print raw, unescaped HTML within a Go HTML template, the html/template package provides the template.HTML type. By converting a string containing HTML to template.HTML, you can instruct the template engine to render it as raw HTML instead of escaping it for safe output.
		// WARNING: Since this Go template function outputs raw HTML, make sure it is safe from attacks such as XSS.
		return template.HTML(bbcode.ConvertBbcodeToHtml(bbcodeStr))
	},
}
