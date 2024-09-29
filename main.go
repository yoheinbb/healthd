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
	sconfig, err := util.NewScriptConfig(*cmd_arg.ScriptConfigFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	interval, err := (strconv.Atoi(strings.Replace(sconfig.Interval, "s", "", -1)))
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
	fmt.Println("  Script          : " + sconfig.Script)
	fmt.Println("  MaintenanceFile : " + sconfig.MaintenanceFile)
	fmt.Println("  CheckInterval   : " + sconfig.Interval)
	fmt.Println("  CommandTimeout  : " + sconfig.Timeout)
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
	ss := util.NewServiceStatus(gconfig, sconfig)
	// scriptをバックグラウンドでcheckInterval間隔で実行
	// start getStatus goroutine
	eg.Go(func() error {
		ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
		ss.GetStatus()
		for {
			select {
			case <-ticker.C:
				ss.GetStatus()
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
