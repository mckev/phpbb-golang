package controller

import (
	"html/template"
	"net/http"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func PostsPage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/posts?t=2
	ctx := r.Context()
	session := GetSession(r)

	// Parse query string. We use queryParams.Get("key") to retrieve the first value for a given query parameter.
	queryParams := r.URL.Query()
	topicId := helper.StrToInt(queryParams.Get("t"), model.INVALID_TOPIC_ID)
	startItem := helper.StrToInt(queryParams.Get("start"), 0)

	// Case specify post id
	postId := helper.StrToInt(queryParams.Get("p"), model.INVALID_POST_ID)
	if postId > 0 {
		post, err := model.GetPost(ctx, postId)
		if err != nil {
			logger.Errorf(ctx, "Error while getting post id %d: %s", postId, err)
			return
		}
		topicId = post.TopicId
		startItem, err = model.CountPostCurItem(ctx, topicId, postId)
		if err != nil {
			logger.Errorf(ctx, "Error while counting current item: %s", err)
			return
		}
	}

	// Prepare template files
	templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/posts.html")
	if err != nil {
		logger.Errorf(ctx, "Error while parsing posts template files: %s", err)
		return
	}

	// Prepare data
	topic, err := model.GetTopic(ctx, topicId)
	if err != nil {
		logger.Errorf(ctx, "Error while getting topic id %d: %s", topicId, err)
		return
	}
	forum, err := model.GetForum(ctx, topic.ForumId)
	if err != nil {
		logger.Errorf(ctx, "Error while getting forum id %d: %s", topic.ForumId, err)
		return
	}
	startItem = forumhelper.ComputeStartItem(startItem, topic.TopicNumPosts, model.MAX_POSTS_PER_PAGE)
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
	forumNavTrails := forumhelper.ComputeForumNavTrails(ctx, forums, forum.ForumId)
	paginations := forumhelper.ComputePaginations(startItem, topic.TopicNumPosts, model.MAX_POSTS_PER_PAGE)
	type PostsPageData struct {
		Forum                   model.Forum
		Topic                   model.Topic
		Posts                   []model.Post
		UsersMap                map[int]model.User
		Paginations             []forumhelper.Pagination
		Session                 model.Session
		RedirectURIForLoginPage string
		ForumNavTrails          []forumhelper.ForumNavTrail
	}
	postsPageData := PostsPageData{
		Forum:                   forum,
		Topic:                   topic,
		Posts:                   posts,
		UsersMap:                usersMap,
		Paginations:             paginations,
		Session:                 session,
		RedirectURIForLoginPage: helper.UrlWithSID(r.URL.RequestURI(), ""),
		ForumNavTrails:          forumNavTrails,
	}

	// Execute template
	// Go HTML Templates:
	//   - Go's html/template package automatically strips HTML comments (<!-- comment -->) during template execution.
	//   - Go templates do not inherently support nested template definitions in the way one might expect from other templating engines. While you can define templates within other templates using {{define}}, they are all effectively "hoisted" to the top level and treated as independent templates within a single namespace. This means you can't directly access a nested template as a property of its parent.
	//     To achieve a similar effect as nested templates, you should define each template separately and then compose them using the {{template}} action. When calling {{template}}, you can pass data to the included template using a dot ., which represents the current data context.
	//     If you need to access data from a parent template, you can either pass it explicitly to the nested template or use variables to store the data before calling the nested template.
	//   - In Go templates, . refers to the current context, which changes within a range loop to the current element being iterated over.
	//     To access the outer struct's fields from within the inner loop, $ is used, which always refers to the root context (the original data passed to the template).
	err = templateOutput.ExecuteTemplate(w, "overall", postsPageData)
	if err != nil {
		logger.Errorf(ctx, "Error while executing posts template: %s", err)
		return
	}
}
