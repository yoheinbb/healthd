package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/yoheinbb/healthd/internal/util"
	"golang.org/x/sync/errgroup"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	// todo: from config
	l := &lumberjack.Logger{Filename: "healthd.log", MaxSize: 100, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
	}
	log.SetOutput(l)

	fmt.Println("## Get Args ##")
	cmd_arg := util.ReadCommandArg()
	fmt.Printf("## Read Global Config %s  ##\n", *cmd_arg.GlobalConfigFile)
	gconfig, err := util.NewGlobalConfig(*cmd_arg.GlobalConfigFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Printf("## Read Script Config %s  ##\n", *cmd_arg.ScriptConfigFile)
	config, err := util.NewScriptConfig(*cmd_arg.ScriptConfigFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	interval, err := (strconv.Atoi(strings.Replace(config.Interval, "s", "", -1)))
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	intervalTime := time.Duration(interval)

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
	fmt.Println("  Script          : " + config.Script)
	fmt.Println("  MaintenanceFile : " + config.MaintenanceFile)
	fmt.Println("  CheckInterval   : " + config.Interval)
	fmt.Println("  CommandTimeout  : " + config.Timeout)
	fmt.Println("")

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	defer done()

	eg, gctx := errgroup.WithContext(ctx)

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

	// Statusを保持する変数
	ss := util.NewServiceStatus(gconfig)
	// scriptをバックグラウンドでcheckInterval間隔で実行
	// start getStatus goroutine
	eg.Go(func() error {
		ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
		getStatus(ss, config)
		for {
			select {
			case <-ticker.C:
				getStatus(ss, config)
			case <-gctx.Done():
				if !errors.Is(gctx.Err(), context.Canceled) {
					return gctx.Err()
				}
				return nil
			}
		}
	})

	// Statusを返却するHttpServerインスタンス生成
	hs := util.NewHttpServer(ss, gconfig)
	// start httServer goroutine
	eg.Go(func() error {
		fmt.Println("exec curl from other console:  `curl localhost" + gconfig.Port + gconfig.URLPath + "`")
		if err := hs.Start(); err != nil {
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

		if err := hs.Shutdown(); err != nil {
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

func getStatus(ss *util.ServiceStatus, config *util.ScriptConfig) {

	script := config.Script
	maintenance_file := config.MaintenanceFile
	cmdTimeout, _ := time.ParseDuration(config.Timeout)

	statusCode, err := util.ExecCommand(script, cmdTimeout)
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
