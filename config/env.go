package config

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
)

type EnvConfiger struct {
	data map[string]string
	mux  *sync.RWMutex
}

func newEnvConfiger() Configer {
	config := &EnvConfiger{}
	config.data = make(map[string]string)
	config.mux = new(sync.RWMutex)
	config.mux.Lock()
	defer config.mux.Unlock()
	config.load()
	return config
}

// 读.env文件，如果有就加载，没有就直接使用系统的
func (this *EnvConfiger) load() {
	file, err := os.Open(".env")
	defer file.Close()
	if err != nil {
		log.Warn("can not found .env file, use system env directly")
		return
	}
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, " ")
		if len(line) > 0 {
			if line[0] == '#' {
				continue
			}
			if strings.Contains(line, "=") {
				// is a key value
				kvSet := strings.SplitN(line, "=", 2)
				key := strings.Trim(kvSet[0], " ")
				value := strings.Trim(kvSet[1], " ")
				log.Info("set env", key, value)
				this.Set(key, value)
			} else {
				// white space or comment
				// do nothing
			}
		}
		if err != nil || err == io.EOF {
			break
		}
	}
	log.Info("load .env finish")
}

func (this *EnvConfiger) Get(key string) string {
	return os.Getenv(key)
}
func (this *EnvConfiger) Set(key, value string) {
	os.Setenv(key, value)
}
func (this *EnvConfiger) Bool(key string) bool {
	return toBool(this.Get(key))
}
func (this *EnvConfiger) Int(key string) int {
	return toInt(this.Get(key))
}
func (this *EnvConfiger) Int64(key string) int64 {
	return toInt64(this.Get(key))
}
func (this *EnvConfiger) Float(key string) float64 {
	return toFloat(this.Get(key))
}
func (this *EnvConfiger) String(key string) string {
	return this.Get(key)
}
