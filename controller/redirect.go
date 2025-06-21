package controller

import (
	"html/template"
	"net/http"

	"phpbb-golang/internal/logger"
)

func RedirectPage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/redirect?url=https%3A%2F%2Fwww.google.com%2Fsearch%3Fq%3Dhow%2Bto%2Bmake%2Ba%2Braspberry%2Bpi%2Bweb%2Bserver%26hl%3Den%26source%3Dhp%26ei%3Dabcdef
	// How to encode:  encoded := url.QueryEscape("https://www.google.com/search?q=how+to+make+a+raspberry+pi+web+server&hl=en&source=hp&ei=abcdef")
	// How encode works:  It replaces following characters  : %3A, / %2F, ? %3F, = %3D, & %26, + %2B
	ctx := r.Context()
	queryParams := r.URL.Query()
	redirectURIForRedirectPage := queryParams.Get("url")
	if redirectURIForRedirectPage == "" {
		redirectURIForRedirectPage = "/"
	}
	type RedirectPageData struct {
		RedirectURIForRedirectPage string
	}
	redirectPageData := RedirectPageData{
		RedirectURIForRedirectPage: redirectURIForRedirectPage,
	}
	templateOutput, err := template.ParseFiles("./view/templates/redirect.html")
	if err != nil {
		logger.Errorf(ctx, "Error while parsing redirect template file: %s", err)
		return
	}
	err = templateOutput.Execute(w, redirectPageData)
	if err != nil {
		logger.Errorf(ctx, "Error while executing redirect template: %s", err)
		return
	}
}
