package util

type ServiceStatus struct {
    Status string
    GlobalConfig *GlobalConfig
}
func (self *ServiceStatus) SetMaintenance() {
    self.Status = self.GlobalConfig.RetFailed
}
func (self *ServiceStatus) SetInservice() {
    self.Status = self.GlobalConfig.RetSuccess
}
func NewServiceStatus(gconfig *GlobalConfig) *ServiceStatus {
    return &ServiceStatus{ Status: "MAINTENANCE", GlobalConfig: gconfig }
}
