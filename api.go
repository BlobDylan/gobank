package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2NvdW50TnVtYmVyIjo1OTMxLCJleHBpcmVzQXQiOjE1MDAwMH0.qUyEtqJu7jGyCEPAwVWNsojaO-h7Bw53U3idAEDSYpo
func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandlerFunc(s.HandleAccount))

	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandlerFunc(s.HandleAccountByID), s.db))

	router.HandleFunc("/transfer", makeHTTPHandlerFunc(s.HandleTransfer))

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

func (s *APIServer) HandleAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == http.MethodGet {
		return s.HandleGetAccountByID(w, r)
	}
	if r.Method == http.MethodDelete {
		return s.HandleDeleteAccountByID(w, r)
	}
	return fmt.Errorf("unsupported method %s", r.Method)
}

func (s *APIServer) HandleGetAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}
	account, err := s.db.GetAccount(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) HandleDeleteAccountByID(w http.ResponseWriter, r *http.Request) error {
	id, err := getIDFromRequest(r)
	if err != nil {
		return err
	}
	if err := s.db.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
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

	tokenstring, err := createJWTToken(account)
	if err != nil {
		return err
	}

	fmt.Println("Token:", tokenstring)

	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferreq := new(transferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferreq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferreq)
}

func Unauthorized(w http.ResponseWriter, e error) {
	fmt.Println("Unauthorized: ", e)
	WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Unauthorized"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, db storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Checking JWT token")
		tokenstring := r.Header.Get("jwt-token")
		token, err := validateJWTToken(tokenstring)
		if err != nil {
			Unauthorized(w, err)
			return
		}
		if !token.Valid {
			Unauthorized(w, err)
			return
		}
		userID, err := getIDFromRequest(r)
		if err != nil {
			Unauthorized(w, err)
			return
		}

		account, err := db.GetAccount(userID)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, ApiError{Error: "Unauthorized"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			Unauthorized(w, err)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWTToken(tokenstring string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenstring, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func createJWTToken(account *account) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt":     150000,
		"accountNumber": account.Number,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string `json:"error"`
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

func getIDFromRequest(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return 0, fmt.Errorf("id not found")
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("invalid id")
	}

	return id, nil
}
