package dump

type LinkMeasurer interface {
	GetLinks() (map[string]interface{}, error)
}
type NodeMeasurer interface {
	// GetNode returns a gob-serializable interface{}
	GetNode() (interface{}, error)
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

type NodeArgs struct{}
type NodeReply struct {
	NodeData interface{}
}

func (s *Service) DumpNodes(args NodeArgs, reply *NodeReply) (err error) {
	reply.NodeData, err = s.nm.GetNode()
	if err != nil {
		return err
	}
	return nil
}

type LinkArgs struct{}
type LinkReply struct {
	LinkData map[string]interface{}
}

func (s *Service) DumpLinks(args LinkArgs, reply *LinkReply) (err error) {
	reply.LinkData, err = s.lm.GetLinks()
	if err != nil {
		return err
	}
	return nil
}
