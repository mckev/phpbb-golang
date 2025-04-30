package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"

	"golang-bb/internal/logger"
)

func serveTemplate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	httpMethod := r.Method
	urlPath := filepath.Clean(r.URL.Path)
	logger.Debugf(ctx, "%s %s", httpMethod, urlPath)

	templateOutput, _ := template.ParseFiles("./view/templates/overall.html")
	templateOutput.ExecuteTemplate(w, "layout", nil)
}

func main() {
	ctx := context.Background()
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./view/static/assets/"))))
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./view/static/styles/"))))
	http.HandleFunc("/", serveTemplate)
	portNumber := 9000
	logger.Infof(ctx, "Server is listening on port %d", portNumber)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", portNumber), nil)
	if err != nil {
		logger.Fatalf(ctx, "Error while listening on port %d: %s", portNumber, err)
	}
}
