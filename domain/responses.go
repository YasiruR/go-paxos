package domain

type Decision struct {
	SlotID int    `json:"slot_id"`
	Val    string `json:"val"`
}

type Acceptance struct {
	PID        int `json:"pid"`
	PrvPromise struct {
		Exists bool   `json:"exists"`
		ID     int    `json:"id"`
		Val    string `json:"val"`
	} `json:"prv_promise"`
	PrvAccept struct {
		Exists bool   `json:"exists"`
		ID     int    `json:"id"`
		Val    string `json:"val"`
	} `json:"prv_accept"`

	Accepted bool `json:"accepted"`
}
