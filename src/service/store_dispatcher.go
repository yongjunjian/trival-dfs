package service

import(
    . "trival/types"
    "container/list"
    "time"
    "log"
    "errors"
    . "trival/utils"
    "sync"
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


type StoreDispatcher struct{
    stopped bool
    input chan  *StoreReq
    output map[GroupID](map[PartiID] chan *StoreReq)
    mapLock sync.Mutex
    storeWorkers map[GroupID](map[PartiID] *list.List)

}

func NewStoreDispatcher() *StoreDispatcher{
    //StoreDispatcher必须以单例模式运行
    if dispatcher == nil {
        dispatcher = &StoreDispatcher{
            input: make(chan *StoreReq, INPUT_QUEUE_LENGTH),
            output: make(map[GroupID](map[PartiID] chan *StoreReq )),
            mapLock: sync.Mutex{},
            storeWorkers: make(map[GroupID](map[PartiID] *list.List)),
        }
    }
    return dispatcher
}

func (this *StoreDispatcher) getPartiId(timestamp int64) PartiId{
    timestamp := req.Args.Timestamp
    return (PartiID)(time.Unix(timestamp, 0).Format("20060102"))
 }

func (this *StoreDispatcher) splitThread() error{
    split := func( req *StoreReq ){
        groupId := req.Args.GroupId
           if _, existed :=  this.output[groupId]; !existed{
            this.mapLock.Lock()
            this.output[groupId] = make(map[PartiID] chan *StoreReq)
            this.mapLock.Unlock()
        }
        partiId := getPartiId(req.Args.Timestamp)
        if _, existed :=  this.output[groupId][partiId]; !existed{
            this.mapLock.Lock()
            this.output[groupId][partiId] = make(chan *StoreReq, OUTPUT_QUEUE_LENGTH) 
            this.mapLock.UnLock()
        }
        this.output[groupId][partiId]  <- req
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
    go this.splitThread()
    go this.adjustThread()
    this.stopped = false
    return nil
}

func (this *StoreDispatcher) Stop(){
   this.stopped = true
}

func (this *StoreDispatcher) initWorkerIfNotExist(groupId GroupID, partiId PartiID) bool{
    executed =  false
    if _, existed := this.storeWorkers[groupId]; !existed{
        this.workerLock.Lock()
        this.storeWorkers[groupId] = make(map[PartiID] *list.List)
        this.workerLock.UnLock()
        executed = true
    }
    if _, existed := this.storeWorkers[groupId][partiId]; !existed {
        this.workerLock.Lock()
        this.storeWorkers[groupId][partiId] = list.New().Init()
        this.workerLock.UnLock()
        executed = true
    }
    return executed
}

func (this *StoreDispatcher) addWorkerIfNotExist(groupId GroupID, partiId PartiID){
    if this.initWorker(groupId, partiId) {
        sw := NewStoreWorker(
            this.output[groupId][partiId],
            groupId, 
            partiId)
        sw.Start()
        this.workerLock.Lock()
        this.storeWorkers[groupId][partiId].PushFront(sw)
        this.workerLock.UnLock()
    }
}

//注意该接口会多线程调用，必须保证线程安全
func (this *StoreDispatcher) SyncDispatch(args *StoreArgs)* StoreReply {
    if this.stopped {
        return &StoreReply{Err: errors.New("store dispatcher is stopped")}
    } 
    done := make(chan *StoreReply)
    canceled := false
    this.input <-  &StoreReq{Args:args, Done:done, Canceled: &canceled}
    this.addWorkerIfNotExist(args.GroupId, getPartiID(args.Timestamp))
    select{
            case <- time.After(time.Duration(Config().Storage.Timeout)*time.Second):
                log.Printf("timeout to store")
                canceled = true       
                return &StoreReply{Err: errors.New("timeout to wait for storing")}
            case reply := <- done:
                return reply
    }
    return nil
}

func (this *StoreDispatcher) adjustThread(){
	noticeAllExit := func(groupId GroupID, partiId PartiID){
        if  _, existed := this.storeWorkers[groupId];  !existed{
            return
        }
        if  _, existed := this.storeWorkers[groupId][partiId]; !existed{
            return
        }
        queue := this.storeWorkers[groupId][partiId]
        for elem := queue.Front(); elem != nil; elem = elem.Next(){
           elem.Value.(*StoreWorker).Stop()
		}
	}
    adjust := func(){
        //统计每个分区请求队列长度
		//根据队列长度比例，计算每个分区应该拥有的线程数
		//统计每个分区中存活的线程数，顺带清理map中已退出的线程
		total := 0
		queueLen := make(map[GroupID](map[PartiID] int))
        for groupId, queueMap := range this.output{
            for partiId, queue := range queueMap{
			    queueLen[groupId][partiId] = len(queue)
		 	    total += queueLen[groupId][partiId] 
		    }
        }
        threadNumTotal := Config().Storage.ThreadNum
        for groupId, lengthMap := range queueLen{
            for partiId, length := range lengthMap{
    			if length == 0 {
                    noticeAllExit(groupId, partiId)
                    if _, existed := this.storeWorkers[groupId]; existed{
                        if _, existed := this.storeWorkers[groupId][partiId]; existed{
                            delete(this.storeWorkers[groupId], partiId)
                        }
                    }

                    this.mapLock.Lock()
				    delete(this.output[groupId], partiId)
                    this.mapLock.Unlock()
                    continue
			    }

                threadNum := threadNumTotal*length/total
                if threadNum == 0{
                    threadNum = 1 //要保证至少有一个线程在处理分片
                }
                if _, existed := this.storeWorkers[groupId]; !existed{
                    this.storeWorkers[groupId] = make(map[PartiID] *list.List)
                }
                if _, existed := this.storeWorkers[groupId][partiId]; !existed {
                    this.storeWorkers[groupId][partiId] = list.New().Init()
                }
                swList := this.storeWorkers[groupId][partiId]
                for swList.Len() < threadNum {
                    sw := NewStoreWorker(
                        this.output[groupId][partiId],
                        groupId, 
                        partiId)
                    sw.Start()
                    swList.PushFront(sw);
                }
                for swList.Len() > threadNum {
                    elem := swList.Front() 
                    elem.Value.(*StoreWorker).Stop()
                    swList.Remove(elem)
                }
            } 
        }
    }
    interval := time.Duration(Config().Storage.AdjustInterval) * time.Second
    for {
        select {
            case <- time.After( interval):
                if this.stopped {
                    break
                }
                adjust()
        }
    } 
}

func (this *StoreDispatcher) clean(){}
