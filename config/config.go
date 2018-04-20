package config

import (
	"flag"
	"strconv"
	"strings"
)

type CONFIG_MODE int32

const (
	INI CONFIG_MODE = iota
	JSON
	ENV
)

const _DEFAULT_CONFIG = "default_config"

type ConfigPool struct {
	Pool map[string]Configer
}

type Configer interface {
	Get(string) string
	Set(string, string)
	Bool(string) bool
	Int(string) int
	Int64(string) int64
	Float(string) float64
	String(string) string
}

var Config *ConfigPool

func init() {
	if Config != nil {
		return
	}
	Config = &ConfigPool{}
	Config.Pool = make(map[string]Configer)
}

func InitDefault(mode CONFIG_MODE, operationParam, defaultValue, tips string) Configer {
	confFile := flag.String(operationParam, defaultValue, tips)
	flag.Parse()
	if *confFile == "" {
		panic("run with -h to find usage")
	}
	AddConfiger(mode, _DEFAULT_CONFIG, *confFile)
	return Default()
}

func Default() Configer {
	return Use(_DEFAULT_CONFIG)
}

func Use(name string) Configer {
	return Config.Pool[name]
}

func AddConfiger(mode CONFIG_MODE, name, file string) {
	log.Info("add configer", mode, file)
	var configer Configer
	switch mode {
	case INI:
		configer = newIniConfiger(file)
		break
	case JSON:
		// TODO
	case ENV:
		// TODO
		configer = newEnvConfiger()
	default:
		configer = nil
	}
	if configer == nil {
		log.Error("error configer", mode, name, file)
		panic("configer is nil")
	}
	Config.Pool[name] = configer
}

func toInt(value string) int {
	ret, err := strconv.Atoi(value)
	if err != nil {
		log.Error("parse to int error", value, err.Error())
	}
	return ret
}
func toInt64(value string) int64 {
	ret, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Error("parse to int64 error", value, err.Error())
	}
	return ret
}
func toFloat(value string) float64 {
	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Error("parse to float error", value, err.Error())
	}
	return ret
}
func toBool(value string) bool {
	ret := strings.ToLower(value)
	switch ret {
	case "1", "true", "y", "on", "yes":
		return true
	case "0", "false", "n", "off", "no":
		return false
	default:
		log.Error("parse to bool error", value)
	}
	return false
}
