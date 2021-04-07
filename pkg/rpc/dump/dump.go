package dump

type LinkMetrics interface {
	Links() []string
	// LinkMetric returns a gob-serializable interface{}
	LinkMetric(s string) (interface{}, error)
}
type NodeMetrics interface {
	// NodeMetric returns a gob-serializable interface{}
	NodeMetric() (interface{}, error)
}

var ServiceName = "dump"

// Service .
type Service struct {
	nm NodeMetrics
	lm LinkMetrics
}

func (s *Service) RegisterNodeMetrics(nm NodeMetrics) {
	s.nm = nm
}

func (s *Service) RegisterLinkMetrics(lm LinkMetrics) {
	s.lm = lm
}

type NodeArgs struct{}
type NodeReply struct {
	NodeData interface{}
}

func (s *Service) DumpNodes(args NodeArgs, reply *NodeReply) (err error) {
	reply.NodeData, err = s.nm.NodeMetric()
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
	links := s.lm.Links()
	data := make(map[string]interface{}, len(links))
	for _, link := range links {
		data[link], err = s.lm.LinkMetric(link)
		if err != nil {
			return err
		}
	}
	return nil
}
