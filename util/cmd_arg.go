package util

import (
	"flag"
)

type CmdArg struct {
	GlobalConfigFile *string
	ScriptConfigFile *string
}

// CommandLineオプションの読み込み
func ReadCommandArg() *CmdArg {
	cmdarg := new(CmdArg)
	cmdarg.GlobalConfigFile = flag.String("global-config", "conf/global.json", "global-config: json global-config for http port, url path.")
	cmdarg.ScriptConfigFile = flag.String("script-config", "conf/script.json", "script-config: json sciprt config for script path and script exec settings.")
	flag.Parse()

	return cmdarg
}
