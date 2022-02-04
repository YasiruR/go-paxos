package roles

import (
	"github.com/go-paxos/domain"
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
func (r *Replica) HandleRequest(req domain.Request) {

}

func (r *Replica) send() {

}
