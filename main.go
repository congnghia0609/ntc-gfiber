/**
 *
 * @author nghiatc
 * @since Jan 5, 2021
 */

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/congnghia0609/ntc-gconf/nconf"
	"github.com/natefinch/lumberjack"
)

// initNConf init file config
func initNConf() {
	_, b, _, _ := runtime.Caller(0)
	wdir := filepath.Dir(b)
	fmt.Println("wdir:", wdir)
	nconf.Init(wdir)
}

// https://github.com/natefinch/lumberjack
func initLogger() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "/data/log/ntc-gfiber/ntc-gfiber.log",
		MaxSize:    10,   // 10 megabytes. Defaults to 100 MB.
		MaxBackups: 3,    // maximum number of old log files to retain.
		MaxAge:     28,   // maximum number of days to retain old log files
		Compress:   true, // disabled by default
	})
}

// increaseLimit increase resources limitations: ulimit -aH
func increaseLimit() {
	var rlimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	rlimit.Cur = rlimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit); err != nil {
		panic(err)
	}
	log.Printf("rlimit.Max = %d\n", rlimit.Max)
	log.Printf("rlimit.Cur = %d\n", rlimit.Cur)
}

func main() {
	////// -------------------- Init System -------------------- //////
	// Increase resources limitations
	increaseLimit()

	// Init NConf
	initNConf()

	//// Init Logger
	if "development" != nconf.GetEnv() {
		log.Printf("============== LogFile: /data/log/ntc-gfiber/ntc-gfiber.log")
		initLogger()
	}

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	////// -------------------- Start WebServer -------------------- //////
	// StartWebServer
	// go server.StartWebServer("webserver")

	// Hang thread Main.
	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C) SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)
	// Block until we receive our signal.
	<-c
	log.Println("################# End Main #################")
}
