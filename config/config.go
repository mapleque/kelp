package config

type CONFIG_MODE int32

const (
	INI = iota
	XML
	JSON
)

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

func Use(name string) Configer {
	return Config.Pool[name]
}

func AddConfiger(mode CONFIG_MODE, name, file string) {
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
	Config.Pool[name] = configer
}
