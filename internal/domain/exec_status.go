package domain

type CmdExecStatus struct {
	Status string
}

const (
	Failed      = "FAILED"
	Success     = "SUCCESS"
	Maintenance = "MAINTENANCE"
)

func NewCmdExecStatus() *CmdExecStatus {
	return &CmdExecStatus{Status: Maintenance}

}

func (cs *CmdExecStatus) SetFailed() {
	cs.Status = Failed
}

func (cs *CmdExecStatus) SetSucess() {
	cs.Status = Success
}
