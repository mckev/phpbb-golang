package main

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"phpbb-golang/examples/myforum"
	"phpbb-golang/internal/bbcode"
	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpMethod := r.Method
	urlPath := filepath.Clean(r.URL.Path)
	queryParams := r.URL.Query()
	logger.Debugf(ctx, "%s %s", httpMethod, urlPath)

	// Template Functions
	funcMap := template.FuncMap{
		"fnAdd": func(x, y int) int {
			return x + y
		},
		"fnUnixTimeToStr": func(unixTime int64) string {
			return helper.UnixTimeToStr(unixTime)
		},
		"fnBbcodeToHtml": func(bbcodeStr string) template.HTML {
			// To print raw, unescaped HTML within a Go HTML template, the html/template package provides the template.HTML type. By converting a string containing HTML to template.HTML, you can instruct the template engine to render it as raw HTML instead of escaping it for safe output.
			// WARNING: Since this Go template function outputs raw HTML, make sure it is safe from attacks such as XSS.
			return template.HTML(bbcode.ConvertBbcodeToHtml(bbcodeStr))
		},
	}

	if urlPath == "/" {
		io.WriteString(w, "Welcome to Golang BB!")

	} else if urlPath == "/redirect" {
		// To try: http://localhost:9000/redirect?url=https%3A%2F%2Fwww.google.com%2Fsearch%3Fq%3Dhow%2Bto%2Bmake%2Ba%2Braspberry%2Bpi%2Bweb%2Bserver%26hl%3Den%26source%3Dhp%26ei%3Dabcdef
		// How to encode:  encoded := url.QueryEscape(`https://www.google.com/search?q=how+to+make+a+raspberry+pi+web+server&hl=en&source=hp&ei=abcdef`)
		// How encode works:  It replaces following characters  : %3A, / %2F, ? %3F, = %3D, & %26, + %2B
		redirectUrl := queryParams.Get("url")
		if redirectUrl == "" {
			redirectUrl = "/"
		}
		type RedirectPageData struct {
			RedirectUrl string
		}
		redirectPageData := RedirectPageData{
			RedirectUrl: redirectUrl,
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

	} else if urlPath == "/myforum/main" {
		// To try: http://localhost:9000/myforum/main
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/main.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}
	} else if urlPath == "/myforum/forums" {
		// To try: http://localhost:9000/myforum/forums
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/forums.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}
	} else if urlPath == "/myforum/topics" {
		// To try: http://localhost:9000/myforum/topics
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/topics.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}
	} else if urlPath == "/myforum/posts" {
		// To try: http://localhost:9000/myforum/posts
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/posts.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}

	} else if urlPath == "/topics" {
		// To try: http://localhost:9000/topics?f=10
		forumId := helper.StrToInt(queryParams.Get("f"), model.INVALID_FORUM_ID)
		startItem := helper.StrToInt(queryParams.Get("start"), 0)

		// Prepare template files
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/topics.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing topics template files: %s", err)
			return
		}

		// Prepare data
		forum, err := model.GetForum(ctx, forumId)
		if err != nil {
			logger.Errorf(ctx, "Error while getting forum: %s", err)
		}
		startItem = forumhelper.FixStartItem(startItem, forum.ForumNumTopics, model.MAX_TOPICS_PER_PAGE)
		topics, err := model.ListTopics(ctx, forumId, startItem)
		if err != nil {
			logger.Errorf(ctx, "Error while listing topics: %s", err)
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
		forumNavTrails, err := model.ComputeForumNavTrails(ctx, forumId)
		if err != nil {
			logger.Errorf(ctx, "Error while computing Forum Nav Trails for forum id %d: %s", forumId, err)
		}
		topicPaginations := forumhelper.ComputePaginations(startItem, forum.ForumNumTopics, model.MAX_TOPICS_PER_PAGE)
		type TopicsPageData struct {
			Forum            model.Forum
			TopicsWithInfo   []TopicWithInfo
			ForumNavTrails   []model.ForumNavTrail
			TopicPaginations []forumhelper.Pagination
		}
		topicsPageData := TopicsPageData{
			Forum:            forum,
			TopicsWithInfo:   topicsWithInfo,
			ForumNavTrails:   forumNavTrails,
			TopicPaginations: topicPaginations,
		}

		// Execute template
		err = templateOutput.ExecuteTemplate(w, "overall", topicsPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing topics template: %s", err)
			return
		}

	} else if urlPath == "/posts" {
		// To try: http://localhost:9000/posts?t=2
		// Parse query string. We use queryParams.Get("key") to retrieve the first value for a given query parameter.
		topicId := helper.StrToInt(queryParams.Get("t"), model.INVALID_TOPIC_ID)
		startItem := helper.StrToInt(queryParams.Get("start"), 0)

		// Prepare template files
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/posts.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing posts template files: %s", err)
			return
		}

		// Prepare data
		topic, err := model.GetTopic(ctx, topicId)
		if err != nil {
			logger.Errorf(ctx, "Error while getting topic: %s", err)
		}
		forum, err := model.GetForum(ctx, topic.ForumId)
		if err != nil {
			logger.Errorf(ctx, "Error while getting forum: %s", err)
		}
		startItem = forumhelper.FixStartItem(startItem, topic.TopicNumPosts, model.MAX_POSTS_PER_PAGE)
		posts, err := model.ListPosts(ctx, topicId, startItem)
		if err != nil {
			logger.Errorf(ctx, "Error while listing posts: %s", err)
			return
		}
		users, err := model.ListUsers(ctx, topicId)
		usersMap := map[int]model.User{} // Convert users from a list into a map
		for _, user := range users {
			usersMap[user.UserId] = user
		}
		forumNavTrails, err := model.ComputeForumNavTrails(ctx, forum.ForumId)
		if err != nil {
			logger.Errorf(ctx, "Error while computing Forum Nav Trails for forum id %d: %s", forum.ForumId, err)
		}
		paginations := forumhelper.ComputePaginations(startItem, topic.TopicNumPosts, model.MAX_POSTS_PER_PAGE)
		type PostsPageData struct {
			Forum          model.Forum
			Topic          model.Topic
			Posts          []model.Post
			UsersMap       map[int]model.User
			ForumNavTrails []model.ForumNavTrail
			Paginations    []forumhelper.Pagination
		}
		postsPageData := PostsPageData{
			Forum:          forum,
			Topic:          topic,
			Posts:          posts,
			UsersMap:       usersMap,
			ForumNavTrails: forumNavTrails,
			Paginations:    paginations,
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
	} else {
		logger.Errorf(ctx, "URL Path not supported: %s %s", httpMethod, urlPath)
		return
	}
}

func main() {
	ctx := context.Background()

	// MyForum example
	myforum.InitMyforum(ctx)
	myforum.DebugMyforum(ctx)

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./view/static/assets/"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./view/static/images/"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./view/static/styles/"))))
	http.HandleFunc("/", serveTemplate)
	portNumber := 9000
	logger.Infof(ctx, "Server is listening on port %d", portNumber)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", portNumber), nil)
	if err != nil {
		logger.Fatalf(ctx, "Error while listening on port %d: %s", portNumber, err)
	}
}
