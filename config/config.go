package config

import "flag"

type _CONFIG_MODE int32

const (
	INI _CONFIG_MODE = iota
	XML
	JSON
)

const _DEFAULT_CONFIG = "default_config"

type _ConfigPool struct {
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

var config *_ConfigPool

func init() {
	if config != nil {
		return
	}
	config = &_ConfigPool{}
	config.Pool = make(map[string]Configer)
}

// This method provide a convinience way to use config package when you just need configer value all in one.
// Extra, it can read configer file path from arguments in cmd.
// Usage:
//		configer := config.InitDefault(config.INI, "ini", "./config.ini", "config file path with param --ini")
//		value := configer.Get("section.key")
//		// use the value as string
//
//		// when in other scope
//		configer := config.Default()
//		// use this configer as you wish
func InitDefault(mode _CONFIG_MODE, operationParam, defaultValue, tips string) Configer {
	confFile := flag.String(operationParam, defaultValue, tips)
	flag.Parse()
	if *confFile == "" {
		panic("run with -h to find usage")
	}
	Add(INI, _DEFAULT_CONFIG, *confFile)
	return Default()
}

// This method get default configer
func Default() Configer {
	return Use(_DEFAULT_CONFIG)
}

func Use(name string) Configer {
	return config.Pool[name]
}

// This method is an normal way to user config package.
// You can add several configer and use them seperately.
// Usage:
//		config.Add(config.INI, "configer_name_1", "./config.ini")
//		config.Add(config.INI, "configer_name_2", "./config.ini")
//		// ...
//		// in other context
//		configer1 := config.Use("configer_name_1")
//		configer2 := config.Use("configer_name_2")
func Add(mode _CONFIG_MODE, name, file string) {
	log.Info("add configer", mode, file)
	var configer Configer
	switch mode {
	case INI:
		configer = NewIniConfiger(file)
		break
	case XML:
		// TODO
	case JSON:
		// TODO
	default:
		configer = nil
	}
	if configer == nil {
		log.Error("error configer", mode, name, file)
		panic("configer is nil")
	}
	config.Pool[name] = configer
}
