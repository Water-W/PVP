package influxdb

import (
	"context"
	"fmt"

	"time"

	"github.com/Water-W/PVP/pkg/biz"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	// protocol "github.com/influxdata/line-protocol"
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
		protocols := make(map[string]map[string]float64)
		nodes := make(map[string]string)
		for k, v := range v1.Reply.Node.(map[string]interface{}) {
			if k == "ID" {
				nodes[k] = v.(string)
			}
			if k == "Protocols" {
				for k1, v1 := range v.(map[string]interface{}) {
					p := make(map[string]float64)
					for k2, v2 := range v1.(map[string]interface{}) {
						p[k2] = v2.(float64)
					}
					protocols[k1] = p
				}
			}
		}

		links := make(map[string]map[string]float64)
		for k, v := range v1.Reply.Links {
			m := make(map[string]float64)
			for k1, v1 := range v.(map[string]interface{}) {
				m[k1] = v1.(float64)
			}
			links[k] = m
		}
		// 最后刷新数据库的写入
		defer writeAPI.Flush()
		// 每个link存一条，每个node 存一条，每个protocol存一条。
		nodename := nodes["ID"]
		for k := range links {
			p := influxdb2.NewPoint("dump",
				map[string]string{
					"kind": "edge",
					"from": nodename,
					"to":   k,
				},
				map[string]interface{}{
					"RateIn":   links[k]["RateIn"],
					"RateOut":  links[k]["RateOut"],
					"TotalIn":  links[k]["TotalIn"],
					"TotalOut": links[k]["TotalOut"],
				},
				time.Now())
			writeAPI.WritePoint(p)
		}
		for k := range protocols {
			ty := "other"
			name := k
			if k == "" {
				ty = "total"
				name = "total"
			}
			p := influxdb2.NewPoint("dump",
				map[string]string{
					"kind":     "node",
					"protocol": ty,
					"name": name,
				},
				map[string]interface{}{
					"RateIn":   protocols[k]["RateIn"],
					"RateOut":  protocols[k]["RateOut"],
					"TotalIn":  protocols[k]["TotalIn"],
					"TotalOut": protocols[k]["TotalOut"],
				},
				time.Now())
			writeAPI.WritePoint(p)
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

func Querydata() string {

	client := getclient()
	// timestart := "2021-05-11 07:00:00.850"
	// timestop := "2021-05-11 16:27:01.828"
	query := fmt.Sprintf("from(bucket: \"%v\") |> range(start: -2h) |> filter(fn: (r) => r[\"_measurement\"] == \"dump\")", bucket)

	
	// Get query client
	queryAPI := client.QueryAPI(org)
	// get QueryTableResult
	// result, err := queryAPI.Query(context.Background(), query)
	result, err := queryAPI.QueryRaw(context.Background(), query, influxdb2.DefaultDialect())
	if err == nil {
        fmt.Println("QueryResult:")
        fmt.Println(result)
    } else {
        panic(err)
    }
	defer client.Close()
	return result
}
