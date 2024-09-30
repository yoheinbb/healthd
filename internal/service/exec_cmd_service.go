package service

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/util"
)

type ExecCmdService struct {
	Status          *domain.Status
	Interval        string
	Timeout         string
	MaintenanceFile string
	Script          string
}

func NewExecCmdService(status *domain.Status, interval, timeout, maintenanceFile, script string) *ExecCmdService {
	return &ExecCmdService{
		Status:          status,
		Interval:        interval,
		Timeout:         timeout,
		MaintenanceFile: maintenanceFile,
		Script:          script,
	}
}

// scriptをバックグラウンドでcheckInterval間隔で実行
// Statusメンバ変数を更新する
func (ecs *ExecCmdService) Start(ctx context.Context) error {

	interval, err := (strconv.Atoi(strings.Replace(ecs.Interval, "s", "", -1)))
	if err != nil {
		return err
	}
	intervalTime := time.Duration(interval)

	ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
	if err := ecs.updateStatus(); err != nil {
		return err
	}
	for {
		select {
		case <-ticker.C:
			ecs.updateStatus()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (ecs *ExecCmdService) updateStatus() error {
	// check maintenance file
	if checkFileStatus(ecs.MaintenanceFile) {
		ecs.Status.SetFailed()
		log.Println("maintenance file exits : " + ecs.MaintenanceFile)
		return nil
	}

	// execute command
	cmdTimeout, _ := time.ParseDuration(ecs.Timeout)

	statusCode, err := util.ExecCommand(ecs.Script, cmdTimeout)
	if err != nil {
		if !strings.Contains(err.Error(), "timeout") {
			return err
		}
		log.Printf("%v", err)
	}

	if statusCode == 0 {
		ecs.Status.SetSucess()
	} else {
		ecs.Status.SetFailed()
	}
	// fmt.Printf("exit code : %d, script path : %s\n", statusCode, ecs.Script)
	// fmt.Printf("status    : %s\n", ecs.Status.GetStatus())
	return nil
}

func checkFileStatus(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
