package roles

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-paxos/domain"
	"github.com/go-paxos/logger"
	"net/http"
	"sync"
)

// todo add all errs to one file
const (
	errNoLeader    = `no leader found in the replica`
	errInvalidSlot = `received a decision for an invalid slot`
)

type Replica struct {
	hostname string
	log      []string
	leaders  []string
	client   *http.Client
	lock     *sync.Mutex
}

func NewReplica() *Replica {
	// add leaders
	return &Replica{}
}

// HandleRequest builds a request from the client value and forwards the request received to a leader
func (r *Replica) HandleRequest(val string) error {
	req := r.buildRequest(val)
	for {
		ok, err := r.send(req)
		if err != nil {
			return logger.ErrorWithLine(err)
		}

		if ok {
			return nil
		}

		r.lock.Lock()
		if req.SlotID == len(r.log) {
			req.SlotID++
		} else {
			req.SlotID = len(r.log)
		}
		r.lock.Unlock()
	}
}

// buildRequest builds the request from the client value
func (r *Replica) buildRequest(val string) domain.Request {
	return domain.Request{
		Replica: r.hostname,
		SlotID:  len(r.log),
		Val:     val,
	}
}

// Sends the request to first leader found in the leader list. If the list is empty, an error is returned with success as false
func (r *Replica) send(replicaReq domain.Request) (ok bool, err error) {
	if len(r.leaders) == 0 {
		return false, logger.ErrorWithLine(errors.New(errNoLeader))
	}

	data, err := json.Marshal(replicaReq)
	if err != nil {
		return false, logger.ErrorWithLine(err)
	}

	req, err := http.NewRequest(http.MethodPost, `http://`+r.leaders[0]+domain.RequestLeaderEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return false, logger.ErrorWithLine(err)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return false, logger.ErrorWithLine(err)
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

// Update updates the log of the current replica when a decision is made by the leaders
func (r *Replica) Update(dec domain.Decision) error {
	if dec.SlotID != len(r.log) {
		return logger.ErrorWithLine(errors.New(fmt.Sprintf(`%s (slot: %d, log size: %d)`, errInvalidSlot, dec.SlotID, len(r.log))))
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	r.log = append(r.log, dec.Val)

	return nil
}
