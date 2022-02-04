package domain

type Request struct {
	SlotID int `json:"slot_id"`
	Value  int `json:"value"`
}
