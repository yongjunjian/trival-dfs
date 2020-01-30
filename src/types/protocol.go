package types
import (
	"os"
)
type StoreArgs struct{
    GroupId GroupID
    FileName string
    Timestamp int64
    FileData []byte
}

type StoreReply struct{
    Id string
    Err error
}

type RetrieveArgs struct{
    GroupId GroupID
    PartiId PartiID 
    BlockId BlockID
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
    PartiList []Parti
}
type Parti struct{
    Date string
    LimitedBlockNum int
    BlockList []Block
}
type Block struct{
    Id BlockID
    Size int64
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

type StoreReq struct{
   Args *StoreArgs
   Canceled *bool
   Done chan *StoreReply
}
