package infrastructure

import (
	"log"
	"os"
	"time"

	"github.com/yoheinbb/healthd/internal/util"
)

type ExecCmdRepository struct {
	statusCode      int
	maintenanceFile string
	script          string
	timeout         time.Duration
}

func NewExecCmdRepository(maintenanceFile, script, timeout string) (*ExecCmdRepository, error) {

	ttimeout, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	return &ExecCmdRepository{
		statusCode:      -1,
		maintenanceFile: maintenanceFile,
		script:          script,
		timeout:         ttimeout,
	}, nil
}

func (ecr *ExecCmdRepository) GetStatus() int {
	return ecr.statusCode
}

func (ecr *ExecCmdRepository) UpdateStatus() error {
	// check maintenance file
	if checkFileStatus(ecr.maintenanceFile) {
		ecr.statusCode = 1
		log.Println("maintenance file exits : " + ecr.maintenanceFile)
		return nil
	}

	statusCode, err := util.ExecCommand(ecr.script, ecr.timeout)
	if err != nil {
		return err
	}
	ecr.statusCode = statusCode

	// fmt.Printf("exit code : %d, script path : %s\n", statusCode, ecs.Script)
	// fmt.Printf("status    : %s\n", ecs.Status.GetStatus())
	return nil
}

func checkFileStatus(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
