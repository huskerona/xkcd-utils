package main

import (
	"flag"
	"fmt"
	xkcdlib "github.com/huskerona/xkcd/xkcd-lib"
	"log"
	"os"
	"os/user"
	"time"
)

var verbose = flag.Bool("v", false, "verbose printout of the operations.")
var showOIP = flag.Bool("o", false, "show offline index path and quit")
var writeLog = flag.Bool("l", false, "write to log file (xkcd.log)")

func main() {
	flag.Parse()
	xkcdlib.SetVerboseFlag(*verbose)

	setOfflineFolderForUser()

	if err := xkcdlib.IsOfflineIndexAvailable(xkcdlib.OfflineIndexPath); err != nil {
		log.Fatal(err)
	}

	logFile, err := setLogFile()

	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()
	xkcdlib.SetLogFile(logFile)

	if *showOIP {
		fmt.Printf("Index folder at: %s\n", xkcdlib.OfflineIndexPath)
		os.Exit(0)
	}

	start := time.Now()

	if err := xkcdlib.SyncAll(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nDuration: %s\n", time.Since(start))

	os.Exit(0)
}

func setOfflineFolderForUser() {
	currentUser, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	xkcdlib.SetOfflineIndexPath(currentUser.HomeDir)
}

func setLogFile() (logFile *os.File, err error) {

	if *writeLog {
		logPath := fmt.Sprintf("%s%cxkcd.log", xkcdlib.OfflineIndexPath, os.PathSeparator)

		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			if logFile, err = os.Create(logPath); err != nil {
				return nil, err
			}
		} else {
			if logFile, err = os.Open(logPath); err != nil {
				return nil, err
			}
		}
	} else {
		if logFile, err = os.Create(os.DevNull); err != nil {
			log.Fatal(err)
		}
	}

	return
}