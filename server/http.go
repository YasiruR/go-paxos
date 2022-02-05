package server

import (
	"encoding/json"
	"github.com/go-paxos/domain"
	"github.com/go-paxos/roles"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	traceableContext "github.com/tryfix/traceable-context"
	"io/ioutil"
	"net/http"
)

type Server struct {
	replica *roles.Replica
	leader  *roles.Leader
	logger  log.Logger
}

func Init() {
	s := &Server{}
	r := mux.NewRouter()
	r.HandleFunc(domain.RequestReplicaEndpoint, s.handleClientRequest).Methods(http.MethodPost)
	r.HandleFunc(domain.UpdateReplicaEndpoint, s.handleUpdateReplica).Methods(http.MethodPost)

	r.HandleFunc(domain.RequestLeaderEndpoint, s.handleReplicaRequest).Methods(http.MethodPost)
	r.HandleFunc(domain.PrepareEndpoint, s.handlePrepare).Methods(http.MethodPost)
	r.HandleFunc(domain.AcceptEndpoint, s.handleAccept).Methods(http.MethodPost)
}

// todo change errors

// handleClientRequest handles the client request with a string value in raw body and passes the decoded value to replica
// to initiate the procedure.
func (s *Server) handleClientRequest(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.replica.HandleRequest(string(data))
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleUpdateReplica updates the state of the current node whenever a consensus is reached and sent by leaders
func (s *Server) handleUpdateReplica(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var dec domain.Decision
	err = json.Unmarshal(data, &dec)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = s.replica.Update(dec)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handleReplicaRequest handles the request by a replica and forwards to the leader layer to proceed with a proposal
func (s *Server) handleReplicaRequest(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var req domain.Request
	err = json.Unmarshal(data, &req)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	ok, err := s.leader.Propose(req)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		s.logger.DebugContext(ctx, `proposed value was not chosen`)
		w.WriteHeader(http.StatusNotAcceptable)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// handlePrepare handles the prepare request sent by the proposer to an acceptor with initialization of a proposal
func (s *Server) handlePrepare(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var prop domain.Proposal
	err = json.Unmarshal(data, &prop)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	promise, err := s.leader.HandlePrepare(prop)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&promise)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// handleAccept handles accept requests by the proposer to confirm a proposal
func (s *Server) handleAccept(w http.ResponseWriter, r *http.Request) {
	ctx := traceableContext.WithUUID(uuid.New())
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var prop domain.Proposal
	err = json.Unmarshal(data, &prop)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	accept, err := s.leader.HandleAccept(prop)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(&accept)
	if err != nil {
		s.logger.ErrorContext(ctx, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
