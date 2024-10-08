package main

import (
	"encoding/json"
	"net/http"
)

func (s *serverGateway) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// This marshalIndent feature might affects the performance of the application
	// js, err := json.MarshalIndent(data, "-", "\t")
	js, err := json.Marshal(data)

	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
