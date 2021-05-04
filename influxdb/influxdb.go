package influxdb

import (
	"context"
	"fmt"

	"github.com/Water-W/PVP/pkg/biz"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"time"
	// "encoding/json"
)

// You can generate a Token from the "Tokens Tab" in the UI
const token = "DviVETfMVeZcyF1gbVhw5ibuac2-3z5vQynr8D50B9p8RQlYr7lc5qPo8lBgtRrP1M5JTNolctbJABi-3W27RQ=="
const bucket = "bucket1"
const org = "pku"

func getclient() influxdb2.Client {
	client := influxdb2.NewClient("http://localhost:8086", token)
	return client
}

func write_dump(client influxdb2.Client, data []biz.DumpResult) {
	writeAPI := client.WriteAPI(org, bucket)
	for _, v1 := range data {
		// ipdatafrom := v1.From (数据来源ip)

		// 断言interface
		nodes := make(map[string]string)
		for k, v := range v1.Reply.Node.(map[string]interface{}) {
			nodes[k] = v.(string)
		}
		links := make(map[string]map[string]float64)
		for k, v := range v1.Reply.Links {
			m := make(map[string]float64)
			for k1, v1 := range v.(map[string]interface{}) {
				m[k1] = v1.(float64)
			}
			links[k] = m
		}
		// fmt.Printf("%#v\n", links)
		nodename := nodes["ID"]
		for s2 := range v1.Reply.Links {
			p := influxdb2.NewPoint("dump",
				map[string]string{
					"kind": "edge",
					"from": nodename,
					"to":   s2,
				},
				map[string]interface{}{
					"RateIn":   links[s2]["RateIn"],
					"RateOut":  links[s2]["RateOut"],
					"TotalIn":  links[s2]["TotalIn"],
					"TotalOut": links[s2]["TotalOut"],
				},
				time.Now())
			writeAPI.WritePoint(p)
			writeAPI.Flush()
		}
	}
}

func Storedata(data []biz.DumpResult) {

	client := getclient()

	// create point and write data
	write_dump(client, data)

	// always close client at the end
	defer client.Close()
}

func Querydata() {

	client := getclient()
	query := fmt.Sprintf("from(bucket:\"%v\")|> range(start: -1h) |> filter(fn: (r) => r._measurement == \"stat\")", bucket)
	// Get query client
	queryAPI := client.QueryAPI(org)
	// get QueryTableResult
	result, err := queryAPI.Query(context.Background(), query)

	if err == nil {
		// Iterate over query response
		for result.Next() {
			// Notice when group key has changed
			if result.TableChanged() {
				fmt.Printf("table: %s\n", result.TableMetadata().String())
			}
			// Access data
			fmt.Printf("value: %v\n", result.Record().Value())
		}
		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %v\n", result.Err().Error())
		}
	} else {
		panic(err)
	}

	defer client.Close()
}
