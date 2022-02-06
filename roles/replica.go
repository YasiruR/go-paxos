package roles

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-paxos/domain"
	"github.com/go-paxos/logger"
	"github.com/tryfix/log"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type Replica struct {
	hostname string
	log      []string
	leaders  []string
	client   *http.Client
	lock     *sync.Mutex
	logger   log.Logger
}

func NewReplica(hostname string, leaders []string, logger log.Logger) *Replica {
	return &Replica{
		hostname: hostname,
		leaders:  leaders,
		client:   &http.Client{Timeout: time.Duration(domain.Config.ReplicaTimeout) * time.Second},
		lock:     &sync.Mutex{},
		logger:   logger,
	}
}

// HandleRequest builds a request from the client value and forwards the request received to a leader
func (r *Replica) HandleRequest(ctx context.Context, val string) error {
	req := r.buildRequest(val)
	for {
		dec, ok, err := r.send(req)
		if err != nil {
			return logger.ErrorWithLine(err)
		}

		if ok {
			err = r.Update(ctx, dec)
			if err != nil {
				return logger.ErrorWithLine(err)
			}
			return nil
		}

		time.Sleep(1000 * time.Millisecond)
		req.SlotID++
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
func (r *Replica) send(replicaReq domain.Request) (dec domain.Decision, ok bool, err error) {
	if len(r.leaders) == 0 {
		return domain.Decision{}, false, logger.ErrorWithLine(errors.New(errNoLeader))
	}

	data, err := json.Marshal(replicaReq)
	if err != nil {
		return domain.Decision{}, false, logger.ErrorWithLine(err)
	}

	req, err := http.NewRequest(http.MethodPost, `http://`+r.leaders[0]+domain.RequestLeaderEndpoint, bytes.NewBuffer(data))
	if err != nil {
		return domain.Decision{}, false, logger.ErrorWithLine(err)
	}

	res, err := r.client.Do(req)
	if err != nil {
		return domain.Decision{}, false, logger.ErrorWithLine(err)
	}
	defer res.Body.Close()

	resData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return domain.Decision{}, false, logger.ErrorWithLine(err)
	}

	if res.StatusCode != http.StatusOK {
		return domain.Decision{}, false, nil
	}

	err = json.Unmarshal(resData, &dec)
	if err != nil {
		return domain.Decision{}, false, logger.ErrorWithLine(errors.New(fmt.Sprintf(`%s for value %s (res: %s)`, err.Error(), replicaReq.Val, string(resData))))
	}

	return dec, res.StatusCode == http.StatusOK, nil
}

// Update updates the log of the current replica when a decision is made by the leaders
func (r *Replica) Update(ctx context.Context, dec domain.Decision) error {
	if dec.SlotID != len(r.log) {
		return logger.ErrorWithLine(errors.New(fmt.Sprintf(`%s (slot: %d, log size: %d)`, errInvalidSlot, dec.SlotID, len(r.log))))
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	r.log = append(r.log, dec.Val)

	r.logger.DebugContext(ctx, `replica state updated`, r.log)
	return nil
}
