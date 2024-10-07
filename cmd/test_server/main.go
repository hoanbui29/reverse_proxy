package main

import (
	"fmt"
	"hoanbui29/reverse_proxy/internal/logger"
	"net/http"
	"os"
)

type app struct {
	logger *logger.Logger
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Hello from %s", os.Getenv("HOSTNAME"))))
	})
	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), mux)

	if err != nil {
		panic(err)
	}
}
