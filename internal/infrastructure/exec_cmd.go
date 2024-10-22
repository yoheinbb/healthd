package infrastructure

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/yoheinbb/healthd/internal/util"
)

type ExecCmdRepository struct {
	statusCode      int
	maintenanceFile string
	script          string
	timeout         time.Duration
	logger          *slog.Logger
}

func NewExecCmdRepository(maintenanceFile, script, timeout string, logger *slog.Logger) (*ExecCmdRepository, error) {

	ttimeout, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, err
	}
	return &ExecCmdRepository{
		statusCode:      -1,
		maintenanceFile: maintenanceFile,
		script:          script,
		timeout:         ttimeout,
		logger:          logger,
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

	ecr.logger.Debug(fmt.Sprintf("exit code : %d, script path : %s", statusCode, ecr.script))
	ecr.logger.Debug(fmt.Sprintf("status    : %d", ecr.statusCode))
	return nil
}

func checkFileStatus(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
