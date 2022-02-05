package server

import (
	"context"
	"encoding/json"
	"github.com/go-paxos/domain"
	"github.com/go-paxos/logger"
	"github.com/go-paxos/roles"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/tryfix/log"
	traceableContext "github.com/tryfix/traceable-context"
	"io/ioutil"
	"net/http"
	"strconv"
)

type server struct {
	leader  *roles.Leader
	replica *roles.Replica
	logger  log.Logger
}

func Init(ctx context.Context, port int, leader *roles.Leader, replica *roles.Replica) {
	s := &server{leader: leader, replica: replica, logger: logger.Log}

	r := mux.NewRouter()
	r.HandleFunc(domain.RequestReplicaEndpoint, s.handleClientRequest).Methods(http.MethodPost)
	r.HandleFunc(domain.UpdateReplicaEndpoint, s.handleUpdateReplica).Methods(http.MethodPost)

	r.HandleFunc(domain.RequestLeaderEndpoint, s.handleReplicaRequest).Methods(http.MethodPost)
	r.HandleFunc(domain.PrepareEndpoint, s.handlePrepare).Methods(http.MethodPost)
	r.HandleFunc(domain.AcceptEndpoint, s.handleAccept).Methods(http.MethodPost)

	s.logger.InfoContext(ctx, `initializing http server`)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), r))
}

// todo change errors

// handleClientRequest handles the client request with a string value in raw body and passes the decoded value to replica
// to initiate the procedure.
func (s *server) handleClientRequest(w http.ResponseWriter, r *http.Request) {
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
func (s *server) handleUpdateReplica(w http.ResponseWriter, r *http.Request) {
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
func (s *server) handleReplicaRequest(w http.ResponseWriter, r *http.Request) {
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
func (s *server) handlePrepare(w http.ResponseWriter, r *http.Request) {
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
func (s *server) handleAccept(w http.ResponseWriter, r *http.Request) {
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
