package biz

import (
	"context"

	"github.com/Water-W/PVP/pkg/net"
	"github.com/Water-W/PVP/pkg/rpc/echo"
)

// biz.MasterController provides serveral control methods
// for master user to call rpc and view network.
type MasterController struct {
	nm *net.Master
}

func NewMasterController(nm *net.Master) *MasterController {
	return &MasterController{
		nm: nm,
	}
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
