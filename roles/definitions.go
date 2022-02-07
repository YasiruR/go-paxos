package roles

const (
	idLimit     = 1000
	typePrepare = `prepare`
	typeAccept  = `accept`

	errFutureSlot      = `received future slot`
	errBroadcast       = `sending decision to replicas failed`
	errRequestAcceptor = `received non-2xx code for acceptor response`
	errInvalidProposal = `acceptor received an older proposal`

	errNoLeader        = `no leader found in the replica`
	errInvalidDecision = `received a decision for an invalid slot`
)
