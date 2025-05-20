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
	"phpbb-golang/internal/helper"
	"phpbb-golang/internal/logger"
	"phpbb-golang/model"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpMethod := r.Method
	urlPath := filepath.Clean(r.URL.Path)
	logger.Debugf(ctx, "%s %s", httpMethod, urlPath)

	if urlPath == "/" {
		io.WriteString(w, "Welcome to Golang BB!")
	} else if urlPath == "/myforum/main" {
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
	} else if urlPath == "/posts" {
		// Template Functions
		funcMap := template.FuncMap{
			"fnUnixTimeToStr": func(unixTime int64) string {
				return helper.UnixTimeToStr(unixTime)
			},
			"fnBbcodeToHtml": func(bbcodeStr string) template.HTML {
				// To print raw, unescaped HTML within a Go HTML template, the html/template package provides the template.HTML type. By converting a string containing HTML to template.HTML, you can instruct the template engine to render it as raw HTML instead of escaping it for safe output.
				// Warning: Since this Go template function outputs raw HTML, make sure it is safe from attacks such as XSS.
				return template.HTML(bbcode.ConvertBbcodeToHtml(bbcodeStr))
			},
		}

		// Prepare template files
		templateOutput, err := template.New("").Funcs(funcMap).ParseFiles("./view/templates/overall.html", "./view/templates/posts.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}

		// Prepare data
		TOPIC_ID := 2
		posts, err := model.ListPosts(ctx, TOPIC_ID)
		if err != nil {
			logger.Errorf(ctx, "Error while listing posts: %s", err)
			return
		}
		users, err := model.ListUsers(ctx, TOPIC_ID)
		usersMap := map[int]model.User{} // Convert users from a list into a map
		for _, user := range users {
			usersMap[user.UserId] = user
		}
		type PostsPageData struct {
			Posts    []model.Post
			UsersMap map[int]model.User
		}
		postsPageData := PostsPageData{
			Posts:    posts,
			UsersMap: usersMap,
		}

		// Go HTML Templates:
		//   - Go's html/template package automatically strips HTML comments (<!-- comment -->) during template execution.
		//   - Go templates do not inherently support nested template definitions in the way one might expect from other templating engines. While you can define templates within other templates using {{define}}, they are all effectively "hoisted" to the top level and treated as independent templates within a single namespace. This means you can't directly access a nested template as a property of its parent.
		//     To achieve a similar effect as nested templates, you should define each template separately and then compose them using the {{template}} action. When calling {{template}}, you can pass data to the included template using a dot ., which represents the current data context.
		//     If you need to access data from a parent template, you can either pass it explicitly to the nested template or use variables to store the data before calling the nested template.
		//   - In Go templates, . refers to the current context, which changes within a range loop to the current element being iterated over.
		//     To access the outer struct's fields from within the inner loop, $ is used, which always refers to the root context (the original data passed to the template).
		err = templateOutput.ExecuteTemplate(w, "overall", postsPageData)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}
	} else {
		logger.Errorf(ctx, "URL Path not supported: %s", urlPath)
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
