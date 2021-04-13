package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Water-W/PVP/pkg/biz"
	// "github.com/Water-W/PVP/pkg/metrics/json"
)

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "helloworld~")
}

func dumpHandler(w http.ResponseWriter, r *http.Request) {
	
}

func myhttp(ctrl *biz.MasterController) {
	//TODO原生
	//设置json模式，可以通过master给worker配置模式。
	//dump参数设置
	
	// select {}

	//获取dump的信息
	http.HandleFunc("/dump", dumpHandler)
	
	//设置json mode
	http.HandleFunc("/setjsonmode", set_json_mode)
	
	//loclahost:8080/hello
	http.HandleFunc("/hello", HelloHandler)
	
	//启动静态文件服务,可以访问localhosta:8080/frontend/main.html来查看
	http.Handle("/", http.FileServer(http.Dir("../")))
	http.ListenAndServe(":" + strconv.Itoa(httpPort), nil)
}

//设置json mode,使其应用在worker上
func set_json_mode(w http.ResponseWriter, r *http.Request) {
	//获取json mode


	//调用pkg json修改

}
