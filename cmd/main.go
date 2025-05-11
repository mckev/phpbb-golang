package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"phpbb-golang/examples/myforum"
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
	} else {
		logger.Errorf(ctx, "URL Path not supported: %s", urlPath)
		return
	}
}

func init_example(ctx context.Context) {
	err := myforum.PopulateDb(ctx)
	if err != nil {
		logger.Errorf(ctx, "Error while populating database: %s", err)
		return
	}
}

func main() {
	ctx := context.Background()

	{
		init_example(ctx)
		forums, err := model.ListForums(ctx)
		if err != nil {
			logger.Errorf(ctx, "Error while listing forums: %s", err)
		}
		for _, forum := range forums {
			logger.Infof(ctx, "%s", helper.JsonDumps(forum))
		}
	}

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
