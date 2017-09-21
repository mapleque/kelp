package config

import (
	"testing"
)

func TestRun(t *testing.T) {
	configer := InitDefault(INI, "ini", "./config.ini", "your config file")
	if configer.Get("test.KEY") != "value" {
		t.Fatal("run failed")
	}
}
