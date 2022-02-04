package server

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Init() {
	r := mux.NewRouter()
	r.HandleFunc(`/replica`, handleRequest).Methods(http.MethodPost)
	r.HandleFunc(`/leader/request`, handleLeaderRequest).Methods(http.MethodPost)
	r.HandleFunc(`/leader/proposal`, handleProposal).Methods(http.MethodPost)
}

// endpoint for replica to receive request
func handleRequest(w http.ResponseWriter, r *http.Request) {

}

// endpoint for leader to receive request
func handleLeaderRequest(w http.ResponseWriter, r *http.Request) {

}

// endpoint for leader to receive proposal
func handleProposal(w http.ResponseWriter, r *http.Request) {

}
