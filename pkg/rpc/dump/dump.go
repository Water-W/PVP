package dump

import (
	"encoding/gob"

	"github.com/Water-W/PVP/pkg/log"
)

var logger = log.Get("dump")

func init() {
	gob.Register(map[string]interface{}{})
}

type LinkMeasurer interface {
	GetLinks() (map[string]interface{}, error)
}
type NodeMeasurer interface {
	// GetNode returns a gob-serializable interface{}
	GetNode() (interface{}, error)
}

type Measurer interface {
	LinkMeasurer
	NodeMeasurer
}

var ServiceName = "dump"

// Service .
type Service struct {
	nm NodeMeasurer
	lm LinkMeasurer
}

func (s *Service) RegisterNodeMeasurer(nm NodeMeasurer) {
	s.nm = nm
}

func (s *Service) RegisterLinkMeasurer(lm LinkMeasurer) {
	s.lm = lm
}

func (s *Service) RegisterMeasurer(m Measurer) {
	s.nm = m
	s.lm = m
}

type Args struct{}
type Reply struct {
	Node       interface{}
	Links      map[string]interface{}
	ErrMessage string
}

func (s *Service) Dump(args Args, reply *Reply) (err error) {
	logger.Debugf("dump request incoming")
	reply.Node, err = s.nm.GetNode()
	if err != nil {
		logger.Infof("err=%+v", err)
		reply.ErrMessage = err.Error()
		return nil
	}
	reply.Links, err = s.lm.GetLinks()
	if err != nil {
		logger.Infof("err=%+v", err)
		reply.ErrMessage = err.Error()
		return nil
	}
	return nil
}
