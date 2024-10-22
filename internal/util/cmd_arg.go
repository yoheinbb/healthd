package util

import (
	"flag"
	"log/slog"
)

type CmdArg struct {
	GlobalConfigFile *string
	ScriptConfigFile *string
	LogLevel         slog.Leveler
}

// CommandLineオプションの読み込み
func ReadCommandArg() *CmdArg {
	cmdarg := new(CmdArg)
	cmdarg.GlobalConfigFile = flag.String("global-config", "conf/global.json", "global-config: json global-config for http port, url path.")
	cmdarg.ScriptConfigFile = flag.String("script-config", "conf/script.json", "script-config: json sciprt config for script path and script exec settings.")
	logLevel := flag.Bool("v", false, "v: verbose logs")
	flag.Parse()

	if *logLevel {
		cmdarg.LogLevel = slog.LevelDebug
	} else {
		cmdarg.LogLevel = slog.LevelInfo
	}

	return cmdarg
}
