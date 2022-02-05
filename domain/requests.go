package domain

type Request struct {
	Replica string `json:"replica"`
	SlotID  int    `json:"slot_id"`
	Val     string `json:"value"`
}

type Proposal struct {
	ID     int    `json:"id"`
	SlotID int    `json:"slot_id"`
	Val    string `json:"val"`
}
