package roles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-paxos/domain"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	typePrepare = `prepare`
	typeAccept  = `accept`
)

type Leader struct {
	id       int
	promised struct {
		id  int
		val string
	}
	accepted struct {
		id  int
		val string
	}
	replicas []string
	leaders  []string // except the current one
	client   *http.Client
	lock     *sync.Mutex
}

func NewLeader() *Leader {
	return &Leader{}
}

func (l *Leader) Propose(req domain.Request) (success bool) {
	// create a proposal with p_id=timestamp,leader id
	// send it to other leaders
	// wait for majority responses
	// if there's a negative response, do nothing - return chosen one to replica at the end (wait for a channel?)

	// if all responses are promises, send accept(p_id, val) to each of these acceptors
	// if there's an accepted value, send accept(p_id, accepted_val) - check if this is really required

	// wait for majority responses
	// if all accept, return val to all replicas

	var dec domain.Decision
	prop := l.newProposal(req.Val)
	accepted, rejected, valid := l.validatePromises(l.send(typePrepare, prop))
	if valid {
		if accepted > rejected {
			accepted, rejected = l.validateAccepts(l.send(typeAccept, prop))
			if accepted > rejected {
				dec.SlotID = req.SlotID
				dec.Val = req.Val
				l.broadcastDecision(dec)
				return true
			}
		}
	}

	return false
}

func (l *Leader) broadcastDecision(dec domain.Decision) {
	data, err := json.Marshal(dec)
	if err != nil {
		// err
	}

	for _, replica := range l.replicas {
		// todo can do in parallel
		req, err := http.NewRequest(http.MethodPost, `http://`+replica+domain.UpdateReplicaEndpoint, bytes.NewBuffer(data))
		if err != nil {
			// err
		}

		res, err := l.client.Do(req)
		if err != nil {
			// err
		}

		defer res.Body.Close() // todo close in each return
		// check status
	}
}

func (l *Leader) newProposal(val string) domain.Proposal {
	ts := time.Now().Second()
	pId, err := strconv.Atoi(fmt.Sprintf(`%d%d`, ts, l.id))
	if err != nil {
		// err
	}

	return domain.Proposal{ID: pId, Val: val}
}

func (l *Leader) send(typ string, prop domain.Proposal) []domain.Acceptance {
	data, err := json.Marshal(prop)
	if err != nil {
		// err
	}

	var endpoint string
	if typ == typePrepare {
		endpoint = domain.PrepareEndpoint
	} else {
		endpoint = domain.AcceptEndpoint
	}

	var resList []domain.Acceptance
	for _, acceptor := range l.leaders {
		// todo do this in parallel
		req, err := http.NewRequest(http.MethodPost, `http://`+acceptor+endpoint, bytes.NewBuffer(data))
		if err != nil {
			// err
		}

		// todo majority is enough
		res, err := l.client.Do(req)
		if err != nil {
			// err
		}
		defer res.Body.Close() // todo close in each return

		resData, err := ioutil.ReadAll(res.Body)
		if err != nil {
			// err
		}

		var response domain.Acceptance
		err = json.Unmarshal(resData, &response)
		if err != nil {
			// err
		}
		resList = append(resList, response)
	}

	return resList
}

func (l *Leader) validatePromises(resList []domain.Acceptance) (accepted, rejected int, valid bool) {
	accepted, rejected = 0, 0
	for _, promise := range resList {
		if promise.PrvAccept.Exists {
			// todo close all res bodies in the outer func
			if promise.PrvAccept.ID >= promise.PID {
				return accepted, rejected, false
			}
			rejected++
			continue
		}

		if promise.PrvPromise.Exists {
			if promise.PrvPromise.ID >= promise.PID {
				return accepted, rejected, false
			}
			rejected++
			continue
		}
		accepted++
	}

	return accepted, rejected, true
}

func (l *Leader) validateAccepts(resList []domain.Acceptance) (accepted, rejected int) {
	accepted, rejected = 0, 0
	for _, accept := range resList {
		if accept.Accepted {
			accepted++
			continue
		}
		rejected++
	}

	return accepted, rejected
}

func (l *Leader) HandlePrepare(prop domain.Proposal) domain.Acceptance {
	// proceed if first prepare message
	// if not, check if p_id > accepted_id
	// if not true, send a negative response

	// if true, check if it has already accepted
	// if accepted, return a response with accepted_id, val
	// if not, return promise(p_id)

	var res domain.Acceptance
	res.PID = prop.ID
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.promised.id >= prop.ID {
		res.PrvPromise.Exists = true
		res.PrvPromise.ID = l.promised.id
		res.PrvPromise.Val = l.promised.val
	} else {
		l.promised.id = prop.ID
		l.promised.val = prop.Val
	}

	if l.accepted.id != 0 {
		res.PrvAccept.Exists = true
		res.PrvAccept.ID = l.accepted.id
		res.PrvAccept.Val = l.accepted.val
	}

	return res
}

func (l *Leader) HandleAccept(p domain.Proposal) domain.Acceptance {
	// if p_id > accepted_id, store accepted(p_id, val) and return response accept()
	// if not, reply negative response

	var res domain.Acceptance
	res.PID = p.ID
	l.lock.Lock()
	defer l.lock.Unlock()

	if l.promised.id >= p.ID {
		res.Accepted = false
		return res
	}

	if l.accepted.id != 0 {
		res.Accepted = false
		return res
	}

	l.accepted.id = p.ID
	l.accepted.val = p.Val
	res.Accepted = true

	return res
}
