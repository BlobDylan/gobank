package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, ApiError{Error: err.Error()})
		}

	}
}

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.HandleAccount))

	log.Println("Starting server on", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.HandleGetAccount(w, r)
	}
	if r.Method == http.MethodPost {
		return s.HandleCreateAccount(w, r)
	}
	if r.Method == http.MethodDelete {
		return s.HandleDeleteAccount(w, r)
	}
	if r.Method == http.MethodPut {
		return s.HandleTransfer(w, r)
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) HandleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
