package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime/trace"
	"strconv"
	"syscall"
	"time"

	"github.com/sound-of-destiny/qlsc_zhxf/pkg/log"
	"github.com/sound-of-destiny/qlsc_zhxf/pkg/setting"
)

var version = "0.0.1"
var commit = "NA"
var buildstamp string

var configFile = flag.String("config", "", "path to config file")
var homePath = flag.String("homepath", "", "path to grafana install/home path, defaults to working directory")
var pidFile = flag.String("pidfile", "", "path to pid file")

func main() {
	buildstampInt64, _ := strconv.ParseInt(buildstamp, 10, 64)
	if buildstampInt64 == 0 {
		buildstampInt64 = time.Now().Unix()
	}

	setting.BuildVersion = version
	setting.BuildCommit = commit
	setting.BuildStamp = buildstampInt64

	server := NewServer()

	go listenToSystemSignals(server)

	err := server.Run()

	code := server.Exit(err)

	trace.Stop()
	log.Close()

	os.Exit(code)
}

func listenToSystemSignals(server *NewServerImpl) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case _ = <-sighupChan:
			log.Reload()
		case sig := <-signalChan:
			server.Shutdown(fmt.Sprintf("System signal: %s", sig))
		}
	}
}
