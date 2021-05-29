package influxdb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Water-W/PVP/pkg/biz"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	// "encoding/json"
)

// You can generate a Token from the "Tokens Tab" in the UI
const token = "SEzvFcp4YjFz8zH_FLnnJp3OAhOCZjQj6bKG5wNiyzQgCzqINWUBLM_y1eexPAK-HY1GZOd7ecOJLYn8Nz6DUQ=="
const bucket = "bucket"
const org = "pku"

func getclient() influxdb2.Client {
	client := influxdb2.NewClient("http://localhost:8086", token)
	return client
}

var name2id map[string]int
var num = 0

type myprotocol struct {
	Name     string
	RateIn   float64 `json:"RateIn"`
	RateOut  float64
	TotalIn  float64
	TotalOut float64
}
type node struct {
	Id        int
	Name      string
	Target    []int
	Protocols []myprotocol
}
type link struct {
	Sourceid int
	Targetid int
	RateIn   float64
	RateOut  float64
	TotalIn  float64
	TotalOut float64
}
type tdMap struct {
	Nodes []node
	Links []link
}

func pushNode(str string) bool {
	_, ok := name2id[str]
	if ok {
		return true
	} else {
		name2id[str] = num
		num += 1
		return false
	}
}
func write_dump(client influxdb2.Client, data []biz.DumpResult) {
	thistime := time.Now()
	thisUnixtime := fmt.Sprintf("%v", thistime.Unix())
	writeAPI := client.WriteAPI(org, bucket)
	defer writeAPI.Flush()
	//整和成 {nodes, links}存入

	name2id = make(map[string]int)
	num = 0
	var mytdMap tdMap
	var nodes []node
	var links []link
	// 同一个时刻各个服务器传来消息data[]
	for _, v1 := range data {
		// 对于每个服务器的数据
		// ipdatafrom := v1.From (数据来源ip)
		// fmt.Print(k1, "\n")
		// fmt.Print(v1, "\n")
		// protocols nodename links
		var myprotocols []myprotocol
		var mytarget []int
		mytarget = []int{}
		nodename := ""
		var nodeid int
		for k, v := range v1.Reply.Node.(map[string]interface{}) {
			if k == "ID" {
				nodename = v.(string)
			}
			if k == "Protocols" {
				for k1, v1 := range v.(map[string]interface{}) {
					i := v1.(map[string]interface{})
					myprotocols = append(myprotocols, myprotocol{
						Name:     k1,
						RateIn:   i["RateIn"].(float64),
						RateOut:  i["RateOut"].(float64),
						TotalIn:  i["TotalIn"].(float64),
						TotalOut: i["TotalOut"].(float64),
					})
				}
			}
		}
		edges := make(map[string]map[string]float64)
		for k, v := range v1.Reply.Links {
			m := make(map[string]float64)
			for k1, v1 := range v.(map[string]interface{}) {
				m[k1] = v1.(float64)
			}
			edges[k] = m
		}

		// 先加入本节点
		if !pushNode(nodename) {
			nodes = append(nodes, node{
				Id:        num - 1,
				Name:      nodename,
				Protocols: myprotocols,
			})
			nodeid = num - 1
		} else {
			nodeid = name2id[nodename]
		}
		// links {key{in ,out{float}}}
		for k2, v2 := range edges {
			if !pushNode(k2) {
				nodes = append(nodes, node{
					Id:   num - 1,
					Name: nodename,
				})
			}
			links = append(links, link{
				Sourceid: nodeid,
				Targetid: name2id[k2],
				RateIn:   v2["RateIt"],
				RateOut:  v2["RateOut"],
				TotalIn:  v2["TotalIn"],
				TotalOut: v2["TotalOut"],
			})
			mytarget = append(mytarget, name2id[k2])
		}
		// fmt.Print(nodeid, "nodeid\n")
		// fmt.Print(mytarget, "\ntarget\n")
		nodes[nodeid].Target = mytarget
	}
	mytdMap = tdMap{Nodes: nodes, Links: links}
	j, err := json.Marshal(mytdMap)
	if err != nil {
		fmt.Print(err)
		return
	}
	// fmt.Print("json\n")
	// fmt.Print(string(j))
	// fmt.Print("\n")
	p := influxdb2.NewPoint("dump",
		map[string]string{
			"kind":     "tdMap",
			"unixtime": thisUnixtime,
		},
		map[string]interface{}{
			"tdMap": j,
		},
		thistime)
	writeAPI.WritePoint(p)
}

func Storedata(data []biz.DumpResult) {

	client := getclient()

	// create point and write data
	write_dump(client, data)

	// always close client at the end
	defer client.Close()
}

func Querydata() ([]map[string]interface{}, error) {
	client := getclient()
	querystring := fmt.Sprintf("from(bucket: \"%v\") |> range(start: -5m) |> filter(fn: (r) => r[\"_measurement\"] == \"dump\")", bucket)

	// Get query client
	queryAPI := client.QueryAPI(org)
	// get QueryTableResult
	fmt.Print("开始查\n")
	result, err := queryAPI.Query(context.Background(), querystring)
	fmt.Print("结束查\n")

	//查询到的每条记录
	var records []map[string]interface{}
	if err == nil {
		// 对于每一条记录进行处理

		// _, ok := map[key]
		fmt.Print("开始处理\n")
		for result.Next() {
			records = append(records, result.Record().Values())
		}
		fmt.Print("结束处理\n")

		fmt.Print("\n输出records\n")
		for kk, vv := range records {
			fmt.Print(kk)
			fmt.Print("\n")
			fmt.Print(vv)
			fmt.Print("\n")

		}
		fmt.Print("\n结束\n")

		// check for an error
		if result.Err() != nil {
			fmt.Printf("query parsing error: %s\n", result.Err().Error())
		}

	} else {
		panic(err)
	}
	//将查询到的数据进行处理

	defer client.Close()
	return records, err
}

// result, err := queryAPI.QueryRaw(context.Background(), query, influxdb2.DefaultDialect())

// if err == nil {
// 	fmt.Println("QueryResult:")
// 	fmt.Println(result)
// } else {
// 	panic(err)
// }
