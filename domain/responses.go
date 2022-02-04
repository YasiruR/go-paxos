package domain

type Promise struct {
	PID      int
	Success  bool
	Accepted struct {
		ID  int
		Val string
	}
}
