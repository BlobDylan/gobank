package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.HandleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandlerFunc(s.HandleGetAccountByID))

	log.Println("Starting server on", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.HandleGetAccounts(w, r)
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

func (s *APIServer) HandleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.db.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	fmt.Println("Getting account", id)

	return WriteJSON(w, http.StatusOK, &account{})
}

func (s *APIServer) HandleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := new(createAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accReq); err != nil {
		return err
	}

	account := NewAccount(accReq.Email)

	if err := s.db.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}
func (s *APIServer) HandleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
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
	db         storage
}

func NewAPIServer(listenAddr string, db storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		db:         db,
	}
}
