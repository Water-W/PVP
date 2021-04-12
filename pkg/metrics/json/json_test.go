package json

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleJSON = `
{
	"ID": "12D3KooWK5qUN9LnjhrfPEZ8pBBQje4uNsCq7M451EhjX3B5465y",
	"Peers": {
	  "12D3KooWEZTcPNnpiYyQnrqvUUvxAj3ubdrBHya5ybT3LhoejHaj": {
		"TotalIn": 0,
		"TotalOut": 72,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "12D3KooWEcRenx7qaFc2SQz5j2XtyPmgiAQjEBacSFLwSPTB1jjY": {
		"TotalIn": 26602,
		"TotalOut": 26536,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "12D3KooWJ8zusgRGAN2Uf8BZzPrmqreuCyKq4JxAhknZVdtBPD1G": {
		"TotalIn": 28459,
		"TotalOut": 28006,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "12D3KooWJmjesbNxBo852J4riWYuWixiG5YBgySu6DyTrGPx4nZ1": {
		"TotalIn": 26220,
		"TotalOut": 27249,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "12D3KooWJmrrBY3vLHSbTJ1sWmoXnuo4iF3defQnxo7rmdf8qisH": {
		"TotalIn": 26567,
		"TotalOut": 27824,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "12D3KooWK17ZSFhbuCNxuvdQH5Q4wgTN9kJDku2rqBx6xxiz4Kju": {
		"TotalIn": 0,
		"TotalOut": 36,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "12D3KooWNKitJ4dfkqzuqqfxPqTUXSPKJm6ShF45RVotUG4d8uHA": {
		"TotalIn": 26048,
		"TotalOut": 26372,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "12D3KooWNzuhB4oeykHiTVd2wT8FCqyPHMhtnaCyJCRb3ZLK3Uqo": {
		"TotalIn": 0,
		"TotalOut": 36,
		"RateIn": 0,
		"RateOut": 0
	  }
	},
	"Protocols": {
	  "": {
		"TotalIn": 2695,
		"TotalOut": 1420,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/fastboot/0.0.1": {
		"TotalIn": 340,
		"TotalOut": 1214,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/ipfs/id/1.0.0": {
		"TotalIn": 3170,
		"TotalOut": 3404,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/ipfs/id/push/1.0.0": {
		"TotalIn": 6122,
		"TotalOut": 3515,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/ipfs/kad/1.0.0": {
		"TotalIn": 8665,
		"TotalOut": 13230,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/libp2p/autonat/1.0.0": {
		"TotalIn": 719,
		"TotalOut": 760,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/meshsub/1.1.0": {
		"TotalIn": 111740,
		"TotalOut": 111343,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/p2p/id/delta/1.0.0": {
		"TotalIn": 315,
		"TotalOut": 515,
		"RateIn": 0,
		"RateOut": 0
	  },
	  "/topic-wires/1.0.0": {
		"TotalIn": 130,
		"TotalOut": 730,
		"RateIn": 0,
		"RateOut": 0
	  }
	}
}`

type mockSource struct{}

func (m *mockSource) Source() io.Reader {
	return bytes.NewReader([]byte(exampleJSON))
}

func TestJsonMeasurer(t *testing.T) {
	expNode := map[string]interface{}{
		"ID": "12D3KooWK5qUN9LnjhrfPEZ8pBBQje4uNsCq7M451EhjX3B5465y",
	}
	expLinks := map[string]interface{}{
		"12D3KooWEZTcPNnpiYyQnrqvUUvxAj3ubdrBHya5ybT3LhoejHaj": map[string]interface{}{
			"TotalIn":  float64(0),
			"TotalOut": float64(72),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
		"12D3KooWEcRenx7qaFc2SQz5j2XtyPmgiAQjEBacSFLwSPTB1jjY": map[string]interface{}{
			"TotalIn":  float64(26602),
			"TotalOut": float64(26536),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
		"12D3KooWJ8zusgRGAN2Uf8BZzPrmqreuCyKq4JxAhknZVdtBPD1G": map[string]interface{}{
			"TotalIn":  float64(28459),
			"TotalOut": float64(28006),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
		"12D3KooWJmjesbNxBo852J4riWYuWixiG5YBgySu6DyTrGPx4nZ1": map[string]interface{}{
			"TotalIn":  float64(26220),
			"TotalOut": float64(27249),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
		"12D3KooWJmrrBY3vLHSbTJ1sWmoXnuo4iF3defQnxo7rmdf8qisH": map[string]interface{}{
			"TotalIn":  float64(26567),
			"TotalOut": float64(27824),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
		"12D3KooWK17ZSFhbuCNxuvdQH5Q4wgTN9kJDku2rqBx6xxiz4Kju": map[string]interface{}{
			"TotalIn":  float64(0),
			"TotalOut": float64(36),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
		"12D3KooWNKitJ4dfkqzuqqfxPqTUXSPKJm6ShF45RVotUG4d8uHA": map[string]interface{}{
			"TotalIn":  float64(26048),
			"TotalOut": float64(26372),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
		"12D3KooWNzuhB4oeykHiTVd2wT8FCqyPHMhtnaCyJCRb3ZLK3Uqo": map[string]interface{}{
			"TotalIn":  float64(0),
			"TotalOut": float64(36),
			"RateIn":   float64(0),
			"RateOut":  float64(0),
		},
	}

	j := &Measurer{
		source:     &mockSource{},
		nodeQuery:  `{ID}`,
		linksQuery: `Peers`,
	}

	actNode, err := j.GetNode()
	assert.NoError(t, err)
	assert.Equal(t, expNode, actNode)
	actLinks, err := j.GetLinks()
	assert.NoError(t, err)
	assert.Equal(t, expLinks, actLinks)
}
