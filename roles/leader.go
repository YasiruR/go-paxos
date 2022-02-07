package roles

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-paxos/domain"
	"github.com/go-paxos/logger"
	"github.com/tryfix/log"
	"io/ioutil"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type SlotStatus int

const (
	ValidSlot = iota + 1
	FutureSlot
	InvalidSlot
)

// internal state structure of last promised and accepted proposals
type state struct {
	id   int
	slot int
	val  string
}

type Leader struct {
	id       int
	lastSlot int
	promised state
	accepted state
	leaders  []string // excluding the current node
	replicas []string
	client   *http.Client
	lock     *sync.RWMutex
	logger   log.Logger
}

func NewLeader(hostname string, leaders, replicas []string, logger log.Logger) *Leader {
	return &Leader{
		id:       id(hostname),
		lastSlot: -1,
		promised: state{},
		accepted: state{},
		leaders:  leaders,
		replicas: replicas,
		client:   &http.Client{Timeout: time.Duration(domain.Config.LeaderTimeout) * time.Second},
		lock:     &sync.RWMutex{},
		logger:   logger,
	}
}

func id(hostname string) int {
	sum := sha256.Sum256([]byte(hostname))
	hexVal := hex.EncodeToString(sum[:])
	n := new(big.Int)
	n.SetString(hexVal, 16)

	return int(n.Uint64() % idLimit)
}

/* Proposer functions */

func (l *Leader) ValidateSlot(reqSlot int) (lastSlot int, status SlotStatus) {
	l.lock.RLock()
	defer l.lock.RUnlock()
	if reqSlot > l.lastSlot+1 {
		return l.lastSlot, FutureSlot
	}

	if l.lastSlot+1 != reqSlot {
		return l.lastSlot, InvalidSlot
	}

	return l.lastSlot, ValidSlot
}

// Propose creates the proposal when a replica has requested this leader and carries out the consensus algorithm
func (l *Leader) Propose(ctx context.Context, req domain.Request) (dec domain.Decision, ok bool, err error) {
	prop, err := l.newProposal(req.SlotID, req.Val)
	if err != nil {
		return domain.Decision{}, false, logger.ErrorWithLine(err)
	}

	resList, err := l.send(ctx, typePrepare, prop)
	if err != nil {
		return domain.Decision{}, false, logger.ErrorWithLine(err)
	}

	accepted, rejected, valid := l.validatePromises(resList)
	if valid {
		if accepted > rejected {
			resList, err = l.send(ctx, typeAccept, prop)
			if err != nil {
				return domain.Decision{}, false, logger.ErrorWithLine(err)
			}

			accepted, rejected = l.validateAccepts(resList)
			if accepted > rejected {
				l.logger.DebugContext(ctx, fmt.Sprintf(`requested value %s was proposed and chosen for slot %d`, req.Val, req.SlotID))
				dec.SlotID = req.SlotID
				dec.Val = req.Val

				l.lock.Lock()
				l.lastSlot++
				l.lock.Unlock()

				err = l.broadcastDecision(dec, req.Replica)
				if err != nil {
					return domain.Decision{}, false, logger.ErrorWithLine(err)
				}
				return dec, true, nil
			}
		}
	}

	return domain.Decision{}, false, nil
}

// newProposal creates a proposal with an id in the format of `timestamp`+`leader_id`
func (l *Leader) newProposal(slotID int, val string) (domain.Proposal, error) {
	ts := time.Now().Unix()
	pId, err := strconv.Atoi(fmt.Sprintf(`%d%d`, ts, l.id))
	if err != nil {
		return domain.Proposal{}, logger.ErrorWithLine(err)
	}

	return domain.Proposal{ID: pId, SlotID: slotID, Val: val}, nil
}

// Broadcasts the decision to all the replicas excluding the requested one
func (l *Leader) broadcastDecision(dec domain.Decision, requester string) error {
	data, err := json.Marshal(dec)
	if err != nil {
		return logger.ErrorWithLine(err)
	}

	for _, replica := range l.replicas {
		if replica == requester {
			continue
		}

		// todo can do in parallel
		req, err := http.NewRequest(http.MethodPost, `http://`+replica+domain.UpdateReplicaEndpoint, bytes.NewBuffer(data))
		if err != nil {
			return logger.ErrorWithLine(err)
		}

		res, err := l.client.Do(req)
		if err != nil {
			return logger.ErrorWithLine(err)
		}

		if res.StatusCode != http.StatusOK {
			res.Body.Close()
			return logger.ErrorWithLine(errors.New(fmt.Sprintf(`%s (status: %d)`, errBroadcast, res.StatusCode)))
		}
		res.Body.Close()
	}

	return nil
}

// Sends out the proposal to all acceptors in both phases prepare and accept, excluding the current leader as it does not exist in leader list
func (l *Leader) send(ctx context.Context, typ string, prop domain.Proposal) ([]domain.Acceptance, error) {
	data, err := json.Marshal(prop)
	if err != nil {
		return nil, logger.ErrorWithLine(err)
	}

	var endpoint string
	if typ == typePrepare {
		endpoint = domain.PrepareEndpoint
	} else {
		endpoint = domain.AcceptEndpoint
	}

	var resList []domain.Acceptance
	resChan := make(chan domain.Acceptance)
	errChan := make(chan error)
	wg := &sync.WaitGroup{}
	for _, acceptor := range l.leaders {
		wg.Add(1)
		go func(acceptor string, resChan chan domain.Acceptance, wg *sync.WaitGroup, errChan chan error) {
			req, err := http.NewRequest(http.MethodPost, `http://`+acceptor+endpoint, bytes.NewBuffer(data))
			if err != nil {
				errChan <- errors.New(fmt.Sprintf(`%s for acceptor: %s`, err.Error(), acceptor))
				return
			}

			res, err := l.client.Do(req)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf(`%s for acceptor: %s`, err.Error(), acceptor))
				return
			}

			if res.StatusCode != http.StatusOK {
				res.Body.Close()
				errChan <- errors.New(fmt.Sprintf(`%s (type: %s, status: %d) for acceptor: %s`, errRequestAcceptor, typ, res.StatusCode, acceptor))
				return
			}

			resData, err := ioutil.ReadAll(res.Body)
			if err != nil {
				res.Body.Close()
				errChan <- errors.New(fmt.Sprintf(`%s for acceptor: %s`, err.Error(), acceptor))
				return
			}
			res.Body.Close()

			var response domain.Acceptance
			err = json.Unmarshal(resData, &response)
			if err != nil {
				errChan <- errors.New(fmt.Sprintf(`%s for acceptor: %s`, err.Error(), acceptor))
				return
			}
			wg.Done()
			resChan <- response
		}(acceptor, resChan, wg, errChan)
	}

	wg.Wait()
	for i := 0; i < len(l.leaders); i++ {
		select {
		case res := <-resChan:
			resList = append(resList, res)
		case err = <-errChan:
			l.logger.ErrorContext(ctx, err)
		default:
			continue
		}
	}

	return resList, nil
}

// Validates promises upon receiving them from acceptors and returns number of accepted and rejected cases. This function
// returns false for valid if a different proposer has already started a proposal with a higher id.
func (l *Leader) validatePromises(resList []domain.Acceptance) (accepted, rejected int, valid bool) {
	accepted, rejected = 0, 0
	for _, promise := range resList {
		if promise.PrvAccept.Exists {
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

// Validates accept responses and returns the accepted and rejected cases
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

/* Acceptor functions */

// HandlePrepare handles prepare message requested by a proposer to check if this acceptor has already promised or accepted a proposal
func (l *Leader) HandlePrepare(prop domain.Proposal) (domain.Acceptance, error) {
	var res domain.Acceptance
	res.PID = prop.ID
	l.lock.Lock()
	defer l.lock.Unlock()

	// returns an error if the proposal is for an older slot
	if l.accepted.slot > prop.SlotID {
		return domain.Acceptance{}, logger.ErrorWithLine(errors.New(fmt.Sprintf(`%s (phase: %s, last: %d, requested: %d)`,
			errInvalidProposal, typePrepare, l.accepted.slot, prop.SlotID)))
	}

	if l.promised.slot == prop.SlotID {
		// check if promised id is higher than the requested one since proposer will use this to terminate its proposal
		if l.promised.id >= prop.ID {
			res.PrvPromise.Exists = true
			res.PrvPromise.ID = l.promised.id
			res.PrvPromise.Val = l.promised.val
		} else {
			// as the requested prepare is valid, acceptor updates its state for the same slot
			l.promised.id = prop.ID
			l.promised.val = prop.Val
		}
	} else {
		// if the prepare request is for a new slot
		l.promised.id = prop.ID
		l.promised.slot = prop.SlotID
		l.promised.val = prop.Val
	}

	// if there's an already accepted proposal for the same slot, acceptor just notifies the proposer
	if l.accepted.slot == prop.SlotID && l.accepted.id != 0 {
		res.PrvAccept.Exists = true
		res.PrvAccept.ID = l.accepted.id
		res.PrvAccept.Val = l.accepted.val
	}

	return res, nil
}

// HandleAccept checks if it can accept the confirmation request from a proposer
func (l *Leader) HandleAccept(prop domain.Proposal) (domain.Acceptance, error) {
	// returns an error if the proposal is for an older slot
	if l.accepted.slot > prop.SlotID {
		return domain.Acceptance{}, logger.ErrorWithLine(errors.New(fmt.Sprintf(`%s (phase: %s, last: %d, requested: %d)`,
			errInvalidProposal, typeAccept, l.accepted.slot, prop.SlotID)))
	}

	var res domain.Acceptance
	res.PID = prop.ID
	l.lock.Lock()
	defer l.lock.Unlock()

	// rejects if already promised to a proposal with a higher id for the same slot
	if l.promised.slot == prop.SlotID && l.promised.id > prop.ID {
		res.Accepted = false
		return res, nil
	}

	// rejects if already accepted for the same slot
	if l.accepted.slot == prop.SlotID && l.accepted.id != 0 {
		res.Accepted = false
		return res, nil
	}

	l.accepted.id = prop.ID
	l.accepted.val = prop.Val
	l.accepted.slot = prop.SlotID
	l.lastSlot++
	res.Accepted = true

	return res, nil
}
