package domain

type Request struct {
	SlotID int    `json:"slot_id"`
	Val    string `json:"value"`
}

type Proposal struct {
	ID  int    `json:"id"`
	Val string `json:"val"`
}
