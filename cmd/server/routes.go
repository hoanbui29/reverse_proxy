package main

import "net/http"

func routes() http.Handler {
	mux := http.NewServeMux()
	return mux
}
