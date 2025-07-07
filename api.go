package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

type ApiError struct {
	Error string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//error handling here
			WriteJSON(w, http.StatusBadRequest, ApiError{
				Error: err.Error(),
			})
		}

	}
}

type ApiServer struct {
	listenAddr string
	store      Storage
}

func NewApiServer(listenAddr string, store Storage) *ApiServer {

	return &ApiServer{
		listenAddr: listenAddr,
		store:      store,
	}

}
func (s *ApiServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHttpHandleFunc(s.handleAcount))

	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleGetAccount))

	fmt.Printf("JSON API server running on port: %v", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}
func (s *ApiServer) handleAcount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)

	case "POST":
		return s.handleCreateAccount(w, r)

	case "UPDATE":

	case "DELETE":
		return s.handleDeleteAccount(w, r)

	case "T":
		break

	}
	return fmt.Errorf("method not allowed %s", r.Method)

}
func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)

	return WriteJSON(w, 200, vars)

}
func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}
func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
