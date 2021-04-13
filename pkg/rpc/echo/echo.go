package echo

import "github.com/Water-W/PVP/pkg/log"

var ServiceName = "echo"

// Service .
type Service struct{}

type Args struct {
	Data string
}

type Reply struct {
	Data string
}

func (e *Service) Echo(args Args, reply *Reply) error {
	log.Get("echo").Infof("echo request incoming")
	reply.Data = args.Data
	return nil
}
