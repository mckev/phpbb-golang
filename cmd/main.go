package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"text/template"

	"phpbb-golang/examples/myforum"
	"phpbb-golang/internal/logger"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpMethod := r.Method
	urlPath := filepath.Clean(r.URL.Path)
	logger.Debugf(ctx, "%s %s", httpMethod, urlPath)

	if urlPath == "/" {
		io.WriteString(w, "Welcome to Golang BB!")
	} else if urlPath == "/main" {
		templateOutput, _ := template.ParseFiles("./view/templates/overall.html", "./view/templates/main.html")
		templateOutput.ExecuteTemplate(w, "overall", nil)
	} else if urlPath == "/forums" {
		templateOutput, _ := template.ParseFiles("./view/templates/overall.html", "./view/templates/forums.html")
		templateOutput.ExecuteTemplate(w, "overall", nil)
	} else if urlPath == "/topics" {
		templateOutput, _ := template.ParseFiles("./view/templates/overall.html", "./view/templates/topics.html")
		templateOutput.ExecuteTemplate(w, "overall", nil)
	} else if urlPath == "/posts" {
		templateOutput, _ := template.ParseFiles("./view/templates/overall.html", "./view/templates/posts.html")
		templateOutput.ExecuteTemplate(w, "overall", nil)
	}
}

func init_example(ctx context.Context) {
	err := myforum.PopulateDb(ctx)
	if err != nil {
		logger.Errorf(ctx, "Unable to populate forums table: %s", err)
		return
	}
}

func main() {
	ctx := context.Background()

	init_example(ctx)

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
