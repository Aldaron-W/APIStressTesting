package request

import (
	model "api-stress-testing/model"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type HttpRequestSender struct {

}

func (h *HttpRequestSender) SendRequest(id uint64, request model.Request, callbackCh chan<- *model.Result, times uint64, group *sync.WaitGroup) {
	//fmt.Printf("第%d协程开始工作\n", id)
	defer group.Done()
	httpRequest, ok := request.(*model.HTTPRequest)
	for i := uint64(0); i < times; i++ {
		var (
			isSuccess = true
			errorCode = 0;
		)

		if ok {
			resp, requestTime, err := h.send(httpRequest.URL, httpRequest.Method, httpRequest.GetBody(), httpRequest.Headers, httpRequest.Timeout)

			if err != nil {
				isSuccess = false
			}

			if resp == nil {
				errorCode = -1
			} else {
				errorCode = resp.StatusCode
			}

			requestResult := &model.Result{
				Id: i,
				ChannelId: id,
				IsSuccess: isSuccess,
				SpendTime: requestTime,
				ErrorCode: errorCode,
			}

			callbackCh <- requestResult
		}
	}
}

// 发送HTTP请求
func (h *HttpRequestSender) send(url string, method string, body io.Reader, headers map[string]string, timeout time.Duration) (resp *http.Response, requestTime uint64, err error) {
	// 跳过证书验证
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {

		return
	}

	// 设置默认为utf-8编码
	if _, ok := headers["Content-Type"]; !ok {
		if headers == nil {
			headers = make(map[string]string)
		}
		headers["Content-Type"] = "application/x-www-form-urlencoded; charset=utf-8"
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err = client.Do(req)
	requestTime = uint64(time.Now().UnixNano() - startTime.UnixNano())
	if err != nil {
		fmt.Println("请求失败:", err)

		return
	}

	return
}


