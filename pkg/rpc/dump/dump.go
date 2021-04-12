package dump

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
	Node  interface{}
	Links map[string]interface{}
}

func (s *Service) Dump(args Args, reply *Reply) (err error) {
	reply.Node, err = s.nm.GetNode()
	if err != nil {
		return err
	}
	reply.Links, err = s.lm.GetLinks()
	if err != nil {
		return err
	}
	return nil
}
