package main

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
	"time"

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

	// Session
	sessionId := queryParams.Get("sid")
	session := model.Session{}
	if sessionId != "" {
		// Resume user session (for between pages)
		ip, browser, forwardedFor := helper.ExtractUserFingerprint(r)
		var err error
		session, err = model.ResumeSession(ctx, sessionId, ip, browser, forwardedFor)
		if err != nil {
			logger.Debugf(ctx, "Error while resuming user session for session id '%s': %s", sessionId, err)
			session = model.Session{}
			// Falls through
		}
	}

	// Template Functions
	funcMap := template.FuncMap{
		"fnAdd": func(x, y int) int {
			return x + y
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

	if urlPath == "/" {
		// To try: http://localhost:9000/

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
			CurrentTime     int64
			ForumChildNodes []forumhelper.ForumNode
			RedirectURI     string
			SessionId       string
			ForumNavTrails  []forumhelper.ForumNavTrail
		}
		forumsPageData := MainPageData{
			CurrentTime:     currentTime,
			ForumChildNodes: forumChildNodes,
			RedirectURI:     "./",
			SessionId:       session.SessionId,
			ForumNavTrails:  []forumhelper.ForumNavTrail{},
		}

		// Execute template
		err = templateOutput.ExecuteTemplate(w, "overall", forumsPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing forums template: %s", err)
			return
		}

	} else if urlPath == "/redirect" {
		// To try: http://localhost:9000/redirect?url=https%3A%2F%2Fwww.google.com%2Fsearch%3Fq%3Dhow%2Bto%2Bmake%2Ba%2Braspberry%2Bpi%2Bweb%2Bserver%26hl%3Den%26source%3Dhp%26ei%3Dabcdef
		// How to encode:  encoded := url.QueryEscape("https://www.google.com/search?q=how+to+make+a+raspberry+pi+web+server&hl=en&source=hp&ei=abcdef")
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
	} else if urlPath == "/myforum/user_login" {
		// To try: http://localhost:9000/myforum/user_login
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/user_login.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}
	} else if urlPath == "/myforum/user_register" {
		// To try: http://localhost:9000/myforum/user_register
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/user_register.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}
	} else if urlPath == "/myforum/user_register_created" {
		// To try: http://localhost:9000/myforum/user_register_created
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/user_register_created.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}
	} else if urlPath == "/myforum/user_register_activated" {
		// To try: http://localhost:9000/myforum/user_register_activated
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/user_register_activated.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}

	} else if urlPath == "/forums" {
		// To try: http://localhost:9000/forums?f=1
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
			// RedirectURI is useful for returning the user to their original page after they click "Login" and successfully authenticate.
			// And since the Session Id is provided by the login page, there's no need to include it in the RedirectURI.
			RedirectURI    string
			SessionId      string
			ForumNavTrails []forumhelper.ForumNavTrail
		}
		forumsPageData := ForumsPageData{
			Forum:           forum,
			ForumChildNodes: forumChildNodes,
			RedirectURI:     helper.UrlWithSID(r.URL.RequestURI(), ""),
			SessionId:       session.SessionId,
			ForumNavTrails:  forumNavTrails,
		}

		// Execute template
		err = templateOutput.ExecuteTemplate(w, "overall", forumsPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing forums template: %s", err)
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
			Forum            model.Forum
			TopicsWithInfo   []TopicWithInfo
			TopicPaginations []forumhelper.Pagination
			RedirectURI      string
			SessionId        string
			ForumNavTrails   []forumhelper.ForumNavTrail
		}
		topicsPageData := TopicsPageData{
			Forum:            forum,
			TopicsWithInfo:   topicsWithInfo,
			TopicPaginations: topicPaginations,
			RedirectURI:      helper.UrlWithSID(r.URL.RequestURI(), ""),
			SessionId:        session.SessionId,
			ForumNavTrails:   forumNavTrails,
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
			Forum          model.Forum
			Topic          model.Topic
			Posts          []model.Post
			UsersMap       map[int]model.User
			Paginations    []forumhelper.Pagination
			RedirectURI    string
			SessionId      string
			ForumNavTrails []forumhelper.ForumNavTrail
		}
		postsPageData := PostsPageData{
			Forum:          forum,
			Topic:          topic,
			Posts:          posts,
			UsersMap:       usersMap,
			Paginations:    paginations,
			RedirectURI:    helper.UrlWithSID(r.URL.RequestURI(), ""),
			SessionId:      session.SessionId,
			ForumNavTrails: forumNavTrails,
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

	} else if urlPath == "/user_register" {
		// To try: http://localhost:9000/user_register
		type FormData struct {
			Username        string
			NewPassword     string
			PasswordConfirm string
			Email           string
			Errors          []string
		}
		formData := FormData{}
		if httpMethod == "POST" {
			err := r.ParseForm()
			if err != nil {
				logger.Errorf(ctx, "Error while parsing form upon user registration: %s", err)
				return
			}
			formData.Username = strings.TrimSpace(r.Form.Get("username"))
			if len(formData.Username) < 4 {
				formData.Errors = append(formData.Errors, "The username you entered is too short.")
			} else if len(formData.Username) > 20 {
				formData.Errors = append(formData.Errors, "The username you entered is too long.")
			}
			// TODO: Check that username does not have invalid characters, such as space
			formData.NewPassword = strings.TrimSpace(r.Form.Get("new_password"))
			if len(formData.NewPassword) < 8 {
				formData.Errors = append(formData.Errors, "The password you entered is too short.")
			}
			if !helper.IsPasswordValid(formData.NewPassword) {
				formData.Errors = append(formData.Errors, "Password must be at least 8 characters long, must contain letters in mixed case and must contain numbers.")
			}
			formData.PasswordConfirm = strings.TrimSpace(r.Form.Get("password_confirm"))
			if formData.PasswordConfirm != formData.NewPassword {
				formData.Errors = append(formData.Errors, "Password and confirmation do not match.")
			}
			formData.Email = strings.TrimSpace(r.Form.Get("email"))
			if formData.Email != "" && !helper.IsEmailValid(formData.Email) {
				formData.Errors = append(formData.Errors, "The email address format is invalid.")
			}
			userId := model.GUEST_USER_ID
			if len(formData.Errors) == 0 {
				// Insert user into database
				userId, err = model.InsertUser(ctx, formData.Username, formData.NewPassword, formData.Email, "")
				if err != nil {
					if strings.Contains(err.Error(), model.DB_ERROR_UNIQUE_CONSTRAINT) {
						formData.Errors = append(formData.Errors, "This username is already taken. Please choose a different one.")
					} else {
						logger.Errorf(ctx, "Error while inserting user: %s", err)
						formData.Errors = append(formData.Errors, "The system is currently experiencing an issue. Please try again later.")
					}
				}
			}
			if len(formData.Errors) == 0 {
				// Validation successful
				// fmt.Fprintf(w, "Form submitted successfully!\n")
				// fmt.Fprintf(w, "Username: %s\n", formData.Username)
				// fmt.Fprintf(w, "Password: %s\n", formData.NewPassword)
				// fmt.Fprintf(w, "Confirm password: %s\n", formData.PasswordConfirm)
				// fmt.Fprintf(w, "Email address: %s\n", formData.Email)
				// TODO: Handle CSRF token validation

				// Create user session (for user registration and user login)
				ip, browser, forwardedFor := helper.ExtractUserFingerprint(r)
				session, err = model.CreateSession(ctx, userId, ip, browser, forwardedFor)
				if err != nil {
					logger.Errorf(ctx, "Error while creating user session: %s", err)
					return
				}
				err = model.UpdateLastVisitTimeForUser(ctx, userId)
				if err != nil {
					logger.Errorf(ctx, "Error while updating last visit time for user id %d: %s", userId, err)
					return
				}

				templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/user_register_created.html")
				if err != nil {
					logger.Errorf(ctx, "Error while parsing user register created template files: %s", err)
					return
				}
				type UserRegisterPageData struct {
					RedirectURI    string
					SessionId      string
					ForumNavTrails []forumhelper.ForumNavTrail
				}
				userRegisterPageData := UserRegisterPageData{
					RedirectURI:    "./",
					SessionId:      session.SessionId,
					ForumNavTrails: []forumhelper.ForumNavTrail{},
				}
				err = templateOutput.ExecuteTemplate(w, "overall", userRegisterPageData)
				if err != nil {
					logger.Errorf(ctx, "Error while executing user register created template: %s", err)
					return
				}
				return
			}
		}

		// Prepare template files
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/user_register.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing user register template files: %s", err)
			return
		}

		// Prepare data
		type UserRegisterPageData struct {
			FormData       FormData
			RedirectURI    string
			SessionId      string
			ForumNavTrails []forumhelper.ForumNavTrail
		}
		userRegisterPageData := UserRegisterPageData{
			FormData:       formData,
			RedirectURI:    "./",
			SessionId:      session.SessionId,
			ForumNavTrails: []forumhelper.ForumNavTrail{},
		}

		// Execute template
		err = templateOutput.ExecuteTemplate(w, "overall", userRegisterPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing user register template: %s", err)
			return
		}

	} else if urlPath == "/user_login" {
		// To try: http://localhost:9000/user_login
		type FormData struct {
			Username   string
			Password   string
			Errors     []string
			RedirectTo string
		}
		formData := FormData{}
		if httpMethod == "POST" {
			err := r.ParseForm()
			if err != nil {
				logger.Errorf(ctx, "Error while parsing form upon user login: %s", err)
				return
			}
			formData.Username = strings.TrimSpace(r.Form.Get("username"))
			if formData.Username == "" {
				formData.Errors = append(formData.Errors, "You have specified an incorrect username. Please check your username and try again.")
			} else if len(formData.Username) < 4 {
				formData.Errors = append(formData.Errors, "The username you entered is too short.")
			} else if len(formData.Username) > 20 {
				formData.Errors = append(formData.Errors, "The username you entered is too long.")
			}
			formData.Password = strings.TrimSpace(r.Form.Get("password"))
			if formData.Password == "" {
				formData.Errors = append(formData.Errors, "You cannot login without a password.")
			} else if len(formData.Password) < 8 {
				formData.Errors = append(formData.Errors, "The password you entered is too short.")
			}
			formData.RedirectTo = strings.TrimSpace(r.Form.Get("redirect"))
			if formData.RedirectTo == "" {
				formData.RedirectTo = "./"
			}
			var user model.User
			if len(formData.Errors) == 0 {
				user, err = model.GetUserForLogin(ctx, formData.Username)
				if err != nil {
					logger.Errorf(ctx, "Error while retrieving user name '%s' for login: %s", formData.Username, err)
					formData.Errors = append(formData.Errors, "Invalid username or password")
				}
			}
			if len(formData.Errors) == 0 {
				if !helper.IsPasswordCorrect(formData.Password, user.UserPasswordHashed) {
					logger.Errorf(ctx, "Error while login for user name '%s': Password is incorrect", formData.Username)
					formData.Errors = append(formData.Errors, "Invalid username or password")
				}
			}
			if len(formData.Errors) == 0 {
				// Validation successful
				// fmt.Fprintf(w, "Form submitted successfully!\n")
				// fmt.Fprintf(w, "Username: %s\n", formData.Username)
				// fmt.Fprintf(w, "Password: %s\n", formData.Password)
				// fmt.Fprintf(w, "RedirectTo: %s\n", formData.RedirectTo)
				// TODO: Handle CSRF token validation

				// Create user session (for user registration and user login)
				ip, browser, forwardedFor := helper.ExtractUserFingerprint(r)
				userId := user.UserId
				session, err = model.CreateSession(ctx, userId, ip, browser, forwardedFor)
				if err != nil {
					logger.Errorf(ctx, "Error while creating user session: %s", err)
					return
				}
				err = model.UpdateLastVisitTimeForUser(ctx, userId)
				if err != nil {
					logger.Errorf(ctx, "Error while updating last visit time for user id %d: %s", userId, err)
					return
				}

				// Redirect user to their last visited page
				http.Redirect(w, r, helper.UrlWithSID(formData.RedirectTo, session.SessionId), http.StatusFound)
				return
			}
		}

		// Prepare template files
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/user_login.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing user login template files: %s", err)
			return
		}

		// Prepare data
		if formData.RedirectTo == "" {
			// If say, user entered a wrong password, then we shall keep the hidden input "redirect" intact.
			// Otherwise we record query parameter "redirect" as hidden input.
			if queryParams.Get("redirect") == "" {
				formData.RedirectTo = "./"
			} else {
				formData.RedirectTo = queryParams.Get("redirect")
			}
		}
		type UserLoginPageData struct {
			FormData       FormData
			RedirectURI    string
			SessionId      string
			ForumNavTrails []forumhelper.ForumNavTrail
		}
		userLoginPageData := UserLoginPageData{
			FormData:       formData,
			RedirectURI:    "./",
			SessionId:      session.SessionId,
			ForumNavTrails: []forumhelper.ForumNavTrail{},
		}

		// Execute template
		err = templateOutput.ExecuteTemplate(w, "overall", userLoginPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing user login template: %s", err)
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
