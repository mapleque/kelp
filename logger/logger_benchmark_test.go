package logger

import (
	"testing"
)

func BenchmarkLogger(b *testing.B) {
	filepath := "./benchmark_logger.log"
	clear(filepath)
	testLogger := Add("benchmark_logger", filepath).SetRotateSize(100).SetRotateFiles(2)
	for i := 0; i < b.N; i++ {
		testLogger.Log("TAG", "some logs here", i)
	}
}
