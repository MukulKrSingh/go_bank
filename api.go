package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
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

	router.HandleFunc("/login", makeHttpHandleFunc(s.handleLogin))

	router.HandleFunc("/account", withJwt(makeHttpHandleFunc(s.handleAcount)))
	router.HandleFunc("/account/{id}", makeHttpHandleFunc(s.handleGetAccountById))
	router.HandleFunc("/account/delete/{id}", makeHttpHandleFunc(s.handleDeleteAccount))

	router.HandleFunc("/transfer", makeHttpHandleFunc(s.handleTransfer))
	fmt.Printf("JSON API server running on port: %v", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

}
func (s *ApiServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusForbidden, ApiError{
			Error: "INVALID REQUEST",
		})
	}
	loginReq := new(LoginRequest)
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		return WriteJSON(w, http.StatusBadRequest, ApiError{
			Error: "INVALID REQUEST",
		})
	}

	return WriteJSON(w, http.StatusOK, loginReq)

}
func (s *ApiServer) handleAcount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "UPDATE":
		break
	case "DELETE":
		return s.handleDeleteAccount(w, r)

	case "T":
		break

	}
	return fmt.Errorf("method not allowed %s", r.Method)

}
func (s *ApiServer) handleGetAccount(w http.ResponseWriter, _ *http.Request) error {

	res, err := s.store.GetAccounts()
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, err)
	}
	return WriteJSON(w, http.StatusOK, res)

}
func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccReq); err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}
	//1. Validate the pwd
	isValid, err := validatePwd(createAccReq.Password)
	if !isValid || err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	//2. Encrypt the pwd
	encpw, err := encrypt(createAccReq.Password)
	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, ApiError{
			Error: err.Error(),
		})
	}
	pwd := string(encpw)
	//3. Store the pwd in new account in encrypted format
	account := NewAccount(
		createAccReq.FirstName,
		createAccReq.LastName,
		pwd,
	)

	if err := s.store.CreateAccount(account); err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, account)
}
func (s *ApiServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	i, err := getID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}
	res, err := s.store.GetAccountById(i)
	if err != nil {
		return WriteJSON(w, http.StatusNotFound, err)
	}
	return WriteJSON(w, http.StatusOK, res)

}
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		return WriteJSON(w, http.StatusBadRequest, err)
	}

	return WriteJSON(w, http.StatusOK, map[string]int{
		"deleted": id,
	})
}
func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(TransferRequest)
	json.NewDecoder(r.Body).Decode(&transferReq)

	return WriteJSON(w, http.StatusOK, transferReq)
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id is given %s", idStr)
	}
	return id, nil
}
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, ApiError{
		Error: "access denied",
	})
}
func withJwt(handleFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// tokenString := r.Header.Get("x-jwt-token")
		// token, err := validateJwt(tokenString)
		// if err != nil {
		// 	permissionDenied(w)
		// 	return
		// }
		// if !token.Valid {
		// 	permissionDenied(w)
		// 	return
		// }
		fmt.Println("Called with jwt auth")
		//get the user id from the getId
		//match it with jwt claims
		handleFunc(w, r)

	}
}
func createJwt(account *Account) (string, error) {
	var token string
	return token, nil
}
func validateJwt(tokenString string) (*jwt.Token, error) {

	return nil, nil
}
func validatePwd(pwd string) (bool, error) {
	if len(pwd) < 8 {
		return false, fmt.Errorf("password should be min 8 chars long")
	}
	return true, nil
}
func encrypt(v string) ([]byte, error) {
	enc, err := bcrypt.GenerateFromPassword([]byte(v), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return enc, nil
}
