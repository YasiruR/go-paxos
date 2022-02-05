package roles

const (
	idLimit     = 1000
	typePrepare = `prepare`
	typeAccept  = `accept`

	errInvalidSlotLeader = `leader received a request for an invalid slot`
	errBroadcast         = `sending decision to replicas failed`
	errRequestAcceptor   = `received non-2xx code for acceptor response`
	errInvalidProposal   = `acceptor received an older proposal`

	errNoLeader    = `no leader found in the replica`
	errInvalidSlot = `received a decision for an invalid slot`
)
