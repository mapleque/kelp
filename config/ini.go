package config

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

type IniConfiger struct {
	file string
	data map[string]IniGroup
	mux  *sync.RWMutex
}

type IniGroup map[string]string

var _DEFAULT_GROUP = "default"

func newIniConfiger(file string) Configer {
	config := &IniConfiger{}
	config.file = file
	config.data = make(map[string]IniGroup)
	config.mux = new(sync.RWMutex)
	config.mux.Lock()
	defer config.mux.Unlock()
	config.load()
	return config
}

func (config *IniConfiger) load() {
	// 读ini文件
	file, err := os.Open(config.file)
	defer file.Close()
	if err != nil {
		panic("config file is not exist " + config.file)
	}
	buf := bufio.NewReader(file)
	currentGroup := _DEFAULT_GROUP
	for {
		line, err := buf.ReadString('\n')
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, " ")
		if strings.Contains(line, ";") {
			line = strings.Split(line, ";")[0]
		}
		if len(line) > 0 {
			start := strings.Index(line, "[")
			end := strings.LastIndex(line, "]")
			if start >= 0 && end > start {
				// is a group
				currentGroup = string([]rune(line)[start+1 : end-start])
			} else if strings.Contains(line, "=") {
				// is a key value
				kvSet := strings.SplitN(line, "=", 2)
				key := strings.Trim(kvSet[0], " ")
				value := strings.Trim(kvSet[1], " ")
				config.set(currentGroup+"."+key, value)
			} else {
				// white space or comment
				// do nothing
			}
		}
		if err != nil || err == io.EOF {
			break
		}
	}
	log.Info("load config file finish", config.file)
}

// key accept group.element
func (config *IniConfiger) Get(key string) string {
	if len(key) < 1 {
		log.Error("config key is empty")
		return ""
	}
	config.mux.RLock()
	defer config.mux.RUnlock()
	group, element := escapeKey(key)
	if g, ok := config.data[group]; ok {
		if v, ok := g[element]; ok {
			return v
		}
	}
	log.Error("config property not exist", key)
	return ""
}

// key accept group.element
func (config *IniConfiger) Set(key, value string) {
	if len(key) < 1 {
		log.Error("config key is empty")
		return
	}
	config.mux.Lock()
	defer config.mux.Unlock()
	config.set(key, value)
}

func (config *IniConfiger) set(key, value string) {
	group, element := escapeKey(key)
	if _, ok := config.data[group]; !ok {
		config.data[group] = make(map[string]string)
	}
	config.data[group][element] = value
	log.Info("set config", group+","+key+"="+value)
}

func (config *IniConfiger) Bool(key string) bool {
	return toBool(config.Get(key))
}
func (config *IniConfiger) Int(key string) int {
	return toInt(config.Get(key))
}
func (config *IniConfiger) Int64(key string) int64 {
	return toInt64(config.Get(key))
}
func (config *IniConfiger) Float(key string) float64 {
	return toFloat(config.Get(key))
}
func (config *IniConfiger) String(key string) string {
	return config.Get(key)
}

func escapeKey(key string) (string, string) {
	groupKey := strings.Split(strings.ToLower(key), ".")
	var group, element string
	if len(groupKey) > 1 {
		group = groupKey[0]
		element = groupKey[1]
	} else {
		group = _DEFAULT_GROUP
		element = groupKey[0]
	}
	return group, element
}
