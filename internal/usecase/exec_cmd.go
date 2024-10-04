package usecase

import (
	"context"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yoheinbb/healthd/internal/util"
)

type ExecCmd struct {
	Status  *Status
	Sconfig *util.ScriptConfig
}

func NewExecCmd(status *Status, sconfig *util.ScriptConfig) *ExecCmd {
	return &ExecCmd{
		Status:  status,
		Sconfig: sconfig,
	}
}

// scriptをバックグラウンドでcheckInterval間隔で実行
// Statusメンバ変数を更新する
func (ecs *ExecCmd) Start(ctx context.Context) error {

	interval, err := (strconv.Atoi(strings.Replace(ecs.Sconfig.Interval, "s", "", -1)))
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

func (ecs *ExecCmd) updateStatus() error {
	// check maintenance file
	if checkFileStatus(ecs.Sconfig.MaintenanceFile) {
		ecs.Status.SetFailed()
		log.Println("maintenance file exits : " + ecs.Sconfig.MaintenanceFile)
		return nil
	}

	// execute command
	cmdTimeout, _ := time.ParseDuration(ecs.Sconfig.Timeout)

	statusCode, err := util.ExecCommand(ecs.Sconfig.Script, cmdTimeout)
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
