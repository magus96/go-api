package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIError struct {
	Error string
}

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id := mux.Vars(r)["id"]
		if id != "" {
			return s.handleGetAccountbyId(w, r)
		}
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("Method not allowed: %s", r.Method)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST"{
		return fmt.Errorf("Method not allowed: %s", r.Method)
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountbyNumber(req.Number)
	if err != nil{
		return err
	}
	fmt.Printf("%+v\n", acc)


	return writeJSON(w, http.StatusOK, req)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountbyId(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		id, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			return fmt.Errorf("Invalid id: %d", id)
		}
		acc, err := s.store.GetAccountbyID(id)
		if err != nil {
			return err
		}
		return writeJSON(w, http.StatusOK, acc)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	} else {
		return fmt.Errorf("Invalid method")
	}
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	accReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(accReq); err != nil {
		return err
	}

	account, err := newAccount(accReq.FirstName, accReq.LastName, accReq.Password)
	if err != nil {
		return err
	}

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	token, err := createJWT(account)
	if err != nil {
		return err
	}
	fmt.Println(token)
	return writeJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		return fmt.Errorf("Invalid id: %d", id)
	}
	err = s.store.DeleteAccount(id)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, map[string]int{"Deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("Method not allowed: %s", r.Method)
	}
	transferReq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()
	return writeJSON(w, http.StatusOK, transferReq)
}

func newAPIServer(listener string, store Storage) *APIServer {
	return &APIServer{listenAddr: listener, store: store}
}

func (s *APIServer) run() {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTOut(makeHTTPHandleFunc(s.handleGetAccountbyId), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("API Server running on port:", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle err
			writeJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}
