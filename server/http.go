package server

import (
	"encoding/json"
	"github.com/go-paxos/domain"
	"github.com/go-paxos/roles"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
)

type Server struct {
	replica *roles.Replica
	leader  *roles.Leader
}

func Init() {
	s := &Server{}
	r := mux.NewRouter()
	r.HandleFunc(`/replica`, s.handleRequest).Methods(http.MethodPost)
	r.HandleFunc(`/leader/request`, s.handleLeaderRequest).Methods(http.MethodPost)
	r.HandleFunc(domain.PrepareEndpoint, s.handlePrepare).Methods(http.MethodPost)
	r.HandleFunc(domain.AcceptEndpoint, s.handleAccept).Methods(http.MethodPost)
}

// endpoint for replica to receive request
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// err
	}

}

// endpoint for leader to receive request
func (s *Server) handleLeaderRequest(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// err
	}

	var req domain.Request
	err = json.Unmarshal(data, &req)
	if err != nil {
		// err
	}

	// send req to init proposal
}

// endpoint for leader to receive proposal
func (s *Server) handlePrepare(w http.ResponseWriter, r *http.Request) {

}

// endpoint for leader to receive accept
func (s *Server) handleAccept(w http.ResponseWriter, r *http.Request) {

}
