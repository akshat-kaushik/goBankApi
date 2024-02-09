package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
)

func WriteJson(w http.ResponseWriter, status int, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type APIserver struct {
	listenAddress string
	store         Storage
}

func NewAPIserver(listenAddress string, store Storage) *APIserver {
	return &APIserver{
		listenAddress: listenAddress,
		store:         store,
	}
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandler(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJson(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *APIserver) Start() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))
	router.HandleFunc("/account/{id}", JWTauth(makeHTTPHandler(s.handleAccountByID)))
	router.HandleFunc("/transfer/{accountNumber}", makeHTTPHandler(handleTrasferRequest))

	log.Println("Starting server on", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("unsupported method")
}

func (s *APIserver) handleAccountByID(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountByID(w, r)
	}
	// if r.Method == "PUT" {
	// 	return s.handleUpdateAccount(w, r)
	// }
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("unsupported method")
}

func (s *APIserver) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {

	accountID, err := getId(r)
	if err != nil {
		return err
	}

	account, err := s.store.getAccountByID(accountID)
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.getAccounts()
	if err != nil {
		return err
	}
	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var createAccountReq createAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}
	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName, createAccountReq.Age)
	if err := s.store.createAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	accountID, err := getId(r)
	if err != nil {
		return err
	}

	err = s.store.deleteAccount(accountID)
	if err != nil {
		return WriteJson(w, http.StatusBadRequest, "account not found")
	}

	return WriteJson(w, http.StatusOK, map[string]int{"deleted": accountID})
}

func getId(r *http.Request) (int, error) {
	id := mux.Vars(r)["id"]
	accountID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("invalid account id")
	}
	return accountID, nil
}

func handleTrasferRequest(w http.ResponseWriter, r *http.Request) error {
	transferReq := &TransereRequest{}
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return nil
}

func JWTauth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := validateToken(tokenString)
		if err != nil {
			WriteJson(w, http.StatusUnauthorized, ApiError{Error: "Invalid token"})
			return
		}
		fmt.Println("JWTauth")
		handlerFunc(w, r)
	}
}

	var secret = "my_secret"

	func validateToken(tokenString string) (*jwt.Token, error) {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil {
			return nil, err
		}
		return token, nil
	} 
