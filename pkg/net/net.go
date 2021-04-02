package net

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"reflect"
	"sync"
	"time"

	pvplog "github.com/Water-W/PVP/pkg/log"
)

var (
	log = pvplog.Get("net")
)

var Timeout = 10 * time.Second

type Master struct {
	ln      net.Listener
	rw      sync.RWMutex
	workers map[string]*workerItem
}

func NewMaster() (*Master, error) {
	return &Master{
		rw:      sync.RWMutex{},
		workers: make(map[string]*workerItem),
	}, nil
}

func (m *Master) Listen(port int) error {
	tcpaddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		return err
	}
	go m.serve(ln)
	return nil
}

func (m *Master) Close() error {
	if m.ln != nil {
		return m.ln.Close()
	}
	return nil
}

func (m *Master) serve(l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
			return
		}
		m.handleIncomingConn(conn)
	}
}

func (m *Master) handleIncomingConn(c net.Conn) {
	log.Debugf("new conn incoming: %s", c.RemoteAddr().String())
	cli := rpc.NewClient(c)
	item := &workerItem{
		addr:   c.RemoteAddr().String(),
		conn:   c,
		client: cli,
	}
	m.rw.Lock()
	defer m.rw.Unlock()
	m.workers[item.addr] = item
}

// func (m *Master) handleDelAddr(addr string) {
// 	m.rw.Lock()
// 	defer m.rw.Unlock()
// 	delete(m.workers, addr)
// }

func (m *Master) all() []*workerItem {
	m.rw.RLock()
	defer m.rw.RUnlock()
	out := make([]*workerItem, 0, len(m.workers))
	for _, item := range m.workers {
		out = append(out, item)
	}
	return out
}

func (m *Master) ForAll(methodName string, args interface{}, reply interface{}) (<-chan RpcCall, error) {
	tctx, _ := context.WithTimeout(context.Background(), Timeout)
	return m.forAll(tctx, methodName, args, reply)
}

func (m *Master) forAll(
	ctx context.Context,
	methodName string,
	args interface{},
	reply interface{},
) (<-chan RpcCall, error) {
	items := m.all()
	out := make(chan RpcCall, len(items))
	// create an reply for each call
	dups := makeZeroValueDuplicates(reply, len(items))
	cctx, cc := context.WithCancel(ctx)
	var wg sync.WaitGroup
	wg.Add(len(items))
	for i, item := range items {
		go func(ctx context.Context, item *workerItem, dupReply interface{}) {
			defer wg.Done() //
			// call it and wait until it is done
			var call *rpc.Call
			select {
			case call = <-item.client.Go(methodName, args, dupReply, nil).Done:
				log.Debugf("%s:rpc done", item.addr)
			case <-ctx.Done():
				return
			}
			m.checkRpcCall(call)
			out <- RpcCall{
				Addr: item.addr,
				Call: call,
			}
		}(cctx, item, dups[i])
	}
	// wait for all rpc are done, and then close channel
	go func(wg *sync.WaitGroup) {
		wg.Wait()
		cc()
		close(out)
	}(&wg)
	return out, nil
}

func (m *Master) ForAllSync(
	ctx context.Context,
	methodName string,
	args interface{},
	reply interface{},
) ([]RpcCall, error) {
	callCh, err := m.forAll(ctx, methodName, args, reply)
	if err != nil {
		return []RpcCall{}, err
	}
	out := make([]RpcCall, 0, len(callCh))
	for call := range callCh {
		out = append(out, call)
	}
	return out, nil
}

func (m *Master) WorkerAddrs() []string {
	ws := m.all()
	out := make([]string, len(ws))
	for i := range ws {
		out[i] = ws[i].addr
	}
	return out
}

// checkRpcCall checks rpc and update workerItems.
// It handles the situation that some workers are unavailable.
// It can be safiticated if you consider tons of situations.
func (m *Master) checkRpcCall(c *rpc.Call) {
	// TODO
}

type workerItem struct {
	addr   string
	conn   net.Conn
	client *rpc.Client
}

type RpcCall struct {
	Addr string
	*rpc.Call
}

func makeZeroValueDuplicates(p interface{}, n int) []interface{} {
	out := make([]interface{}, n)
	t := reflect.TypeOf(p).Elem()
	for i := range out {
		out[i] = reflect.New(t).Interface()
	}
	return out
}

/*===========================================================================*/

type Worker struct {
	s *rpc.Server
}

func NewWorker() (*Worker, error) {
	return &Worker{
		s: rpc.NewServer(),
	}, nil
}

func (w *Worker) Connect(masterAddr string) error {
	conn, err := net.Dial("tcp", masterAddr)
	if err != nil {
		return err
	}
	log.Infof("connected to %s", masterAddr)
	w.s.ServeConn(conn)
	return nil
}

func (w *Worker) Register(service interface{}) error {
	return w.s.Register(service)
}

func (w *Worker) RegisterName(name string, service interface{}) error {
	return w.s.RegisterName(name, service)
}
