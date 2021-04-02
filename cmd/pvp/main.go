package main

import (
	"flag"

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
)

func main() {
	flag.Parse()
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
	if *interactive {
		cli(ctrl)
		return
	}
	if *httpPort != -1 {
		http(ctrl)
		return
	}
	log.Error("no cli and http flag is given. exiting.")
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
