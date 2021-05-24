package biz

import (
	"context"
	"fmt"
	"time"

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

/*===========================================================================*/

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

/*===========================================================================*/

type DumpResult struct {
	From  string
	Reply *dump.Reply
}

type DumpResults struct {
	Results []DumpResult
	Err     error
}

func (d DumpResult) String() string {
	return fmt.Sprintf("{from:%s, Reply:%+v}", d.From, d.Reply)
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

func (m *MasterController) StartPeriodlyDump(interval time.Duration) (resCh <-chan DumpResults, cancel func()) {
	cctx, cancel := context.WithCancel(context.Background())
	ch := make(chan DumpResults, 10)
	go func(ch chan DumpResults, ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				close(ch)
				return
			case <-time.After(interval):
				results, err := m.Dump(ctx)
				ch <- DumpResults{
					Results: results,
					Err:     err,
				}
				continue
			}
		}
	}(ch, cctx)
	return ch, cancel
}

/*===========================================================================*/

func (m *MasterController) WorkerAddrs() []string {
	return m.nm.WorkerAddrs()
}
