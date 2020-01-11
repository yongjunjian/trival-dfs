package service

import(
    . "trival/types"
    "container/list"
    "time"
	"fmt"
	"io/ioutil"
    "log"
    "errors"
	"strconv"
	"os"
	"strings"
    . "trival/utils"
)

const (
   CHECK_EXIT_INTERVAL  = 10
   DEFAULT_SELECT_TIMEOUT = 3
   INPUT_QUEUE_LENGTH = 4096
   OUTPUT_QUEUE_LENGTH = 1024
   BLOCK_FILE_EXT = ".blk"
)
var (
    dispatcher *StoreDispatcher
)
type StoreReq struct{
   args *StoreArgs
   reply *StoreReply
}

type StoreDispatcher struct{
    stopped bool
    input chan  StoreReq
    output map[PartitionLabel](chan StoreReq)
    //空闲分块表
    freeBlock map[PartitionLabel](chan *Block)
    storeWorker map[PartitionLabel](list.List)
}

func NewStoreDispatcher() *StoreDispatcher{
    //StoreDispatcher必须以单例模式运行
    if dispatcher == nil {
        input := make(chan StoreReq, INPUT_QUEUE_LENGTH)
        dispatcher = &StoreDispatcher{input:input}
    }
    return dispatcher
}

func (this *StoreDispatcher) initFreeBlock() error{
    //TODO
    return nil
}

func (this *StoreDispatcher) GetFreeBlock(groupId GroupId, partition PartitionLabel) (*Block, error){
    getMaxId := func(path string) (int, error){
		max := 0
        files, err := ioutil.ReadDir(path)
		if err != nil {
			log.Printf("directory not existed:%s", path)
			return -1, err
		}
    	for _, file := range files {
			if strings.HasSuffix(file.Name(), BLOCK_FILE_EXT){
				idStr := strings.TrimSuffix(file.Name(), BLOCK_FILE_EXT)
			   	id, err := strconv.Atoi(idStr)
				if err != nil {
					log.Printf("illeagal block file with bad name:%s/%s", path, file.Name())
				}else{
					if (id > max){
						max = id
					}
				}
			}
		}
		return max, nil
    }

    createBlock := func(partition PartitionLabel, groupId GroupId) (*Block, error){
        partitionPath := fmt.Sprintf("%s/%d/%s", Config().DataPath, groupId, partition)
        maxBlockId, err := getMaxId(partitionPath)
		if err != nil{
			log.Printf("get max id failed, path=%s", partitionPath)	
			return nil, err
		}
		blockId := maxBlockId+1
        blockPath := fmt.Sprintf("%s/%d.blk", partitionPath, blockId)
		handle, err := os.Create(blockPath)
        if  err != nil{
           log.Printf("create block file failed,path=%s, error=%v", 
                        blockPath, err)
           return nil, err 
        }
        return &Block{
                    ID: blockId,
                    Size: 0,
                    FileNum: 0,
                    Handle: handle, 
                  },nil
    }

 
    key := (PartitionLabel)(fmt.Sprintf("%d_%s", groupId, partition))
    select{
        case block := <- this.freeBlock[key]:
            return block, nil
        case <- time.After(DEFAULT_SELECT_TIMEOUT * time.Second ):
            if block, err:= createBlock( partition, groupId ); err != nil {
               return block, nil
            }else{
                return nil,err
            }
    }
}


func (this *StoreDispatcher) splitThread() error{
    split := func( req StoreReq ){
        timestamp := req.args.Timestamp
        partition := (PartitionLabel)(time.Unix(timestamp, 0).Format("20060102"))
        this.output[partition]  <- req
    }
    for{
        select {
            case req :=  <- this.input:
                split(req)
            case <- time.After(CHECK_EXIT_INTERVAL * time.Second ):
                if this.stopped && len(this.input) == 0{
                    log.Printf("split thread is ready to stop")
                    break
                }
            } 
    }
    return nil
}

func (this *StoreDispatcher) Start() error{
    if err := this.initFreeBlock(); err != nil{
        log.Printf("init free block failed!")
        return  err;
    }
    go this.splitThread()
    go this.adjustThread()
    return nil
}

func (this *StoreDispatcher) Stop(){
   this.stopped = false
   //TODO;等待所有请求处理完,即各分区的请求对列为空
   //TODO:通知所有storeWorker退出
}

//注意该接口会多线程调用，必须保证线程安全
func (this *StoreDispatcher) AddReq(args *StoreArgs, 
                                    reply * StoreReply) error{
    if this.stopped {
        return errors.New("dispatcher is stopped")
    } 
    this.input <-  StoreReq{args:args,reply:reply}
    return nil 
}

func (this *StoreDispatcher) adjustThread(){
    adjust := func(){
        //TODO
    }
    var interval = time.Duration(Config().Storage.AdjustInterval) * time.Second
    for{
        select{
            case <- time.After( interval):
                if this.stopped {
                    break
                }
                adjust()
        }
    } 
}

func (this *StoreDispatcher) clean(){}
