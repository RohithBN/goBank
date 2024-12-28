package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type APIserver struct {
	ListenAddress string
	store         Storage
}
type APIerror struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

type apifunc func(w http.ResponseWriter, r *http.Request) error

func handleApiFunc(f apifunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusInternalServerError, APIerror{Error: err.Error()})
		}
	}
}

func newApiServer(ListenAddress string, store Storage) *APIserver {
	return &APIserver{
		ListenAddress: ListenAddress,
		store:         store,
	}
}

func (s *APIserver) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/account/{id}", handleApiFunc(s.handleGetAccountById))
	router.HandleFunc("/account", handleApiFunc(s.handleAccount))
	fmt.Println("Listening on server", s.ListenAddress)
	http.ListenAndServe(s.ListenAddress, router)
}

func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("Received %s request to /account\n", r.Method)
	if r.Method == "GET" {
		return s.handleGetAccounts(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return nil
}

func (s *APIserver) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	fmt.Printf("Received %s request to /accounts\n", r.Method)
	accounts, err := s.store.getAllAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}
func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := createAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&accReq); err != nil {
		return err
	}
	account := NewAccount(accReq.FirstName, accReq.LastName)
	if err := s.store.createAccount(account); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusAccepted, account)
}
func (s *APIserver) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	account, err := s.store.getAccountById(id)
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, account)
}
func (s *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
