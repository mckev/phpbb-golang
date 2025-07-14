package controller

import (
	"html/template"
	"net/http"
	"net/url"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func TopicsPage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/topics?f=10
	ctx := r.Context()
	session := getSession(r)

	queryParams := r.URL.Query()
	forumId := helper.StrToInt(queryParams.Get("f"), model.INVALID_FORUM_ID)
	startItem := helper.StrToInt(queryParams.Get("start"), 0)

	// Prepare data
	forum, err := model.GetForum(ctx, forumId)
	if err != nil {
		logger.Errorf(ctx, "Error while getting forum id %d: %s", forumId, err)
		return
	}
	startItem = forumhelper.ComputeStartItem(startItem, forum.ForumNumTopics, model.MAX_TOPICS_PER_PAGE)
	topics, err := model.ListTopics(ctx, forumId, startItem)
	if err != nil {
		logger.Errorf(ctx, "Error while listing topics: %s", err)
		return
	}
	type TopicWithInfo struct {
		Topic           model.Topic
		PostPaginations []forumhelper.Pagination
	}
	topicsWithInfo := []TopicWithInfo{}
	for _, topic := range topics {
		postPaginations := forumhelper.ComputePaginations(max(topic.TopicNumPosts-1, 0), topic.TopicNumPosts, model.MAX_POSTS_PER_PAGE)
		topicsWithInfo = append(topicsWithInfo, TopicWithInfo{
			Topic:           topic,
			PostPaginations: postPaginations,
		})
	}
	forums, err := model.ListForums(ctx)
	if err != nil {
		logger.Errorf(ctx, "Error while listing forums: %s", err)
		return
	}
	forumNavTrails := forumhelper.ComputeForumNavTrails(ctx, forums, forumId)
	topicPaginations := forumhelper.ComputePaginations(startItem, forum.ForumNumTopics, model.MAX_TOPICS_PER_PAGE)
	type TopicsPageData struct {
		Forum                   model.Forum
		TopicsWithInfo          []TopicWithInfo
		TopicPaginations        []forumhelper.Pagination
		Session                 model.Session
		RedirectURIForLoginPage string
		ForumNavTrails          []forumhelper.ForumNavTrail
	}
	topicsPageData := TopicsPageData{
		Forum:                   forum,
		TopicsWithInfo:          topicsWithInfo,
		TopicPaginations:        topicPaginations,
		Session:                 session,
		RedirectURIForLoginPage: url.QueryEscape(helper.UrlWithSID(r.URL.RequestURI(), helper.NO_SID)),
		ForumNavTrails:          forumNavTrails,
	}

	// Execute template
	templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/topics.html")
	if err != nil {
		logger.Errorf(ctx, "Error while parsing topics template files: %s", err)
		return
	}
	err = templateOutput.ExecuteTemplate(w, "overall", topicsPageData)
	if err != nil {
		logger.Errorf(ctx, "Error while executing topics template: %s", err)
		return
	}
}
