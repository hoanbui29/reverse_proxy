package main

import (
	"context"
	"fmt"
	"hoanbui29/reverse_proxy/internal/config"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"
)

func (s serverGateway) getResource(r *http.Request) *config.Resource {
	for _, v := range s.config.Resources {
		if strings.HasPrefix(r.URL.Path, v.Prefix) {
			return &v
		}
	}
	return nil
}

func forward(r *http.Request, host string, timeout int) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Duration(timeout)*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, r.Method, host+r.URL.Path, r.Body)
	req.Header = r.Header
	req.Header.Set("X-Forwarded-Host", r.Host)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)

	return resp, nil
}

func (s serverGateway) proxy() http.Handler {
	factoryMap := make(map[string]IResourceFactory)
	fmt.Println(factoryMap)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resource := s.getResource(r)
		if resource == nil {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		if !slices.Contains(resource.Methods, r.Method) {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var factory IResourceFactory
		var err error

		fmt.Println(factoryMap)

		if v, ok := factoryMap[resource.Prefix]; !ok {
			factory, err = GetResourceFactory(resource)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			factoryMap[resource.Prefix] = factory
		} else {
			factory = v
		}

		host := factory.Next()
		resp, err := forward(r, host, resource.Timeout)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		for k, v := range resp.Header {
			w.Header()[k] = v
		}

		w.WriteHeader(resp.StatusCode)

		// Copy the response body to the client
		_, err = io.Copy(w, resp.Body)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
