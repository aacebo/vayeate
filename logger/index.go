package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const path = "./logs"
const logFlag = log.Ldate | log.Ltime | log.Lshortfile
const fileFlag = os.O_RDWR | os.O_CREATE | os.O_APPEND

type Logger struct {
	info  *log.Logger
	warn  *log.Logger
	debug *log.Logger
	error *log.Logger
}

func New(prefix string) *Logger {
	err := os.Mkdir(path, 0755)

	if err != nil && os.IsExist(err) == false {
		panic(err)
	}

	infof, err := os.OpenFile(fmt.Sprintf("%s/info.log", path), fileFlag, 0666)

	if err != nil {
		panic(err)
	}

	warnf, err := os.OpenFile(fmt.Sprintf("%s/warn.log", path), fileFlag, 0666)

	if err != nil {
		panic(err)
	}

	debugf, err := os.OpenFile(fmt.Sprintf("%s/debug.log", path), fileFlag, 0666)

	if err != nil {
		panic(err)
	}

	errf, err := os.OpenFile(fmt.Sprintf("%s/error.log", path), fileFlag, 0666)

	if err != nil {
		panic(err)
	}

	return &Logger{
		info:  log.New(io.MultiWriter(os.Stdout, infof), fmt.Sprintf("%s %s ", prefix, "info"), logFlag),
		warn:  log.New(io.MultiWriter(os.Stdout, warnf), fmt.Sprintf("%s %s ", prefix, "warn"), logFlag),
		debug: log.New(io.MultiWriter(os.Stdout, debugf), fmt.Sprintf("%s %s ", prefix, "debug"), logFlag),
		error: log.New(io.MultiWriter(os.Stderr, errf), fmt.Sprintf("%s %s ", prefix, "error"), logFlag),
	}
}

func (self *Logger) Info(v ...any) {
	self.info.Print(v...)
}

func (self *Logger) Infof(format string, v ...any) {
	self.info.Printf(format, v...)
}

func (self *Logger) Infoln(v ...any) {
	self.info.Println(v...)
}

func (self *Logger) Warn(v ...any) {
	self.warn.Print(v...)
}

func (self *Logger) Warnf(format string, v ...any) {
	self.warn.Printf(format, v...)
}

func (self *Logger) Warnln(v ...any) {
	self.warn.Println(v...)
}

func (self *Logger) Debug(v ...any) {
	self.debug.Print(v...)
}

func (self *Logger) Debugf(format string, v ...any) {
	self.debug.Printf(format, v...)
}

func (self *Logger) Debugln(v ...any) {
	self.debug.Println(v...)
}

func (self *Logger) Error(v ...any) {
	self.error.Fatal(v...)
}

func (self *Logger) Errorf(format string, v ...any) {
	self.error.Fatalf(format, v...)
}

func (self *Logger) Errorln(v ...any) {
	self.error.Fatalln(v...)
}
