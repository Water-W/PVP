package biz

import (
	"github.com/Water-W/PVP/pkg/metrics/json"
	"github.com/Water-W/PVP/pkg/net"
	"github.com/Water-W/PVP/pkg/rpc/dump"
	"github.com/Water-W/PVP/pkg/rpc/echo"
	"github.com/Water-W/PVP/pkg/source/http"
)

type WorkerConfig struct {
	MasterAddr string
	URL        string // remove it from config
	NodeQuery  string
	LinksQuery string
}

type WorkerController struct {
	c  *WorkerConfig
	w  *net.Worker
	ds *dump.Service
}

func NewWorkerController(c *WorkerConfig) (*WorkerController, error) {
	// code here looks silly
	worker, err := net.NewWorker()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	hs := http.NewSource(c.URL)
	jm := &json.Measurer{}
	jm.SetSource(hs)
	jm.SetNodeQuery(c.NodeQuery)

	ds := &dump.Service{}
	ds.RegisterMeasurer(jm)

	err = worker.RegisterName(dump.ServiceName, ds)
	if err != nil {
		return nil, err
	}

	es := &echo.Service{}
	err = worker.RegisterName(echo.ServiceName, es)
	if err != nil {
		return nil, err
	}

	return &WorkerController{
		c:  c,
		w:  worker,
		ds: ds,
	}, nil
}

func (wc *WorkerController) ConnectAndServe() error {
	return wc.w.ConnectAndServe(wc.c.MasterAddr)
}
