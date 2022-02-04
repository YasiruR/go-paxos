package roles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-paxos/domain"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Leader struct {
	id       int
	prepared struct {
		id  int
		val int
	}
	accepted struct {
		id  int
		val int
	}
	replicas []string
	leaders  []string // except the current one
	client   *http.Client
}

func NewLeader() *Leader {
	return &Leader{}
}

func (l *Leader) InitProposal(req domain.Request) {
	// create a proposal with p_id=timestamp,leader id
	// send it to other leaders
	// wait for majority responses
	// if there's a negative response, do nothing - return chosen one to replica at the end (wait for a channel?)

	// if all responses are promises, send accept(p_id, val) to each of these acceptors
	// if there's an accepted value, send accept(p_id, accepted_val) - check if this is really required

	// wait for majority responses
	// if all accept, return val to all replicas
}

func (l *Leader) newProposal(val string) domain.Proposal {
	ts := time.Now().Second()
	pId, err := strconv.Atoi(fmt.Sprintf(`%d%d`, ts, l.id))
	if err != nil {
		// err
	}

	return domain.Proposal{PID: pId, Val: val}
}

func (l *Leader) sendPrepare(prop domain.Proposal) []*http.Response {
	data, err := json.Marshal(prop)
	if err != nil {
		// err
	}

	var resList []*http.Response
	for _, acceptor := range l.leaders {
		// todo do this in parallel
		req, err := http.NewRequest(http.MethodPost, `http://`+acceptor+domain.PrepareEndpoint, bytes.NewBuffer(data))
		if err != nil {
			// err
		}

		// todo majority is enough
		res, err := l.client.Do(req)
		if err != nil {
			// err
		}
		resList = append(resList, res)
	}

	return resList
}

func (l *Leader) HandlePrepare() {
	// proceed if first prepare message
	// if not, check if p_id > accepted_id
	// if not true, send a negative response

	// if true, check if it has already accepted
	// if accepted, return a response with accepted_id, val
	// if not, return promise(p_id)
}

func (l *Leader) validatePromises(resList []*http.Response) bool {
	// check if res can be passed without closing
	for _, res := range resList {
		data, err := ioutil.ReadAll(res.Body)
		if err != nil {
			// err
		}

		var promise domain.Promise
		err = json.Unmarshal(data, &promise)
		if err != nil {
			// err
		}

		if !promise.Success {
			// todo close all res bodies in the outer func
			return false
		}
	}

	return true
}

func (l *Leader) sendAccept() {

}

func (l *Leader) HandleAccept() {
	// if p_id > accepted_id, store accepted(p_id, val) and return response accept()
	// if not, reply negative response
}
