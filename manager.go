package main

import (
	"api-stress-testing/analyze"
	"api-stress-testing/model"
	request2 "api-stress-testing/request"
	"sync"
)

func manager(request model.Request, concurrency uint64, requestTime uint64) {
	requestSender, err := request2.SenderFactory(request.GetProtocol())
	if err != nil {
		return
	}
	var (
		ch = make(chan *model.Result, 20)
		senderGroup sync.WaitGroup
		analyzeGroup sync.WaitGroup
	)

	// 启动收集与分析协程
	analyzeGroup.Add(1)
	go analyze.AnalyzeResult(concurrency, ch, &analyzeGroup)

	// 启动并发协程发送请求
	for i := uint64(0); i < concurrency; i++ {
		senderGroup.Add(1)
		go requestSender.SendRequest(i, request, ch, requestTime, &senderGroup)
	}

	senderGroup.Wait();

	close(ch)

	analyzeGroup.Wait();
	return
}