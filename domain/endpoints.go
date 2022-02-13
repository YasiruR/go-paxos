package domain

const (
	RequestReplicaEndpoint = `/replica/request`
	UpdateReplicaEndpoint  = `/replica/update`
	RequestLeaderEndpoint  = `/leader/request`
	PrepareEndpoint        = `/leader/prepare`
	AcceptEndpoint         = `/leader/accept`
	TermEndpoint           = `/internal/terminate`
)
