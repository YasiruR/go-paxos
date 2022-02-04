package domain

type Request struct {
	SlotID int    `json:"slot_id"`
	Val    string `json:"value"`
}

type Proposal struct {
	PID int    `json:"pid"`
	Val string `json:"val"`
}
