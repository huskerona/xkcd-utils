package xkcd_lib

import (
	"fmt"
	"log"
	"os"
	"time"
)

var isVerbose bool
var logFile *os.File

func SetVerboseFlag(verbose bool) {
	isVerbose = verbose
}

func SetLogFile(file *os.File) {
	logFile = file
}

func ShowVerbose(message string) {

	if isVerbose {
		fmt.Fprint(os.Stdout, message)
	}
}

func trace(msg string) func() {
	start := time.Now()
	if _, err := logFile.WriteString(fmt.Sprintf("start %s\n", msg)); err != nil {
		log.Fatal(err)
	}

	return func() {
		logFile.WriteString(fmt.Sprintf("exit %s (%s)\n", msg, time.Since(start)))
	}
}
