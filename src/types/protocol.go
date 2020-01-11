package types
import (
	"os"
)
type StoreArgs struct{
    GroupId int 
    FileName string
    Timestamp int64
    FileData []byte
}

type StoreReply struct{
    Id string
}

type RetrieveArgs struct{
    GroupId int 
    Partition string
    BlockId int
    Offset  int64
}

type RetrieveReply struct{
    Code int
    Message string
    FileName string
    FileData []byte
}

type StatArgs struct{
}

type StatInfo struct{
    GroupList []Group
    DiskUsed int
    DiskLimited int
}
type Group struct{
    PartitionList []Partition
}
type Partition struct{
    Date string
    LimitedBlockNum int
    BlockList []Block
}
type Block struct{
    ID int
    Size int
    FileNum int
    Handle *os.File
}
type StatReply struct{
    Statistic StatInfo
}

type SyncArgs struct{
    NodeId int
    Status int
}

type SyncReply struct{
}

