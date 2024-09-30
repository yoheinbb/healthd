package domain

type Status struct {
	s string
}

const (
	Failed      = "FAILED"
	Success     = "SUCCESS"
	Maintenance = "MAINTENANCE"
)

func NewStatus() *Status {
	return &Status{s: Maintenance}

}

func (cs *Status) SetFailed() {
	cs.s = Failed
}

func (cs *Status) SetSucess() {
	cs.s = Success
}

func (cs *Status) GetStatus() string {
	return cs.s
}
