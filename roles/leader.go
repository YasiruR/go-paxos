package roles

type Leader struct{}

func NewLeader() *Leader {
	return &Leader{}
}

func (l *Leader) InitProposal() {

}

func (l *Leader) sendPrepare() {

}

func (l *Leader) HandlePrepare() {

}

func (l *Leader) sendAccept() {

}

func (l *Leader) HandleAccept() {

}