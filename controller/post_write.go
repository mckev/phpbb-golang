package controller

import (
	"html/template"
	"net/http"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func PostWritePage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/post_write?t=2
	ctx := r.Context()
	session := getSession(r)

	queryParams := r.URL.Query()
	topicId := helper.StrToInt(queryParams.Get("t"), model.INVALID_TOPIC_ID)

	// Prepare template files
	templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/post_write.html")
	if err != nil {
		logger.Errorf(ctx, "Error while parsing post write template files: %s", err)
		return
	}

	// Prepare data
	topic, err := model.GetTopic(ctx, topicId)
	if err != nil {
		logger.Errorf(ctx, "Error while getting topic id %d: %s", topicId, err)
		return
	}
	startItem := forumhelper.ComputeStartItem(topic.TopicNumPosts, topic.TopicNumPosts, model.MAX_POSTS_PER_PAGE)
	posts, err := model.ListPosts(ctx, topicId, startItem)
	if err != nil {
		logger.Errorf(ctx, "Error while listing posts: %s", err)
		return
	}
	users, err := model.ListUsersOfTopic(ctx, topicId)
	usersMap := map[int]model.User{} // Convert users from a list into a map
	for _, user := range users {
		usersMap[user.UserId] = user
	}
	forums, err := model.ListForums(ctx)
	if err != nil {
		logger.Errorf(ctx, "Error while listing forums: %s", err)
		return
	}
	forumNavTrails := forumhelper.ComputeForumNavTrails(ctx, forums, topic.ForumId)
	type PostWritePageData struct {
		Topic                   model.Topic
		Posts                   []model.Post
		UsersMap                map[int]model.User
		Session                 model.Session
		RedirectURIForLoginPage string
		ForumNavTrails          []forumhelper.ForumNavTrail
	}
	postWritePageData := PostWritePageData{
		Topic:                   topic,
		Posts:                   posts,
		UsersMap:                usersMap,
		Session:                 session,
		RedirectURIForLoginPage: helper.UrlWithSID(r.URL.RequestURI(), ""),
		ForumNavTrails:          forumNavTrails,
	}

	// Execute template
	err = templateOutput.ExecuteTemplate(w, "overall", postWritePageData)
	if err != nil {
		logger.Errorf(ctx, "Error while executing post write template: %s", err)
		return
	}
}
