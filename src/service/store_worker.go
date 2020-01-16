package service

import (
	//    . "trival/utils"
	. "trival/types"
    "log"
)

var (
    blockManager = NewBlockManager()
)

type StoreWorker struct {
	//退出信号
	exit chan bool
	//当前在处理的对列
	reqQueue chan StoreReq
	//当前所在的分组
	groupId GroupID
	//当前所在的分区
	partiId PartiID
	//当前所在的块
	block Block

}

func NewStoreWorker(reqQueue chan StoreReq,
					groupId GroupID, 
					partiId PartiID) *StoreWorker {
	return &StoreWorker{
		reqQueue: reqQueue,
		groupId: groupId,
		partiId: partiId,
	}
}

func (this *StoreWorker) Stop(){
	this.exit <- true
}

func (this *StoreWorker) Start() error{
	block, err := blockManager.GetFreeBlock(
							this.groupId, 
							this.partiId)
	if err != nil {
		log.Printf("get free block failed:%v", err)
		return err
	}
    log.Printf("get free block:%d", block.ID)	
	//检查是否转移分区的时间间隔
	go func(){
		var req StoreReq
		for {
			select {
			case <-this.exit:
				block.Handle.Sync()
				break
			case req = <-this.reqQueue:
				this.storeFile(req)
			}
		}
		log.Printf("store worker exit, group:%d, partition:%s, block:%d",
				this.groupId, this.partiId, block.ID)
		blockManager.AddFreeBlock(this.groupId, this.partiId, block)
	}()
	return nil
}

func (this *StoreWorker) storeFile(req StoreReq) error {
	//TODO
	return nil
}
