package echo

import "encoding/gob"

var ServiceName = "echo"

// Service .
type Service struct{}

type Args struct {
	Data string
}

type Reply struct {
	Data string
}

func init() {
	gob.Register(Args{})
	gob.Register(&Reply{})
}

func (e *Service) Echo(args Args, reply *Reply) error {
	reply.Data = args.Data
	return nil
}
