package util

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadFile(absoluteFileName string) []string {
	file, err := os.Open(absoluteFileName)
	defer file.Close()
	lines := []string{}
	if err != nil {
		log.Error("file is not exist", absoluteFileName)
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
	log.Info("read file finish", absoluteFileName)
	return lines
}
