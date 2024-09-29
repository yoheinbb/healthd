package util

import (
	"fmt"
	"log"
	"os"
	"time"
)

type ServiceStatus struct {
	Status       string
	GlobalConfig *GlobalConfig
	ScriptConfig *ScriptConfig
}

func (ss *ServiceStatus) SetMaintenance() {
	ss.Status = ss.GlobalConfig.RetFailed
}
func (ss *ServiceStatus) SetInservice() {
	ss.Status = ss.GlobalConfig.RetSuccess
}
func NewServiceStatus(gconfig *GlobalConfig, sconfig *ScriptConfig) *ServiceStatus {
	return &ServiceStatus{Status: "MAINTENANCE", GlobalConfig: gconfig, ScriptConfig: sconfig}
}

func (ss *ServiceStatus) GetStatus() {

	script := ss.ScriptConfig.Script
	maintenance_file := ss.ScriptConfig.MaintenanceFile
	cmdTimeout, _ := time.ParseDuration(ss.ScriptConfig.Timeout)

	statusCode, err := ExecCommand(script, cmdTimeout)
	if err != nil {
		fmt.Printf("%s", err)
		log.Printf("%s", err)
	}

	if checkFileStatus(maintenance_file) {
		ss.SetMaintenance()
		log.Println("maintenance file exits : " + maintenance_file)
	} else if statusCode == 0 {
		ss.SetInservice()
	} else {
		ss.SetMaintenance()
	}
	//fmt.Printf("exit code : %d, script path : %s\n", statusCode, script)
	//fmt.Printf("status    : %s\n", ss.Status)
}

func checkFileStatus(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
