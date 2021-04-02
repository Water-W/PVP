package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/Water-W/PVP/pkg/biz"
	pvplog "github.com/Water-W/PVP/pkg/log"
	"github.com/Water-W/PVP/pkg/net"
	"github.com/Water-W/PVP/pkg/rpc/echo"
)

var (
	log = pvplog.Get("main")

	// master args
	listenPort  = flag.Int("l", 8000, "master listening port")
	httpPort    = flag.Int("p", -1, "http port")
	interactive = flag.Bool("i", true, "start a simple cli for master")

	// worker args
	masterIP = flag.String("m", "", "connects to master with ip:port")

	// log level
	loggerLevel = flag.String("L", "info", "logger level")
)

func main() {
	flag.Parse()

	setLoggerLevel()

	if *masterIP == "" {
		master()
	} else {
		worker()
	}
}

func master() {
	m, err := net.NewMaster()
	if err != nil {
		log.Error(err)
		return
	}
	ctrl := biz.NewMasterController(m)
	err = m.Listen(*listenPort)
	if err != nil {
		log.Error(err)
		return
	}
	if !*interactive || *httpPort == 1 {
		log.Error("no cli and http flag is given. exiting.")
		return
	}
	if *interactive {
		go cli(ctrl)
	}
	if *httpPort != -1 {
		go http(ctrl)
	}
	intCh := make(chan os.Signal, 1)
	signal.Notify(intCh, os.Interrupt)
	<-intCh
}

/*===========================================================================*/

func worker() {
	w, err := net.NewWorker()
	if err != nil {
		log.Error(err)
		return
	}
	registerService(w)
	err = w.Connect(*masterIP)
	if err != nil {
		log.Error(err)
		return
	}
	// block forever
	select {}
}

func registerService(w *net.Worker) {
	w.RegisterName(echo.ServiceName, &echo.Service{})
	// TODO
}

/*===========================================================================*/
func setLoggerLevel() {

	pvplog.SetLoggerLevel(*loggerLevel)

}
