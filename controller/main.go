package controller

import (
	"html/template"
	"net/http"
	"time"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func MainPage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/
	ctx := r.Context()
	session := GetSession(r)

	// Prepare template files
	templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/main.html")
	if err != nil {
		logger.Errorf(ctx, "Error while parsing forums template files: %s", err)
		return
	}

	// Prepare data
	now := time.Now().UTC()
	currentTime := now.Unix()
	forums, err := model.ListForums(ctx)
	if err != nil {
		logger.Errorf(ctx, "Error while listing forums: %s", err)
		return
	}
	forumChildNodes := forumhelper.ComputeForumChildNodes(ctx, forums, model.ROOT_FORUM_ID, 0)
	type MainPageData struct {
		CurrentTime             int64
		ForumChildNodes         []forumhelper.ForumNode
		Session                 model.Session
		RedirectURIForLoginPage string
		ForumNavTrails          []forumhelper.ForumNavTrail
	}
	forumsPageData := MainPageData{
		CurrentTime:             currentTime,
		ForumChildNodes:         forumChildNodes,
		Session:                 session,
		RedirectURIForLoginPage: "./",
		ForumNavTrails:          []forumhelper.ForumNavTrail{},
	}

	// Execute template
	err = templateOutput.ExecuteTemplate(w, "overall", forumsPageData)
	if err != nil {
		logger.Errorf(ctx, "Error while executing forums template: %s", err)
		return
	}
}
