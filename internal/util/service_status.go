package util

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type ServiceStatus struct {
	Status       string
	GlobalConfig *GlobalConfig
	ScriptConfig *ScriptConfig
}

func NewServiceStatus(gconfig *GlobalConfig, sconfig *ScriptConfig) *ServiceStatus {
	return &ServiceStatus{Status: "MAINTENANCE", GlobalConfig: gconfig, ScriptConfig: sconfig}
}
func (ss *ServiceStatus) SetMaintenance() {
	ss.Status = ss.GlobalConfig.RetFailed
}
func (ss *ServiceStatus) SetInservice() {
	ss.Status = ss.GlobalConfig.RetSuccess
}

// scriptをバックグラウンドでcheckInterval間隔で実行
// Statusメンバ変数を更新する
func (ss *ServiceStatus) Start(ctx context.Context) error {
	interval, err := (strconv.Atoi(strings.Replace(ss.ScriptConfig.Interval, "s", "", -1)))
	if err != nil {
		return err
	}
	intervalTime := time.Duration(interval)

	ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
	ss.updateStatus()
	for {
		select {
		case <-ticker.C:
			ss.updateStatus()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
func (ss *ServiceStatus) updateStatus() {

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
