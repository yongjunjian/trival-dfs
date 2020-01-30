package service

import (
	. "trival/utils"
	. "trival/types"
    "encoding/binary"
    "log"
    "bytes"
    "time"
)

const (
    //储存字节总长度所用的字节数
    BYTES_OF_TOTAL_LEN = 4
    //储存文件名长度所用的字节数
    BYTES_OF_FILENAME_LEN =  2
    FILE_DELIMITER = "splt"
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
    log.Printf("get free block:%d", block.Id)	
	go func(){
		var req StoreReq
		for {
			select {
            case <-this.exit:
				block.Handle.Sync()
				break
			case req = <-this.reqQueue:
                if fileId,err  := this.storeFile(req); err != nil{
                    req.Done <- &StoreReply{
                        Err: err,
                    } 
                }else{
                     req.Done <- &StoreReply{
                        Id: fileId,
                        Err: nil,
                    } 
                }
			}
		}
		log.Printf("store worker exit, group:%d, partition:%s, block:%d",
				this.groupId, this.partiId, block.Id)
		blockManager.AddFreeBlock(this.groupId, this.partiId, block)
	}()
	return nil
}

func (this *StoreWorker) storeFile(req StoreReq) (string, error) {
    if this.groupId != req.Args.GroupId{
        log.Panic("groupId mismach")
    }
    
    writer :=  this.block.Handle
    var cell Cell
    cell.FileName = ([]byte)(req.Args.FileName)
    cell.FileNameLen = int16(len(cell.FileName))
    cell.Data  = req.Args.FileData
    cell.Delimeter = ([] byte)(FILE_DELIMITER)
    cell.TotalLen = int32(int(cell.FileNameLen)+len(cell.Data)+BYTES_OF_TOTAL_LEN+BYTES_OF_FILENAME_LEN)
    offsetStart, err := writer.Seek(0, 1)
    if err != nil {
        log.Printf("seek offset failed:%v", err)
        return "", err
    }
    offsetEnd, err := writer.Seek(0, 1)
    realTotalLen := int32(offsetStart-offsetEnd)
    if realTotalLen != cell.TotalLen{
        log.Panicf("write bytes is not expected: %d vs. %d",realTotalLen, cell.TotalLen ) 
    }
    if err := binary.Write(writer, binary.LittleEndian, cell); err != nil{
        return "", err
    } 
    generateId := func(version int, 
                        nodeId int,
                        groupId GroupID,
                        timestamp int64,
                        blockId BlockID,
                        offset int64) string{
        var fileId FileID
        fileId.Version = uint16(version)
        fileId.NodeId = uint16(nodeId)
        fileId.GroupId = uint16(groupId)
        fileId.BlockId = uint16(blockId)
        fileId.Offset = offset
        ts := time.Unix(timestamp/1000, 0)
        fileId.Year = uint16(ts.Year())
        fileId.Month = uint8(ts.Month())
        fileId.Day = uint8(ts.Day())
        bufWriter := bytes.NewBufferString("") 
        if err := binary.Write(bufWriter, binary.LittleEndian, fileId); err != nil{
            log.Panicf("write string to buffer err:%v", err)
        }
        return  bufWriter.String()
    }

    id := generateId(DATABASE_VERSION, 
                    Config().Storage.DataNodeId,
                    this.groupId,
                    req.Args.Timestamp,
                    this.block.Id,
                    offsetStart,
                )
	return id,nil
}
