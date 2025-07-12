package controller

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"phpbb-golang/internal/forumhelper"
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

type PostWriteMode int

const (
	POST_WRITE_MODE_COMPOSE PostWriteMode = iota
	POST_WRITE_MODE_QUOTE_A_POST
	POST_WRITE_MODE_PREVIEW
	POST_WRITE_MODE_SUBMIT
)

func PostWritePage(w http.ResponseWriter, r *http.Request) {
	// To try: http://localhost:9000/post_write?t=2
	ctx := r.Context()
	session := getSession(r)
	queryParams := r.URL.Query()
	topicId := helper.StrToInt(queryParams.Get("t"), model.INVALID_TOPIC_ID)

	mode := POST_WRITE_MODE_COMPOSE
	type FormData struct {
		Subject string
		Message string
		Errors  []string
	}
	formData := FormData{}

	switch r.Method {
	case "POST":
		err := r.ParseForm()
		if err != nil {
			logger.Errorf(ctx, "Error while parsing form upon post write: %s", err)
			return
		}
		formData.Subject = strings.TrimSpace(r.PostForm.Get("subject"))
		if formData.Subject == "" {
			formData.Errors = append(formData.Errors, "Subject is required. Please provide a subject to proceed.")
		}
		formData.Message = strings.TrimSpace(r.PostForm.Get("message"))
		if formData.Message == "" {
			formData.Errors = append(formData.Errors, "Message content is required. Please enter your message to proceed.")
		}
		if r.PostForm.Get("post") == "Submit" {
			mode = POST_WRITE_MODE_SUBMIT
		}
		if r.PostForm.Get("preview") == "Preview" {
			mode = POST_WRITE_MODE_PREVIEW
		}
		var postId int
		if len(formData.Errors) == 0 && mode == POST_WRITE_MODE_SUBMIT {
			// Validation successful

			// fmt.Fprintf(w, "Form submitted successfully!\n")
			// fmt.Fprintf(w, "Subject: %s\n", formData.Subject)
			// fmt.Fprintf(w, "Message: %s\n", formData.Message)
			// TODO: Handle CSRF token validation

			topic, err := model.GetTopic(ctx, topicId)
			if err != nil {
				logger.Errorf(ctx, "Error while getting topic id %d: %s", topicId, err)
			}
			postId, err = InsertPost(ctx, topicId, topic.ForumId, formData.Subject, formData.Message, session.SessionUserId, session.SessionUserName)
			if err != nil {
				logger.Errorf(ctx, "Error while inserting post subject '%s' with topic id %d for user id %d: %s", formData.Subject, topicId, session.SessionUserId, err)
				formData.Errors = append(formData.Errors, "We're experiencing a temporary issue. Your request couldn't be completed.")
			}
		}
		if len(formData.Errors) == 0 && mode == POST_WRITE_MODE_SUBMIT {
			http.Redirect(w, r, helper.UrlWithSID(fmt.Sprintf("./posts?p=%d#p%d", postId, postId), session.SessionId), http.StatusFound)
			return
		}

		fallthrough

	case "GET":
		// Case reply (quote) from a specific Post
		if queryParams.Get("mode") == "quote" {
			mode = POST_WRITE_MODE_QUOTE_A_POST
			postId := helper.StrToInt(queryParams.Get("p"), model.INVALID_POST_ID)
			if postId > 0 {
				post, err := model.GetPost(ctx, postId)
				if err != nil {
					logger.Errorf(ctx, "Error while getting post id %d: %s", postId, err)
					return
				}
				topicId = post.TopicId
				formData.Subject = post.PostSubject
				formData.Message = fmt.Sprintf("[blockquote user_name=%s user_id=%d post_id=%d time=%d]\n%s\n[/blockquote]", helper.FormatAttributeValue(post.PostUserName), post.PostUserId, post.PostId, post.PostTime, post.PostText)
			}
		}

		// Prepare data
		topic, err := model.GetTopic(ctx, topicId)
		if err != nil {
			logger.Errorf(ctx, "Error while getting topic id %d: %s", topicId, err)
			return
		}
		if formData.Subject == "" {
			formData.Subject = fmt.Sprintf("Re: %s", topic.TopicTitle)
		}

		// Topic Review panel
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
			FormData                FormData
			Posts                   []model.Post
			UsersMap                map[int]model.User
			Session                 model.Session
			RedirectURIForLoginPage string
			ForumNavTrails          []forumhelper.ForumNavTrail
		}
		postWritePageData := PostWritePageData{
			Topic:                   topic,
			FormData:                formData,
			Posts:                   posts,
			UsersMap:                usersMap,
			Session:                 session,
			RedirectURIForLoginPage: helper.UrlWithSID(r.URL.RequestURI(), ""),
			ForumNavTrails:          forumNavTrails,
		}

		// Execute template
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/post_write.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing post write template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", postWritePageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing post write template: %s", err)
			return
		}
	}
}

func InsertPost(ctx context.Context, topicId int, forumId int, postSubject string, postText string, postUserId int, postUserName string) (int, error) {
	err := model.CheckIfGuestUser(ctx, postUserId, postUserName)
	if err != nil {
		return model.INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d for user id %d: %s", postSubject, topicId, postUserId, err)
	}
	postId, err := model.InsertPost(ctx, topicId, forumId, postSubject, postText, postUserId, postUserName)
	if err != nil {
		return model.INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d for user id %d: %s", postSubject, topicId, postUserId, err)
	}
	err = model.UpdateLastPostOfTopic(ctx, topicId, postId, postUserId, postUserName)
	if err != nil {
		return model.INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d for user id %d: %s", postSubject, topicId, postUserId, err)
	}
	err = model.UpdateLastPostOfForum(ctx, forumId, postId, postSubject, postUserId, postUserName)
	if err != nil {
		return model.INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d for user id %d: %s", postSubject, topicId, postUserId, err)
	}
	err = model.IncreaseNumPostsForUser(ctx, postUserId)
	if err != nil {
		return model.INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d for user id %d: %s", postSubject, topicId, postUserId, err)
	}
	err = model.IncreaseNumPostsForTopic(ctx, topicId)
	if err != nil {
		return model.INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d for user id %d: %s", postSubject, topicId, postUserId, err)
	}
	err = model.IncreaseNumPostsForForum(ctx, forumId)
	if err != nil {
		return model.INVALID_POST_ID, fmt.Errorf("Error while inserting post subject '%s' with topic id %d for user id %d: %s", postSubject, topicId, postUserId, err)
	}
	return postId, nil
}
