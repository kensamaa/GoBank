package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	jwt "github.com/golang-jwt/jwt/v4"

	//installed by running in command go get github.com/gorilla/mux
	"github.com/gorilla/mux"
)
type APIServer struct {
	listenAddr string
	store storage
}

//return the pointer to api server	instance
func NewAPIServer(listenAddr string,store storage) *APIServer {
	return &APIServer{ 
        listenAddr: listenAddr,
		store:store,
    }
}
func (s *APIServer)Run()  {
	router:=mux.NewRouter()
	router.HandleFunc("/account",makeHTTPHandleFunc(s.handleAccount))

	router.HandleFunc("/login",makeHTTPHandleFunc(s.handleLogin))
	
	router.HandleFunc("/account/{id}",makeHTTPHandleFunc(s.handleGetAccountById))
	log.Println("JSON API Server running on port : ",s.listenAddr)
	http.ListenAndServe(s.listenAddr,router)
}
//handle requests to the API server
func (s *APIServer)handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method=="GET"{
		return s.handleGetAccount(w,r)
	}
	if r.Method=="POST"{
		return s.handleCreateAccount(w,r)
	}
	if r.Method=="DELETE"{
		return s.handleDeleteAccount(w,r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}
func (s * APIServer) handleLogin(w http.ResponseWriter, r *http.Request)error {
	if r.Method != "POST"{
		return fmt.Errorf("method not allowed %s",r.Method)
	}
	var req LoginRequest
	if err:=json.NewDecoder(r.Body).Decode(&req);err !=nil{
		return err
	}
	acc, err := s.store.GetAccountByNumber(int(req.Number))
	if err != nil {
		return err
	}
	if !acc.ValidPassword(req.Password) {
		return fmt.Errorf("not authenticated")
	}

	token, err := createJWT(acc)
	if err != nil {
		return err
	}

	resp := LoginResponse{
		Token:  token,
		Number: acc.Number,
	}

	return WriteJson(w, http.StatusOK, resp)
}
func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}
//GET /accounts
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}
func (s *APIServer) handleGetAccountById (w http.ResponseWriter, r *http.Request)error {
	if r.Method == "GET" {
		id, err := getID(r)
		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)
		if err != nil {
			return err
		}

		return WriteJson(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	account, err := NewAccount(req.FirstName, req.LastName,req.Password)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}
func (s *APIServer)handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
		if err != nil {
			return err
		}
	s.store.DeleteAccount(id)
	return WriteJson(w, http.StatusOK, id)
}
func (s *APIServer)handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}

//helper function 
type ApiError struct {Error string}
type apiFunc func(http.ResponseWriter, *http.Request) error
func WriteJson(w http.ResponseWriter, status int,v any) error {
	w.Header().Add("Content-Type", "application/json")
    w.WriteHeader(status)
    return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		if err:=f(w,r); err!=nil{
			//handle the error
			WriteJson(w,http.StatusBadRequest,ApiError{Error:err.Error()})
		}
	}
}
func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}