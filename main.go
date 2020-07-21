package main

import (
	"api-stress-testing/model"
	"flag"
)

var (
	concurrency uint64	// 并发数(个)
	requestTime uint64	// 单个协程请求数(次)
	requestURI string	// 请求地址。当前支持(http/https)
)

func main() {
	flag.Uint64Var(&concurrency, "c", 10, "并发数(个)")
	flag.Uint64Var(&requestTime, "t", 100, "单个协程请求数(次)")
	flag.StringVar(&requestURI, "u", "https://www.baidu.com", "请求地址。当前支持(http/https)")

	flag.Parse()

	//Request
	httpRequest := &model.HTTPRequest{
		Method: "GET",
		URL:  requestURI ,
	}

	manager(httpRequest, concurrency, requestTime)

	return
}
