package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/Water-W/PVP/pkg/biz"
	"github.com/tidwall/gjson"
)

func TestWorkerParse(t *testing.T) {
	c := &biz.WorkerConfig{
		URL:        "http://39.104.200.8:2404/report",
		NodeQuery:  `{ID,Protocols}`,
		LinksQuery: `Peers`,
	}
	resp, err := http.Get(c.URL)
	if err != nil {
		panic(err)
	}
	bin, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bin))
	res := gjson.GetBytes(bin, c.NodeQuery)
	fmt.Println(res.String())
}
