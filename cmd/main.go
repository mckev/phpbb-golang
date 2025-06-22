package main

import (
	"context"
	"fmt"
	"net/http"

	"phpbb-golang/controller"
	"phpbb-golang/examples/myforum"
	"phpbb-golang/internal/logger"
)

func httpHandler() http.Handler {
	rootMux := http.NewServeMux()

	// Public (no middleware)
	// Note: The "/" suffix will do prefix match, otherwise it will do exact match
	rootMux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("./view/static/assets/"))))
	rootMux.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("./view/static/images/"))))
	rootMux.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("./view/static/styles/"))))
	rootMux.HandleFunc("/myforum/", controller.MyForumPage)

	// With Session middleware
	sessionMux := http.NewServeMux()
	sessionMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		controller.MainPage(w, r)
	})
	sessionMux.HandleFunc("/redirect", http.HandlerFunc(controller.RedirectPage))
	sessionMux.HandleFunc("/forums", http.HandlerFunc(controller.ForumsPage))
	sessionMux.HandleFunc("/topics", http.HandlerFunc(controller.TopicsPage))
	sessionMux.HandleFunc("/posts", http.HandlerFunc(controller.PostsPage))
	sessionMux.HandleFunc("/user_register", http.HandlerFunc(controller.UserRegisterPage))
	sessionMux.HandleFunc("/user_login", http.HandlerFunc(controller.UserLoginPage))
	rootMux.Handle("/", controller.SessionMiddleware(sessionMux))

	return rootMux
}

func main() {
	ctx := context.Background()

	// MyForum example
	// TODO: Delete after finish development
	myforum.InitMyforum(ctx)
	myforum.DebugMyforum(ctx)

	portNumber := 9000
	logger.Infof(ctx, "Server is listening on port %d", portNumber)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", portNumber), httpHandler())
	if err != nil {
		logger.Fatalf(ctx, "Error while listening on port %d: %s", portNumber, err)
	}
}
