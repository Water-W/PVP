package net

import (
	"context"
	"fmt"
	"net"
	"net/rpc"
	"reflect"
	"sync"

	pvplog "github.com/Water-W/PVP/pkg/log"
)

var log = pvplog.Get("net")

type Master struct {
	rw      sync.RWMutex
	workers map[string]*workerItem
	cancel  func()
}

func NewMaster() (*Master, error) {
	return &Master{
		rw:      sync.RWMutex{},
		workers: make(map[string]*workerItem),
	}, nil
}

func (m *Master) listen(port int) error {
	tcpaddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", tcpaddr)
	if err != nil {
		return err
	}
	cctx, cc := context.WithCancel(context.Background())
	m.cancel = cc
	go m.serve(cctx, ln)
	return nil
}

func (m *Master) serve(ctx context.Context, l net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Error(err)
			select {
			case <-ctx.Done():
				return
			default:
			}
			continue
		}
		m.handleIncomingConn(conn)
	}
}

func (m *Master) handleIncomingConn(c net.Conn) {
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

func (m *Master) handleDelAddr(addr string) {
	m.rw.Lock()
	defer m.rw.Unlock()
	delete(m.workers, addr)
}

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
	items := m.all()
	out := make(chan RpcCall, len(items))
	// create an reply for each call
	dups := makeZeroValueDuplicates(reply, len(items))
	for i, item := range items {
		go func(item *workerItem, dupReply interface{}) {
			// call it and wait until it is done
			call := <-item.client.Go(methodName, args, dupReply, nil).Done
			out <- RpcCall{
				Addr: item.addr,
				Call: call,
			}
		}(item, dups[i])
	}
	return out, nil
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
	w.s.ServeConn(conn)
	return nil
}

func (w *Worker) Register(service interface{}) error {
	return w.s.Register(service)
}
