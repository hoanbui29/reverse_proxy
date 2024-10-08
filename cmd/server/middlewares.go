package main

import (
	"net/http"
	"slices"
	"strings"
)

//rate limiting

func (s serverGateway) getResourceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, v := range s.config.Resources {
			if strings.HasPrefix(r.URL.Path, v.Prefix) {
				r = setResource(r, &v)
				next.ServeHTTP(w, r)
				return
			}
		}
		s.notFound(w, r)
		return
	})
}

func (s serverGateway) checkAllowedMethods(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource := getResource(r)
		if !slices.Contains(resource.Methods, r.Method) {
			s.methodNotAllowed(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s serverGateway) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				s.logger.Error(err.(error), nil)
				s.internalServerError(w, r, err.(error))
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
