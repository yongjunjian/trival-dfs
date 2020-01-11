package service

import (
	//    . "trival/utils"
	. "trival/types"
)

func InitFreeBolck() error {
	return nil
}

type StoreWorker struct {
	//退出信号
	exit chan bool
	//当前在处理的对列
	reqQueue chan StoreReq
	//当前所在的分组
	groupId GroupId
	//当前所在的分区
	partition PartitionLabel
}

func NewStoreWorker(reqQueue chan StoreReq,
					groupId GroupId, 
					partition PartitionLabel) *StoreWorker {
	return &StoreWorker{
		reqQueue: reqQueue,
		groupId: groupId,
		partition: partition,
	}
}

func (this *StoreWorker) Run() {
	//检查是否转移分区的时间间隔
	var req StoreReq
	for {
		select {
		case <-this.exit:
			break
		case req = <-this.reqQueue:
			this.storeFile(req)
		}
	}
}

func (this *StoreWorker) storeFile(req StoreReq) error {
	//TODO
	return nil
}
