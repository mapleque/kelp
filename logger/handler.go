package logger

import (
	"log"
	"os"
)

type handler struct {
	fileHandler *os.File
	logger      *log.Logger
}

func (this *handler) Open(filepath string) {
	if filepath == "" {
		this.logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	} else {
		if fileHandler, err := os.OpenFile(
			filepath,
			os.O_RDWR|os.O_APPEND|os.O_CREATE,
			0666,
		); err != nil {
			panic(err)
		} else {
			this.fileHandler = fileHandler
			this.logger = log.New(fileHandler, "", log.Ldate|log.Ltime)
		}
	}
}

func (this *handler) Close() {
	if this.fileHandler != nil {
		this.fileHandler.Close()
		this.fileHandler = nil
	}
	this.logger = nil
}

func (this *handler) Println(message ...interface{}) {
	if this.logger == nil {
		panic("This is a bug, logger is nil, please create an issue to tell us.")
	} else {
		this.logger.Println(message...)
	}
}
