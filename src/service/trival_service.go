package service

import (
	. "trival/types"
	//    metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/server"
	//    "github.com/smallnest/rpcx/serverplugin"
	"context"
	"fmt"
	"log"
	. "trival/utils"
)

const (
	REGISTRY_PATH = "/trival"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

/*
- 功能:储存文件到本地
- 参数
  - groupId: 分组id
  - fileName:文件名
  - fileData:文件数据
- 返回
  - 出错信息
  - 文件id
**/
func (this *Service) Store(ctx context.Context,
	args *StoreArgs,
	reply *StoreReply) error {
	err, fileId := this.localStore(args.GroupId,
		args.FileName,
		args.Timestamp,
		args.FileData)
	if err != nil {
		log.Printf("store file data failed:%v", err)
		return err
	} else {
		reply.Id = fileId
		return nil
	}
}

func (this *Service) localStore(groupId int,
	fileName string,
	timestamp int64,
	fileData []byte) (error, string) {

	return nil, ""
}

/*
  - 功能:从本地读取文件
  - 参数
    - groupId: 分组id
    - partition: 分区，格式：YYYYMMDD
    - blockId: 分块id
    - offset:文件起始位置在分块中的偏移量
  - 返回
    - 出错信息
    - 文件名
    - 文件数据
**/
func (this *Service) localRetrieve(groupId int,
	partition string,
	blockId int,
	offset int64) (error, string, []byte) {
	return nil, "", nil
}

func (this *Service) Retrieve(ctx context.Context,
	args *RetrieveArgs,
	reply *RetrieveReply) error {
	return nil
}

/*
- 功能: 统计本地的状态信息
  - 参数:无
  - 返回
    - 出错信息
    - 状态信息数据
*/
func (this *Service) localStat() (error, *StatInfo) {
	return nil, nil
}

func (this *Service) Stat(ctx context.Context,
	args *StatArgs,
	reply *StatReply) error {
	return nil
}

/*
- 功能: 接受其它节点广播过来的状态信息,用于请求调度
  - 参数:
    nodeId: 结点编号
    status:服务当前的状态,0-不可读不可写 1-仅可读 2-仅可写 3-可读写
  - 返回
    - 出错信息
    - 状态信息数据
 - 注意;该函数会被频繁调用，响应时间必须要短
*/
func (this *Service) syncStatus(nodeId int, status int) error {
	return nil
}

func (this *Service) SyncStatus(ctx context.Context,
	args *SyncArgs,
	reply *SyncReply) error {
	return nil
}

func ServeRPC() {
	service := NewService()
	srv := server.NewServer()
	addr := fmt.Sprintf("%s:%d", Config().Rpc.IP, Config().Rpc.Port)
	srv.Register(service, "")
	log.Printf("start rpc service, listen on:%s", addr)
	err := srv.Serve("tcp", addr)
	if err != nil {
		log.Fatalf("start rpc service failed", err)
	}
}
