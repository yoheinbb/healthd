package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/infrastructure"
	"github.com/yoheinbb/healthd/internal/presentation"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util"
	"golang.org/x/sync/errgroup"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	logger.Info("## Get Args ##")
	cmd_arg := util.ReadCommandArg()

	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cmd_arg.LogLevel}))
	slog.SetDefault(logger)

	logger.Info(fmt.Sprintf("## Read Global Config %s  ##\n", *cmd_arg.GlobalConfigFile))
	gconfig, err := util.NewGlobalConfig(*cmd_arg.GlobalConfigFile)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("## Read Script Config %s  ##\n", *cmd_arg.ScriptConfigFile))
	sconfig, err := util.NewScriptConfig(*cmd_arg.ScriptConfigFile)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	startMessageTemplate := `
	############################
	###   Start Healthd!!!   ###
	############################

	global_config_file_path : %s
	script_config_file_path : %s

	GlobalConfig Setting
	  Port        : %s
	  URLPath     : %s

	ScriptConfig Setting
	  Script          : %s
	  MaintenanceFile : %s
	  CheckInterval   : %s
	  CommandTimeout  : %s
	`
	startMessage := fmt.Sprintf(startMessageTemplate,
		*cmd_arg.GlobalConfigFile, *cmd_arg.ScriptConfigFile,
		gconfig.Port, gconfig.URLPath,
		sconfig.Script, sconfig.MaintenanceFile, sconfig.Interval, sconfig.Timeout,
	)
	logger.Debug(startMessage)

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer done()

	eg, gctx := errgroup.WithContext(ctx)

	// domain
	d := domain.NewStatus()
	// cmd exec repository
	r, err := infrastructure.NewExecCmdRepository(sconfig.MaintenanceFile, sconfig.Script, sconfig.Timeout, logger)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	// usecase for status
	usecase := usecase.NewStatus(d, r)
	// start getStatus goroutine
	eg.Go(func() error {
		interval, err := (strconv.Atoi(strings.Replace(sconfig.Interval, "s", "", -1)))
		if err != nil {
			return err
		}

		if err := usecase.StartStatusUpdater(gctx, interval); err != nil {
			if !errors.Is(err, context.Canceled) {
				return err
			}
		}

		<-gctx.Done()
		if !errors.Is(gctx.Err(), context.Canceled) {
			return gctx.Err()
		}
		return nil
	})

	// Statusを返却するHttpServerインスタンス生成
	apiServer, err := presentation.NewAPIServer(
		logger,
		gconfig.URLPath,
		gconfig.Port,
		presentation.NewHandler(usecase, gconfig.RetSuccess, gconfig.RetFailed),
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// start httServer goroutine
	eg.Go(func() error {
		logger.Info("HttpServer start")
		logger.Info(fmt.Sprintf("exec curl from other console:  `curl localhost" + gconfig.Port + gconfig.URLPath + "`"))
		if err := apiServer.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				return err
			}
		}

		<-gctx.Done()
		if !errors.Is(gctx.Err(), context.Canceled) {
			return gctx.Err()
		}
		return nil
	})

	// signalを受けたらhttp serverを停止する
	// shutdown httServer goroutine
	eg.Go(func() error {
		<-gctx.Done()
		if !errors.Is(gctx.Err(), context.Canceled) {
			return gctx.Err()
		}

		if err := apiServer.Shutdown(context.Background()); err != nil {
			return err
		}
		return nil
	})

	logger.Info("all component started")

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	logger.Info("exit healthd")
}
