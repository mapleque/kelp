package config

import (
	"testing"
)

func TestEnv(t *testing.T) {
	AddConfiger(ENV, "test_env", "")
	conf := Use("test_env")
	conf.Set("TEST_ENV", "test_env")
	if conf.Get("TEST_ENV") != "test_env" {
		t.Error("test env failed")
	}
}
