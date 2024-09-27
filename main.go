package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/yoheinbb/healthd/internal/util"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	// log file setting & Signal handdle
	// todo: from config
	l := &lumberjack.Logger{Filename: "healthd.log", MaxSize: 100, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
	}
	log.SetOutput(l)
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-c
			if err := l.Rotate(); err != nil {
				fmt.Printf("%v", err)
			}
		}
	}()

	fmt.Println("############################")
	fmt.Println("## Start Healthd!! ##")
	fmt.Println("############################")
	fmt.Println("")

	fmt.Println("## Get Args ##")
	cmd_arg := util.ReadCommandArg()
	fmt.Println("global_config_file_path : " + *cmd_arg.GlobalConfigFile)
	fmt.Println("script_config_file_path : " + *cmd_arg.ScriptConfigFile)
	fmt.Println("")

	fmt.Printf("## Read Global Config %s  ##\n", *cmd_arg.GlobalConfigFile)
	gconfig, err := util.NewGlobalConfig(*cmd_arg.GlobalConfigFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println("GlobalConfig Setting")
	fmt.Println("  Port        : " + gconfig.Port)
	fmt.Println("  URLPath     : " + gconfig.URLPath)

	fmt.Printf("## Read Script Config %s  ##\n", *cmd_arg.ScriptConfigFile)
	config, err := util.NewScriptConfig(*cmd_arg.ScriptConfigFile)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Println("ScriptConfig Setting")
	fmt.Println("  Script          : " + config.Script)
	fmt.Println("  MaintenanceFile : " + config.MaintenanceFile)
	fmt.Println("  CheckInterval   : " + config.Interval)
	fmt.Println("  CommandTimeout  : " + config.Timeout)
	fmt.Println("")

	checkInterval, _ := time.ParseDuration(config.Interval)
	// Statusを保持する変数
	ss := util.NewServiceStatus(gconfig)
	// scriptをバックグラウンドでcheckInterval間隔で実行
	go func() {
		for {
			getStatus(ss, config)
			time.Sleep(checkInterval)
		}
	}()

	// Statusを返却するHttpServerの起動
	hs := util.NewHttpServer(ss, gconfig)

	fmt.Println("exec curl from other console:  `curl localhost" + gconfig.Port + gconfig.URLPath + "`")
	hs.Start()

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
	// fmt.Printf("exit code : %d, script path : %s", statusCode, script)
	// fmt.Printf("status    : %s", ss.Status)
}

func checkFileStatus(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
