package main

import (
	"context"
	"io"
	"net/http"
	"time"
)

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

func (s serverGateway) loadBalancerHandler(w http.ResponseWriter, r *http.Request) {
	resource := getResource(r)
	pool, ok := s.lb.poolMap[resource.Prefix]
	if !ok {
		s.notFound(w, r)
		return
	}
	srv := pool.Next()
	srv.Add(1)
	defer srv.Add(-1)
	resp, err := forward(r, srv.address, resource.Timeout)

	if err != nil {
		s.internalServerError(w, r, err)
		return
	}

	defer resp.Body.Close()
	for k, v := range resp.Header {
		w.Header()[k] = v
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		s.internalServerError(w, r, err)
		return
	}
}
