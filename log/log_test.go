package log

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
)

func clear(absoluteFileName string) {
	os.Remove(absoluteFileName)
}

func readFile(absoluteFileName string) []string {
	file, err := os.Open(absoluteFileName)
	defer file.Close()
	lines := []string{}
	if err != nil {
		return lines
	}
	buf := bufio.NewReader(file)
	for {
		line, err := buf.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		line = strings.Trim(line, "\n")
		line = strings.Trim(line, " ")
		lines = append(lines, line)
	}
	return lines
}

func TestLogText(t *testing.T) {
	name := "test_log.log"
	clear("./" + name)
	Log.Pool = make(map[string]*Logger)
	AddLogger(
		name,
		".",
		1,
		1000000000,
		5,
		0,
	)
	Info("here start the log")
	Debug("debug info")
	Warn("warn info")
	Error("error info")
	lines := readFile("./" + name)
	if len(lines) != 28 {
		t.Error("wrong log lines should 28 but ", len(lines))
		return
	}
	if lines[0][21:25] != "INFO" {
		t.Error("wrong log info should INFO but ", lines[0][21:25])
	}
	if lines[1][21:26] != "DEBUG" {
		t.Error("wrong log info should DEBUG but ", lines[1][21:26])
	}
	if lines[10][21:25] != "WARN" {
		t.Error("wrong log info should WARN but ", lines[10][21:25])
	}
	if lines[19][21:26] != "ERROR" {
		t.Error("wrong log info should ERROR but ", lines[19][21:26])
	}
}

func TestLogger(t *testing.T) {
	name := "test_logger.log"
	clear("./" + name)
	Log.Pool = make(map[string]*Logger)
	AddLogger(
		name,
		".",
		1,
		1000000000,
		5,
		0,
	)
	Log.Get(name).Info("here start the log")
	Log.Get(name).Debug("debug info")
	Log.Get(name).Warn("warn info")
	Log.Get(name).Error("error info")
	lines := readFile("./" + name)
	if len(lines) != 28 {
		t.Error("wrong log lines should 28 but ", len(lines))
		return
	}
	if lines[0][21:25] != "INFO" {
		t.Error("wrong log info should INFO but ", lines[0][21:25])
	}
	if lines[1][21:26] != "DEBUG" {
		t.Error("wrong log info should DEBUG but ", lines[1][21:26])
	}
	if lines[10][21:25] != "WARN" {
		t.Error("wrong log info should WARN but ", lines[10][21:25])
	}
	if lines[19][21:26] != "ERROR" {
		t.Error("wrong log info should ERROR but ", lines[19][21:26])
	}
}

func TestLogLevel(t *testing.T) {
	Log.Pool = make(map[string]*Logger)
	names := []string{
		"test_debug.log",
		"test_info.log",
		"test_warn.log",
		"test_error.log",
	}
	for i, name := range names {
		clear("./" + name)
		AddLogger(
			name,
			".",
			1,
			1000000000,
			i+1,
			i+1,
		)
	}
	Info("here start the log")
	Debug("debug info")
	Warn("warn info")
	Error("error info")
	tags := []string{
		"DEBU",
		"INFO",
		"WARN",
		"ERRO",
	}
	line_num := []int{9, 3, 9, 9}
	for i, name := range names {
		lines := readFile("./" + name)
		if len(lines) != line_num[i] {
			t.Error("wrong log lines", name, "should", line_num, "but", len(lines))
			return
		}
		if lines[0][21:25] != tags[i] {
			t.Error("wrong log info should", tags[i], "but", lines[0][21:25])
		}
	}
}
