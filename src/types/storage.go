package types

type PartiID string
type BlockID uint16
type GroupID uint16

type BlockHeader struct{
    MagicNumber [8]byte 
    Version uint16
}

type Cell struct{
    TotalLen      int32
    FileNameLen   int16
    FileName  []byte
    Data []byte
    Delimeter []byte
}

type FileID struct{
    Version uint16
    NodeId uint16
    GroupId uint16
    Year uint16
    Month uint8
    Day uint8
    BlockId uint16
    Offset int64
}
