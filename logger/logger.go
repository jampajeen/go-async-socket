package logger

import (
	"log"
	"os"
	"runtime"
	"strings"
)

var debug bool

func init() {
	debug = os.Getenv("DEBUG") != ""
}

func Fatal(err error) {
	pc, fn, line, _ := runtime.Caller(1)
	if debug {
		log.Fatalf("[FATAL] %s [%s:%s:%d]", err, runtime.FuncForPC(pc).Name(), fn, line)
	} else {
		log.Fatalf("[FATAL] %s [%s:%d]", err, fn, line)
	}
}

func Error(err error) {
	pc, fn, line, _ := runtime.Caller(1)
	if debug {
		log.Printf("[ERROR] %s [%s:%s:%d]", err, runtime.FuncForPC(pc).Name(), fn, line)
	} else {
		log.Printf("[ERROR] %s [%s:%d]", err, fn, line)
	}
}

func Warn(msg string, vars ...interface{}) {
	log.Printf(strings.Join([]string{"[WARN ]", msg}, " "), vars...)
}

func Info(msg string, vars ...interface{}) {
	log.Printf(strings.Join([]string{"[INFO ]", msg}, " "), vars...)
}

func Debug(msg string, vars ...interface{}) {
	if debug {
		log.Printf(strings.Join([]string{"[DEBUG]", msg}, " "), vars...)
	}
}
