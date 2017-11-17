package config

import (
	"testing"
)

func TestIni(t *testing.T) {
	AddConfiger(INI, "test_ini", "./config.ini")
	conf := Use("test_ini")
	if conf.Get("sec1.V1") != "hello" {
		t.Error("sec1.V1 wrong")
	}
	if conf.Int("sec1.V2") != 1 {
		t.Error("sec1.V2 wrong")
	}
	if conf.Int64("sec1.V2") != 1 {
		t.Error("sec1.V2 wrong")
	}
	if conf.Bool("sec1.V3") != true {
		t.Error("sec1.V3 wrong")
	}
	if conf.String("sec1.V4") != "hello" {
		t.Error("sec1.V4 wrong")
	}
	if conf.Float("sec1.V5") != 3.14 {
		t.Error("sec1.V5 wrong")
	}
	if conf.Get("sec2.V1") != "1:a_d@**." {
		t.Error("sec2.V1 wrong")
	}
}
