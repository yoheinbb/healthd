package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/presentation"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util"
	"golang.org/x/sync/errgroup"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	fmt.Println("## Get Args ##")
	cmd_arg := util.ReadCommandArg()
	fmt.Printf("## Read Global Config %s  ##\n", *cmd_arg.GlobalConfigFile)
	gconfig, err := util.NewGlobalConfig(*cmd_arg.GlobalConfigFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Printf("## Read Script Config %s  ##\n", *cmd_arg.ScriptConfigFile)
	sconfig, err := util.NewScriptConfig(*cmd_arg.ScriptConfigFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println("############################")
	fmt.Println("###   Start Healthd!!!   ###")
	fmt.Println("############################")
	fmt.Println("")

	fmt.Println("global_config_file_path : " + *cmd_arg.GlobalConfigFile)
	fmt.Println("script_config_file_path : " + *cmd_arg.ScriptConfigFile)
	fmt.Println("")

	fmt.Println("GlobalConfig Setting")
	fmt.Println("  Port        : " + gconfig.Port)
	fmt.Println("  URLPath     : " + gconfig.URLPath)

	fmt.Println("ScriptConfig Setting")
	fmt.Println("  Script          : " + sconfig.Script)
	fmt.Println("  MaintenanceFile : " + sconfig.MaintenanceFile)
	fmt.Println("  CheckInterval   : " + sconfig.Interval)
	fmt.Println("  CommandTimeout  : " + sconfig.Timeout)
	fmt.Println("")

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer done()

	eg, gctx := errgroup.WithContext(ctx)

	// todo: from config
	l := &lumberjack.Logger{Filename: "healthd.log", MaxSize: 100, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
	}
	log.SetOutput(l)

	// signal channel for SIGHUP
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP)
	// start log roate goroutine
	eg.Go(func() error {
		for {
			select {
			case <-sig:
				log.Println("log rotate")
				if err := l.Rotate(); err != nil {
					fmt.Printf("%v", err)
				}
			case <-gctx.Done():
				if !errors.Is(gctx.Err(), context.Canceled) {
					return gctx.Err()
				}
				return nil
			}
		}
	})

	// Statusを保持
	status := usecase.NewStatus(domain.NewStatus())
	// コマンドを実施しstatusに保持
	ecs := usecase.NewExecCmd(status, sconfig)
	// start getStatus goroutine
	eg.Go(func() error {
		if err := ecs.Start(gctx); err != nil {
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
	handler := presentation.NewHandler(status, gconfig)
	restAPIServer, err := presentation.NewRestAPIServer(handler, gconfig)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	// start httServer goroutine
	eg.Go(func() error {
		log.Println("HttpServer start")
		fmt.Println("exec curl from other console:  `curl localhost" + gconfig.Port + gconfig.URLPath + "`")
		if err := restAPIServer.ListenAndServe(); err != nil {
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

		if err := restAPIServer.Shutdown(context.Background()); err != nil {
			return err
		}
		return nil
	})

	fmt.Println("all component started")

	if err := eg.Wait(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("exit healthd")
}
