package analyze

import (
	"api-stress-testing/model"
	"fmt"
	"github.com/spenczar/tdigest"
	"sync"
	"time"
)

func AnalyzeResult(concurrent uint64, ch <-chan *model.Result, group *sync.WaitGroup) {
	defer group.Done()

	// 统计数据
	var (
		processingTime uint64 // 处理总时间
		requestTime    uint64 // 请求总时间
		maxTime        uint64 // 最大时长
		minTime        uint64 // 最小时长
		successNum     uint64 // 成功处理数，code为0
		failureNum     uint64 // 处理失败数，code不为0
		chanIdLen      int    // 并发数
		chanIds        = make(map[uint64]bool)
		td             = tdigest.New()
		//errCode        = make(map[int]int)
		startTime = uint64(time.Now().UnixNano())
		ticker    = time.NewTicker(1 * time.Second)
		stopChan  = make(chan bool)
	)

	go func() {
		for {
			select {
			case <-ticker.C:
				endTime := uint64(time.Now().UnixNano())
				requestTime = endTime - startTime
				go analyzeData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIdLen, td)
			case <-stopChan:
				return
			}
		}

	}()

	printTitle()

	for data := range ch {
		// fmt.Println("处理一条数据", data.Id, data.Time, data.IsSucceed, data.ErrCode)
		// 总耗时
		processingTime = processingTime + data.SpendTime
		// 记录最大执行时长
		if maxTime <= data.SpendTime {
			maxTime = data.SpendTime
		}

		// 记录最小执行时长
		if minTime > data.SpendTime || minTime == 0 {
			minTime = data.SpendTime
		}

		// 记录请求数据
		td.Add(float64(data.SpendTime), 1)

		// 是否请求成功
		if data.IsSuccess == true {
			successNum = successNum + 1
		} else {
			failureNum = failureNum + 1
		}

		if _, ok := chanIds[data.ChannelId]; !ok {
			chanIds[data.ChannelId] = true
			chanIdLen = len(chanIds)
		}
	}

	// 数据全部接受完成，停止定时输出统计数据
	stopChan <- true

	endTime := uint64(time.Now().UnixNano())
	requestTime = endTime - startTime

	analyzeData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum, chanIdLen, td)

	fmt.Printf("\n\n")

	fmt.Println("|--------------------------------------------|")
	fmt.Println("|--综合统计--")
	fmt.Println("|\t并发数：", concurrent)
	fmt.Println("|\t总请求数：", successNum + failureNum)
	fmt.Printf("|\t总请时间：%ds \n", requestTime/1e9)
	fmt.Printf("|\t请求成功率：%d \n", ((successNum) / (successNum + failureNum)) * 100)
	fmt.Printf("|\t平均响应时间：%8.2fms \n", float64(processingTime) / float64(successNum*1e6))
	fmt.Printf("|\tTP90：%8.2fms \n", td.Quantile(0.9) / 1e6)
	fmt.Printf("|\tTP95：%8.2fms \n", td.Quantile(0.95) / 1e6)
	fmt.Printf("|\tTP99：%8.2fms \n", td.Quantile(0.99) / 1e6)
	fmt.Println("|--------------------------------------------|")

}

func printTitle() {
	fmt.Println("|--------|--------|--------|--------|--------|--------|--------|--------|--------|--------|--------|")
	fmt.Println("|  时间  | 并发数 | 成功数 | 失败数 |   QPS  |   Max  |   Min  |   Avg  |  TP90  |  TP95  |  TP99  |")
	fmt.Println("|--------|--------|--------|--------|--------|--------|--------|--------|--------|--------|--------|")
}

func analyzeData(concurrent, processingTime, requestTime, maxTime, minTime, successNum, failureNum uint64, chanIdLen int, td *tdigest.TDigest) {
	if processingTime == 0 {
		processingTime = 1
	}

	var (
		qps              float64
		averageTime      float64
		maxTimeFloat     float64
		minTimeFloat     float64
		requestTimeFloat float64
	)

	// 平均 每个协程成功数*总协程数据/总耗时 (每秒)
	if processingTime != 0 {
		qps = float64(successNum*1e9*concurrent) / float64(processingTime)
	}

	// 平均时长 总耗时/总请求数/并发数 纳秒=>毫秒
	if successNum != 0 && concurrent != 0 {
		averageTime = float64(processingTime) / float64(successNum*1e6)
	}

	// TP90
	tp90 := td.Quantile(0.9) / 1e6
	// TP95
	tp95 := td.Quantile(0.95) / 1e6
	// TP99
	tp99 := td.Quantile(0.99) / 1e6

	// 纳秒=>毫秒
	maxTimeFloat = float64(maxTime) / 1e6
	minTimeFloat = float64(minTime) / 1e6
	requestTimeFloat = float64(requestTime) / 1e9

	result := fmt.Sprintf("|%7.0fs│%8d│%8d│%8d│%8.2f│%8.2f│%8.2f│%8.2f│%8.2f│%8.2f│%8.2f│",
		requestTimeFloat, chanIdLen, successNum, failureNum, qps, maxTimeFloat, minTimeFloat, averageTime,tp90, tp95, tp99)
	fmt.Println(result)
}
