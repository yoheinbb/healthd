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
	Status          string
	RetFailed       string
	RetSuccess      string
	Interval        string
	Timeout         string
	MaintenanceFile string
	Script          string
}

func NewServiceStatus(retFailed, retSuccess, interval, timeout, maintenanceFile, script string) *ServiceStatus {
	return &ServiceStatus{
		RetFailed:       retFailed,
		RetSuccess:      retSuccess,
		Interval:        interval,
		Timeout:         timeout,
		MaintenanceFile: maintenanceFile,
		Script:          script,
	}
}

func (ss *ServiceStatus) SetMaintenance() {
	ss.Status = ss.RetFailed
}
func (ss *ServiceStatus) SetInservice() {
	ss.Status = ss.RetSuccess
}

// scriptをバックグラウンドでcheckInterval間隔で実行
// Statusメンバ変数を更新する
func (ss *ServiceStatus) Start(ctx context.Context) error {
	interval, err := (strconv.Atoi(strings.Replace(ss.Interval, "s", "", -1)))
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
	script := ss.Script
	maintenance_file := ss.MaintenanceFile
	cmdTimeout, _ := time.ParseDuration(ss.Timeout)

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
