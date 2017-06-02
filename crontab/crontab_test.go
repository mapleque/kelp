package crontab

import (
	"testing"
	"time"
)

type TestTrigerCase struct {
	t      time.Time
	expr   string
	result bool
}

func TestCrontabTriger(t *testing.T) {
	cases := []TestTrigerCase{
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 58, 0, 0, time.UTC),
			"* * * * *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 58, 0, 0, time.UTC),
			"*/1 * * * *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 0, 0, 0, time.UTC),
			"*/5 * * * *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 58, 1, 0, time.UTC),
			"*/5 * * * *",
			false},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 0, 0, 0, time.UTC),
			"0 */1 * * *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 58, 0, 0, time.UTC),
			"0 */1 * * *",
			false},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 58, 1, 0, time.UTC),
			"0 */1 * * *",
			false},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 0, 0, 0, time.UTC),
			"0 */5 * * *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 58, 0, 0, time.UTC),
			"0 */5 * * *",
			false},
		TestTrigerCase{
			time.Date(2017, 10, 26, 0, 0, 0, 0, time.UTC),
			"0 0 */1 * *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 58, 0, 0, time.UTC),
			"0 0 */1 * *",
			false},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 0, 0, 0, time.UTC),
			"0 10 26 */1 *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 0, 0, 0, time.UTC),
			"0 10 26 */2 *",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 0, 0, 0, time.UTC),
			"0 10 26 */1 4",
			true},
		TestTrigerCase{
			time.Date(2017, 10, 26, 10, 0, 0, 0, time.UTC),
			"0 10 26 */1 5",
			false},
	}

	for _, c := range cases {
		if triger(c.t, c.expr) != c.result {
			if !c.result {
				t.Error("wrong triger", c.t, c.expr)
			} else {
				t.Error("not triger", c.t, c.expr)
			}
		}
	}
}
