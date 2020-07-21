package model

type Result struct {
	Id uint64
	ChannelId uint64
	IsSuccess bool
	SpendTime uint64
	ErrorCode int
}