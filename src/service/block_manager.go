package service

import (
	. "trival/types"
    "time"
	"fmt"
	"io/ioutil"
    "log"
	"strconv"
	"os"
	"strings"
    . "trival/utils"
 
)
type BlockManager struct{
    //空闲分块表
    freeBlock map[GroupID](map[PartiID](chan *Block) )
}

func NewBlockManager() *BlockManager{
    bm := &BlockManager{
        freeBlock: make(map[GroupID](map[PartiID](chan *Block))),
    }
    bm.initFreeBlock()
    return bm
}
func (this *BlockManager) initFreeBlock() error{
	getFileNum := func(file *os.File) int{
		//读取文件总数
		return 0
	}

	dataPath := Config().Storage.DataPath
	groupDirs, err := ioutil.ReadDir(dataPath)
	if err != nil {
		log.Printf("directory not existed:%s", dataPath)
		return err
	}
	for _, groupDir := range groupDirs {
        tmp, err := strconv.Atoi(groupDir.Name())
        if err != nil{
            log.Printf("illeagal group directory:%s", groupDir.Name()) 
            continue
        }
        groupId := GroupID(tmp)
        groupDirPath := fmt.Sprintf("%s/%s", dataPath, groupDir.Name())
		if !groupDir.IsDir(){
			log.Printf("not a directory:%s",groupDirPath)
			continue
		}
    	partiDirs, err := ioutil.ReadDir(groupDirPath)
    	if err != nil {
    		log.Printf("directory not existed:%s", groupDirPath)
    		return err
    	}
        for _, partiDir := range partiDirs {
            partiId := PartiID(partiDir.Name())
            partiDirPath := fmt.Sprintf("%s/%s", groupDirPath, partiDir.Name()) 
            files, _ := ioutil.ReadDir(partiDirPath)
		    for _, file := range files {
		    	if strings.HasSuffix(file.Name(), BLOCK_FILE_EXT){
		    		idStr := strings.TrimSuffix(file.Name(), BLOCK_FILE_EXT)
		    		blockId, err := strconv.Atoi(idStr)
		    		if err != nil {
		    			log.Printf("illeagal block file with bad name:%s/%s", partiDir, file.Name())
		    			continue
		    		}

		    		path := fmt.Sprintf("%s/%s", dataPath, file.Name())
 		    		handle, err := os.Open(path)
		    		if err != nil{
		    			log.Printf("open file %s failed:%v", path, err)
                        continue
		    		}
                    if file.Size() <  Config().Storage.MaxBlockSize{
                        continue
                    }
                    this.freeBlock[groupId][partiId] <- &Block{
                        ID: blockId,
                        Size: file.Size(),
		    			FileNum: getFileNum(handle),
                        Handle: handle, 
                    }
		    	}
		    }
        }
    }
    return nil
}


func (this *BlockManager) AddFreeBlock(
						groupId GroupID, 
						partiId PartiID, 
						block* Block){
	    if block.Size >= Config().Storage.MaxBlockSize{
            return 	
	    }
        if _, existed := this.freeBlock[groupId]; !existed{
            this.freeBlock[groupId] = make(map[PartiID] chan *Block)
        }
        if _, existed := this.freeBlock[groupId][partiId]; !existed{
            this.freeBlock[groupId][partiId] = make(chan *Block)
        }
		this.freeBlock[groupId][partiId] <- block
}
func (this *BlockManager) GetFreeBlock(groupId GroupID, partiId PartiID) (*Block, error){
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
					continue
				}else{
					if (id > max){
						max = id
					}
				}
			}
		}
		return max, nil
    }

    createBlock := func(groupId GroupID,partition PartiID) (*Block, error){
        partitionPath := fmt.Sprintf("%s/%d/%s", Config().Storage.DataPath, groupId, partition)
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

 
    select{
        case block := <- this.freeBlock[groupId][partiId]:
            return block, nil
        case <- time.After(DEFAULT_SELECT_TIMEOUT * time.Second ):
            if block, err:= createBlock( groupId,partiId); err != nil {
               return block, nil
            }else{
                return nil,err
            }
    }
}


