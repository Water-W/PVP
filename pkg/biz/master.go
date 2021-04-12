package biz

import (
	"context"

	"github.com/Water-W/PVP/pkg/log"
	"github.com/Water-W/PVP/pkg/net"
	"github.com/Water-W/PVP/pkg/rpc/dump"
	"github.com/Water-W/PVP/pkg/rpc/echo"
)

var logger = log.Get("ctrl")

type MasterConfig struct {
	ListenPort int
}

// biz.MasterController provides serveral control methods
// for master user to call rpc and view network.
type MasterController struct {
	nm *net.Master
}

func NewMasterController(c *MasterConfig) (*MasterController, error) {
	nm, err := net.NewMaster()
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	err = nm.Listen(c.ListenPort)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &MasterController{
		nm: nm,
	}, nil
}

type EchoResult struct {
	From   string
	Result string
}

func (m *MasterController) Echo(ctx context.Context, s string) ([]EchoResult, error) {
	calls, err := m.nm.ForAllSync(
		ctx,
		echo.ServiceName+".Echo",
		echo.Args{
			Data: s,
		},
		&echo.Reply{},
	)
	out := make([]EchoResult, len(calls))
	for i := range calls {
		out[i] = EchoResult{
			From:   calls[i].Addr,
			Result: calls[i].Reply.(*echo.Reply).Data,
		}
	}
	return out, err
}

func (m *MasterController) WorkerAddrs() []string {
	return m.nm.WorkerAddrs()
}

type DumpResult struct {
	From  string
	Reply *dump.Reply
}

func (m *MasterController) Dump(ctx context.Context) ([]DumpResult, error) {
	calls, err := m.nm.ForAllSync(
		ctx,
		dump.ServiceName+".Dump",
		dump.Args{},
		&dump.Reply{},
	)
	out := make([]DumpResult, len(calls))
	for i := range calls {
		out[i] = DumpResult{
			From:  calls[i].Addr,
			Reply: calls[i].Reply.(*dump.Reply),
		}
	}
	return out, err
}
