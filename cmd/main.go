package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/Water-W/PVP/pkg/biz"
	pvplog "github.com/Water-W/PVP/pkg/log"
)

var (
	log = pvplog.Get("main")

	// master args
	listenPort  = flag.Int("l", 18000, "master listening port")
	httpPort    = flag.Int("p", 18010, "visual html's http port")
	interactive = flag.Bool("i", true, "start a simple cli for master")

	// worker args
	masterIP = flag.String("m", "", "connects to master with ip:port")
	//127.0.0.1:
	datahttp = flag.String("w", "http://39.104.200.8:2404/report", "the network ip:port work connect to")

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
	c := &biz.MasterConfig{
		ListenPort: *listenPort,
	}
	ctrl, err := biz.NewMasterController(c)
	if err != nil {
		log.Errorf("new master controller err:%v", err)
		return
	}

	if !*interactive && *httpPort == -1 {
		log.Error("no cli and http flag is given. exiting.")
		return
	}
	if *interactive {
		go cli(ctrl)
	}
	if *httpPort != -1 {
		go myhttp(ctrl)
	}
	intCh := make(chan os.Signal, 1)
	signal.Notify(intCh, os.Interrupt)
	<-intCh
}

/*===========================================================================*/

func worker() {
	c := &biz.WorkerConfig{
		MasterAddr: *masterIP,
		URL:        *datahttp,
		NodeQuery:  `{ID,Protocols}`,
		LinksQuery: `Peers`,
	}
	wc, err := biz.NewWorkerController(c)
	if err != nil {
		log.Errorf("new worker controller err:%v", err)
		return
	}
	err = wc.ConnectAndServe()
	if err != nil {
		log.Errorf("connect and serve: err:%v", err)
	}
}

/*===========================================================================*/
func setLoggerLevel() {

	pvplog.SetLoggerLevel(*loggerLevel)

}
