package controller

import (
	"html/template"
	"net/http"
	"path/filepath"

	"phpbb-golang/internal/logger"
)

func MyForumPage(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpMethod := r.Method
	urlPath := filepath.Clean(r.URL.Path)
	logger.Debugf(ctx, "%s %s", httpMethod, urlPath)

	switch urlPath {
	case "/myforum/main":
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

	case "/myforum/forums":
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

	case "/myforum/topics":
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

	case "/myforum/posts":
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

	case "/myforum/post_write":
		// To try: http://localhost:9000/myforum/post_write
		templateOutput, err := template.ParseFiles("./examples/myforum/templates/overall.html", "./examples/myforum/templates/post_write.html")
		if err != nil {
			logger.Errorf(ctx, "Error while parsing template files: %s", err)
			return
		}
		err = templateOutput.ExecuteTemplate(w, "overall", nil)
		if err != nil {
			logger.Errorf(ctx, "Error while executing template: %s", err)
			return
		}

	case "/myforum/user_login":
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

	case "/myforum/user_register":
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

	case "/myforum/user_register_created":
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

	case "/myforum/user_register_activated":
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

	default:
		http.NotFound(w, r)
		return
	}
}
