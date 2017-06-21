package crontab

import (
	"strconv"
	"strings"
	"time"
)

type Crontab interface {
	Triger(taskId string)
}

type CrontabWrapper struct {
	Name    string
	Crontab Crontab
	Expr    string
}

var crontabs map[string]*CrontabWrapper

func init() {
	crontabs = make(map[string]*CrontabWrapper)
}

func Regist(expr, name string, crontab Crontab) {
	crontabs[name] = &CrontabWrapper{name, crontab, expr}
}

func GetInfo() map[string]*CrontabWrapper {
	return crontabs
}

func Run() {
	log.Info("crontab starting ...")
	done := make(chan bool, 1)
	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			for _, c := range crontabs {
				go func(c *CrontabWrapper) {
					if triger(t, c.Expr) {
						c.Crontab.Triger(c.Name)
						log.Info("[crontab]", "triger", c)
					}
				}(c)
			}
		}
	}()
	<-done
}

// expr index: min, hour, day, month, week
// each expr define
// <expr> := *|*/n|n|m-n|<expr>,...
func triger(t time.Time, expr string) bool {
	arr := strings.Split(expr, " ")
	if len(arr) != 5 {
		return false
	}
	for i, e := range arr {
		if !match(t, e, i) {
			return false
		}
	}
	return true
}

func match(t time.Time, expr string, index int) bool {
	second := t.Second()
	if second != 0 {
		return false
	}
	needle := 0
	switch index {
	case 0:
		needle = t.Minute()
	case 1:
		needle = t.Hour()
	case 2:
		needle = t.Day()
	case 3:
		needle = int(t.Month())
	case 4:
		needle = int(t.Weekday())
	}
	switch {
	case expr == "*":
		return true
	case strings.Contains(expr, ","):
		arr := strings.Split(expr, ",")
		for _, e := range arr {
			if match(t, e, index) {
				return true
			}
		}
		return false
	case strings.Contains(expr, "/"):
		arr := strings.Split(expr, "/")
		if arr[0] != "*" {
			return false
		}
		mod, err := strconv.Atoi(arr[1])
		if err != nil {
			return false
		}
		return needle%mod == 0
	case strings.Contains(expr, "-"):
		arr := strings.Split(expr, "/")
		min, err := strconv.Atoi(arr[0])
		if err != nil {
			return false
		}
		max, err := strconv.Atoi(arr[1])
		if err != nil {
			return false
		}
		return min < needle && needle < max
	}
	value, err := strconv.Atoi(expr)
	if err != nil {
		return false
	}
	return needle == value
}
