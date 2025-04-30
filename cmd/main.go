package main

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"golang-bb/internal/logger"
)

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logger.Debugf(ctx, "GET /")
	io.WriteString(w, "Welcome to Golang BB!")
}

func main() {
	ctx := context.Background()
	http.HandleFunc("/", getRoot)
	portNumber := 9000
	logger.Infof(ctx, "Server is listening on port %d", portNumber)
	err := http.ListenAndServe(fmt.Sprintf("localhost:%d", portNumber), nil)
	if err != nil {
		logger.Fatalf(ctx, "Error while listening on port %d: %s", portNumber, err)
	}
}
