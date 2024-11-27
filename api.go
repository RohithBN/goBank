package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type APIserver struct {
	ListenAddress string
}
type APIerror struct {
	Error string
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
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

func newApiServer(ListenAddress string) *APIserver {
	return &APIserver{
		ListenAddress: ListenAddress,
	}
}

func (s *APIserver) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/", handleApiFunc(s.handleAccount))
	fmt.Println("Listening on server", s.ListenAddress)
	http.ListenAndServe(s.ListenAddress, router)
}

func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	if r.Method=="DELETE"{
		return s.handleDeleteAccount(w, r)
	}
	return nil
}
func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *APIserver) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
