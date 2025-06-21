package controller

import (
	"html/template"
	"net/http"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func ForumsPage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/forums?f=1
	ctx := r.Context()
	session := GetSession(r)

	queryParams := r.URL.Query()
	forumId := helper.StrToInt(queryParams.Get("f"), model.INVALID_FORUM_ID)
	if forumId == model.INVALID_FORUM_ID {
		http.Redirect(w, r, "./", http.StatusFound)
		return
	}

	// Prepare template files
	templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/forums.html")
	if err != nil {
		logger.Errorf(ctx, "Error while parsing forums template files: %s", err)
		return
	}

	// Prepare data
	forum, err := model.GetForum(ctx, forumId)
	if err != nil {
		logger.Errorf(ctx, "Error while getting forum id %d: %s", forumId, err)
		return
	}
	forums, err := model.ListForums(ctx)
	if err != nil {
		logger.Errorf(ctx, "Error while listing forums: %s", err)
		return
	}
	forumChildNodes := forumhelper.ComputeForumChildNodes(ctx, forums, forumId, 0)
	forumNavTrails := forumhelper.ComputeForumNavTrails(ctx, forums, forumId)
	type ForumsPageData struct {
		Forum           model.Forum
		ForumChildNodes []forumhelper.ForumNode
		Session         model.Session
		// RedirectURIForLoginPage is useful for returning the user to their original page after they click "Login" and successfully authenticate.
		// And since the Session Id is provided by the login page, there's no need to include it in the RedirectURIForLoginPage.
		RedirectURIForLoginPage string
		ForumNavTrails          []forumhelper.ForumNavTrail
	}
	forumsPageData := ForumsPageData{
		Forum:                   forum,
		ForumChildNodes:         forumChildNodes,
		Session:                 session,
		RedirectURIForLoginPage: helper.UrlWithSID(r.URL.RequestURI(), ""),
		ForumNavTrails:          forumNavTrails,
	}

	// Execute template
	err = templateOutput.ExecuteTemplate(w, "overall", forumsPageData)
	if err != nil {
		logger.Errorf(ctx, "Error while executing forums template: %s", err)
		return
	}
}
