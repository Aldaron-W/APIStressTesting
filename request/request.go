package request

import (
	"api-stress-testing/model"
	"errors"
	"sync"
)

type Request interface {
	/**
	 * 发送请求
	 * 参数：
	 * id			请求编号
	 * request 		请求内容对象
	 * callbackCh	请求结果响应Channel
	 * times		请求次数
	 * group		请求WaitGroup
	 */
	SendRequest(id uint64, request model.Request, callbackCh chan <- *model.Result, times uint64, group *sync.WaitGroup)
}

func SenderFactory(protocol string)(request Request, err error)  {
	switch protocol {
	case "http":
		request = new(HttpRequestSender)
		return
	}

	return nil, errors.New("没有实现本协议的Sender")
}