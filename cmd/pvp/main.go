package main

import(
	"flag"
)

func main() {
	// 命令行参数给出
	lisPort := flag.String("l", "", "master listens at 8000")
	httpPort := flag.String("p", "", "http port at 80001")
	masterIP := flag.String("m", "", "connects to master with ip:port")
	flag.Parse()
}
