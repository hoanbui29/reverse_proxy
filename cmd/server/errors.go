package main

import (
	"net/http"
)

type envelope map[string]interface{}

func (app *serverGateway) logError(r *http.Request, err error) {
	app.logger.Error(err, map[string]string{
		"request_method": r.Method,
		"request_url":    r.URL.String(),
	})
}

func (s *serverGateway) errorResponse(w http.ResponseWriter, r *http.Request, code int, message any) {
	env := envelope{
		"error": message,
	}

	err := s.writeJSON(w, code, env, nil)

	if err != nil {
		s.logError(r, err)
		w.WriteHeader(500)
	}
}

func (s *serverGateway) handleFatal(err error) {
	if err != nil {
		s.logger.Fatal(err, nil)
	}
}

func (s *serverGateway) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	s.logError(r, err)
	s.errorResponse(w, r, http.StatusInternalServerError, "internal server error")
}

func (s *serverGateway) notFound(w http.ResponseWriter, r *http.Request) {
	s.errorResponse(w, r, http.StatusNotFound, "not found")
}

func (s *serverGateway) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	s.errorResponse(w, r, http.StatusMethodNotAllowed, "method not allowed")
}
