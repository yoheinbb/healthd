package util

type ServiceStatus struct {
	Status       string
	GlobalConfig *GlobalConfig
}

func (ss *ServiceStatus) SetMaintenance() {
	ss.Status = ss.GlobalConfig.RetFailed
}
func (ss *ServiceStatus) SetInservice() {
	ss.Status = ss.GlobalConfig.RetSuccess
}
func NewServiceStatus(gconfig *GlobalConfig) *ServiceStatus {
	return &ServiceStatus{Status: "MAINTENANCE", GlobalConfig: gconfig}
}
