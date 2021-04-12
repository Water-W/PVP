package json

import (
	"fmt"
	"io/ioutil"

	"github.com/Water-W/PVP/pkg/metrics"

	"github.com/Water-W/PVP/pkg/rpc/dump"
	"github.com/tidwall/gjson"
)

var _ dump.NodeMeasurer = (*Measurer)(nil)
var _ dump.LinkMeasurer = (*Measurer)(nil)

// a measurer using https://github.com/tidwall/gjson to parse json
type Measurer struct {
	source     metrics.Source
	nodeQuery  string
	linksQuery string
}

func (m *Measurer) SetSource(s metrics.Source) {
	m.source = s
}

func (m *Measurer) SetNodeQuery(query string) {
	m.nodeQuery = query
}

func (m *Measurer) SetLinksQuery(query string) {
	m.linksQuery = query
}

func (m *Measurer) GetNode() (interface{}, error) {
	return m.get(m.nodeQuery)
}

func (m *Measurer) GetLinks() (map[string]interface{}, error) {
	out, err := m.get(m.linksQuery)
	if err != nil {
		return map[string]interface{}{}, err
	}
	links, ok := out.(map[string]interface{})
	if !ok {
		return map[string]interface{}{}, fmt.Errorf("links is not a json object")
	}
	return links, err
}

func (m *Measurer) get(query string) (interface{}, error) {
	bin, err := ioutil.ReadAll(m.source.Source())
	if err != nil {
		return nil, err
	}
	if !gjson.ValidBytes(bin) {
		return nil, fmt.Errorf("invalid json")
	}
	res := gjson.GetBytes(bin, query)
	if !res.IsObject() {
		return nil, fmt.Errorf("parse error")
	}
	return res.Value(), nil
}
