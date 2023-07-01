package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

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

	account, err := NewAccount(req.FirstName, req.LastName)
	if err != nil {
		return err
	}
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}
func (s *APIServer)handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
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