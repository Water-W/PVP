package json

import (
	"fmt"
	"io/ioutil"

	"github.com/Water-W/PVP/pkg/metrics"

	"github.com/Water-W/PVP/pkg/rpc/dump"
	"github.com/tidwall/gjson"
)

var _ dump.NodeMeasurer = (*JsonMeasurer)(nil)
var _ dump.LinkMeasurer = (*JsonMeasurer)(nil)

// a measurer using https://github.com/tidwall/gjson to parse json
type JsonMeasurer struct {
	source  metrics.Source
	nodeQL  string
	linksQL string
}

func (m *JsonMeasurer) SetSource(s metrics.Source) {
	m.source = s
}

func (m *JsonMeasurer) SetNodeQL(query string) {
	m.nodeQL = query
}

func (m *JsonMeasurer) SetLinksQL(query string) {
	m.linksQL = query
}

func (m *JsonMeasurer) GetNode() (interface{}, error) {
	return m.get(m.nodeQL)
}

func (m *JsonMeasurer) GetLinks() (map[string]interface{}, error) {
	out, err := m.get(m.linksQL)
	return out.(map[string]interface{}), err
}

func (m *JsonMeasurer) get(query string) (interface{}, error) {
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
