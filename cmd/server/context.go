package main

import (
	"context"
	"hoanbui29/reverse_proxy/internal/config"
	"net/http"
)

var resourceKey = "resource"

func setResource(r *http.Request, resource *config.Resource) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), resourceKey, resource))
}

func getResource(r *http.Request) *config.Resource {
	resource, ok := r.Context().Value(resourceKey).(*config.Resource)
	if !ok {
		panic("resource not found")
	}
	return resource
}
