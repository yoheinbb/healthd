package util

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"time"
)

type IConfig interface {
	validateConfig() error
	setDefaultConfig() IConfig
}

type GlobalConfig struct {
	Port       string `json:"port"`
	URLPath    string `json:"urlpath"`
	RetSuccess string `json:"ret_success"`
	RetFailed  string `json:"ret_failed"`
}

func NewGlobalConfig(config_path string) (*GlobalConfig, error) {
	config, err := newConfig(config_path, new(GlobalConfig))
	if err != nil {
		return nil, err
	}
	return config.(*GlobalConfig), err
}

func (conf *GlobalConfig) validateConfig() error {
	var check bool
	// format check
	check, _ = regexp.MatchString(`^:[0-9]{2,5}`, conf.Port)
	if !check {
		return errors.New("port format error   " + conf.Port)
	}
	check, _ = regexp.MatchString(`^/{1}[a-z]`, conf.URLPath)
	if !check {
		return errors.New("urlpath format error   " + conf.URLPath)
	}

	return nil
}
func (conf *GlobalConfig) setDefaultConfig() IConfig {
	if conf.Port == "" {
		conf.Port = ":80"
	} else {
		conf.Port = ":" + conf.Port
	}
	if conf.URLPath == "" {
		conf.URLPath = "/healthcheck"
	} else {
		conf.URLPath = "/" + conf.URLPath
	}
	if conf.RetSuccess == "" {
		conf.RetSuccess = "INSERVICE"
	}
	if conf.RetFailed == "" {
		conf.RetFailed = "MAINITENANCE"
	}
	return conf
}

type ScriptConfig struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	Script          string `json:"script"`
	MaintenanceFile string `json:"maintenance_file"`
	Interval        string `json:"interval"`
	Timeout         string `json:"timeout"`
}

func NewScriptConfig(config_path string) (*ScriptConfig, error) {
	config, err := newConfig(config_path, new(ScriptConfig))
	if err != nil {
		return nil, err
	}
	return config.(*ScriptConfig), err
}

func (conf *ScriptConfig) validateConfig() error {
	// Idが指定されているか
	if conf.Id == "" {
		return errors.New("check id not defined")
	}
	// Scriptが存在するか
	if _, err := exec.LookPath(conf.Script); err != nil {
		return err
	}
	// Interval, TimeoutがDurationに変換可能か
	if conf.Interval == "" {
		conf.Interval = "5s"
	} else if _, err := time.ParseDuration(conf.Interval); err != nil {
		return err
	}
	if conf.Timeout == "" {
		conf.Timeout = "10s"
	} else if _, err := time.ParseDuration(conf.Timeout); err != nil {
		return err
	}
	return nil
}

func (conf *ScriptConfig) setDefaultConfig() IConfig {
	return conf
}

// Json Configの読み込み、Config構造体への展開
func newConfig(config_path string, config IConfig) (IConfig, error) {
	// file exist check
	file, err := os.ReadFile(config_path)
	if err != nil {
		return nil, err
	}

	// json decode
	if err := json.Unmarshal(file, &config); err != nil {
		return nil, err
	}

	// check & set default values
	config = config.setDefaultConfig()

	// validate
	err = config.validateConfig()
	if err != nil {
		return nil, err
	}

	return config, nil
}
