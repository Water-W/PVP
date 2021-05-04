package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	// "github.com/tidwall/gjson"

	"github.com/Water-W/PVP/pkg/biz"
	// "github.com/Water-W/PVP/pkg/metrics/json"
)

type server struct {
	ctrl *biz.MasterController
}

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "helloworld~")
}

func (s *server) dump(w http.ResponseWriter, r *http.Request) {
	request, err := s.ctrl.Dump(context.Background())
	if err != nil {
		log.Error(err)
		return
	}
	// 转换dump的结果为标准json
	j, err := json.Marshal(request)
	if err != nil {
		log.Error(err)
		return
	}
	w.Write(j)
}

//设置json mode,使其应用在worker上
func (s *server) setjm(w http.ResponseWriter, r *http.Request) {
	//获取json mode
	//调用pkg json修改
	
}

func myhttp(ctrl *biz.MasterController) {
	//TODO原生
	//设置json模式，可以通过master给worker配置模式。
	//dump参数设置
	s := server{ctrl: ctrl}
	// select {}

	//获取dump的信息
	http.HandleFunc("/dump", s.dump)

	//设置json mode
	http.HandleFunc("/setjsonmode", s.setjm)

	//loclahost:8080/hello
	http.HandleFunc("/hello", HelloHandler)

	//启动静态文件服务,可以访问localhosta:8080/frontend/main.html来查看
	http.Handle("/", http.FileServer(http.Dir("../")))
	
	http.ListenAndServe(":"+strconv.Itoa(*httpPort), nil)
}
