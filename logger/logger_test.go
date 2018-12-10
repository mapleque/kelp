package logger

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func clear(absoluteFileName string) {
	os.Remove(absoluteFileName)
}

func assertLogFile(filepath string, size int) error {
	body, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	if len(string(body)) != size {
		return fmt.Errorf("%s content is not expect:%d %d", filepath, len(string(body)), size)
	}
	return nil
}

func TestPackageLog(t *testing.T) {
	filepath := "./test_package_log.log"
	clear(filepath)
	RedirectTo(filepath)
	Log("TAG", "some logs here")
	Debug("some logs here")
	Info("some logs here")
	Warn("some logs here")
	Error("some logs here")
	if err := assertLogFile(filepath, 1639); err != nil {
		t.Error(err)
	}
}

func TestLogger(t *testing.T) {
	filepath := "./test_logger.log"
	clear(filepath)
	testLogger := Add("test_logger", filepath)
	otherLogger := Get("test_logger")
	testLogger.Log("TAG", "some logs here")
	testLogger.Debug("some logs here")
	testLogger.Info("some logs here")
	otherLogger.Warn("some logs here")
	otherLogger.Error("some logs here")
	if err := assertLogFile(filepath, 1453); err != nil {
		t.Error(err)
	}
}

func TestOutput(t *testing.T) {
	filepath := "./test_output.log"
	clear(filepath)
	testLogger := Add("test_output", filepath).SetOutput(false).SetTagOutput("TAG", true)
	testLogger.Log("TAG", "some logs here")
	testLogger.Debug("some logs here")
	testLogger.Info("some logs here")
	testLogger.Warn("some logs here")
	testLogger.Error("some logs here")
	if err := assertLogFile(filepath, 41); err != nil {
		t.Error(err)
	}
}

func TestCallstack(t *testing.T) {
	filepath := "./test_callstack.log"
	clear(filepath)
	testLogger := Add("test_callstack", filepath).WithCallstack(false).WithTagCallstack("TAG", true)
	testLogger.Log("TAG", "some logs here")
	testLogger.Debug("some logs here")
	testLogger.Info("some logs here")
	testLogger.Warn("some logs here")
	testLogger.Error("some logs here")
	if err := assertLogFile(filepath, 563); err != nil {
		t.Error(err)
	}
}

func TestRotate(t *testing.T) {
	filepath := "./test_rotate.log"
	testLogger := Add("test_rotate", filepath).SetRotateSize(50).SetRotateFiles(1)
	testLogger.Log("TAG", "some logs here 1")
	testLogger.Log("TAG", "some logs here 2")
	testLogger.Log("TAG", "some logs here 3")
	testLogger.Log("TAG", "some logs here 4")
	if err := assertLogFile(filepath, 43); err != nil {
		t.Error(err)
	}
	if err := assertLogFile(filepath+".0", 86); err != nil {
		t.Error(err)
	}
}
