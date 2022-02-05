package roles

import (
	"bytes"
	"encoding/json"
	"github.com/go-paxos/domain"
	"net/http"
	"sync"
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

// HandleRequest forwards the request received from a client to a leader
func (r *Replica) HandleRequest(val string) error {
	// prepare request

	// forward to a leader

	// when response returned,
	// output slot_id, val when response is received (or all at the end)
	// if not success, retry with incrementing slot id

	// return when the requested val is chosen

	req := r.buildRequest(val)
	for {
		ok, err := r.send(req)
		if err != nil {
			// err
		}

		if ok {
			break
		}

		r.lock.Lock()
		if req.SlotID == len(r.log) {
			req.SlotID++
		} else {
			req.SlotID = len(r.log)
		}
		r.lock.Unlock()
	}

	return nil
}

func (r *Replica) buildRequest(val string) domain.Request {
	return domain.Request{
		Replica: r.hostname,
		SlotID:  len(r.log),
		Val:     val,
	}
}

func (r *Replica) send(replicaReq domain.Request) (bool, error) {
	if len(r.leaders) == 0 {
		// err
		//return false, err
	}

	data, err := json.Marshal(replicaReq)
	if err != nil {
		// err
	}

	req, err := http.NewRequest(http.MethodPost, `http://`+r.leaders[0]+domain.RequestLeaderEndpoint, bytes.NewBuffer(data))
	if err != nil {
		// err
	}

	res, err := r.client.Do(req)
	if err != nil {
		// err
	}
	defer res.Body.Close()

	return res.StatusCode == http.StatusOK, nil
}

func (r *Replica) Update(dec domain.Decision) error {
	if dec.SlotID != len(r.log) {
		// err
	}

	r.lock.Lock()
	defer r.lock.Unlock()
	r.log = append(r.log, dec.Val)

	return nil
}
