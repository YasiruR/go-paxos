package roles

import (
	"net/http"
)

type Replica struct {
	log     [][]int
	leaders []string
	client  *http.Client
}

func NewReplica() *Replica {
	// add leaders
	return &Replica{}
}

// HandleRequest forwards the request received from a client to a leader
func (r *Replica) HandleRequest(val string) {
	// prepare request

	// forward to a leader

	// when response returned,
	// output slot_id, val when response is received (or all at the end)
	// if not success, retry with incrementing slot id

	// return when the requested val is chosen
}

func (r *Replica) send() {

}
