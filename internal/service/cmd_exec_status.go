package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/util"
)

type CmdExecStatusService struct {
	CmdExecStatus   domain.CmdExecStatus
	Interval        string
	Timeout         string
	MaintenanceFile string
	Script          string
}

func NewCmdExecStatusService(interval, timeout, maintenanceFile, script string) *CmdExecStatusService {
	return &CmdExecStatusService{
		Interval:        interval,
		Timeout:         timeout,
		MaintenanceFile: maintenanceFile,
		Script:          script,
	}
}

// scriptをバックグラウンドでcheckInterval間隔で実行
// Statusメンバ変数を更新する
func (css *CmdExecStatusService) Start(ctx context.Context) error {
	interval, err := (strconv.Atoi(strings.Replace(css.Interval, "s", "", -1)))
	if err != nil {
		return err
	}
	intervalTime := time.Duration(interval)

	ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
	css.updateStatus()
	for {
		select {
		case <-ticker.C:
			css.updateStatus()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (css *CmdExecStatusService) updateStatus() {
	cmdTimeout, _ := time.ParseDuration(css.Timeout)

	statusCode, err := util.ExecCommand(css.Script, cmdTimeout)
	if err != nil {
		fmt.Printf("%s", err)
		log.Printf("%s", err)
	}

	if checkFileStatus(css.MaintenanceFile) {
		css.CmdExecStatus.SetFailed()
		log.Println("maintenance file exits : " + css.MaintenanceFile)
	} else if statusCode == 0 {
		css.CmdExecStatus.SetSucess()
	} else {
		css.CmdExecStatus.SetFailed()
	}
	//fmt.Printf("exit code : %d, script path : %s\n", statusCode, script)
	//fmt.Printf("status    : %s\n", ss.Status)
}

func checkFileStatus(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
