### 针对小文件的分布式储存
#### 功能需求
1. 文件合并储存，减少inode结点的使用，并提升磁盘空间利用率
2. 支持按天批量删除文件
3. 支持分布式储存，最大支持100T
4. 支持高并发上传，单结点并发量大于100 qps
5. 文件下载并发量大于10 qps
5. 支持限定单结点的储存空间
6. 支持快速扩容(加机器+更新配置即可实现扩容）

#### 实现方案
1. 数据文件按天分目录储存,方便实现按天删除
2. 文件的路径和偏移量信息记录在文件id中，方便快速定位文件，文件id格式(40字符)
   字段项 |version|结点编号|分组编号|年|月|日|文件块编号|文件在文件块中的偏移量
   --|--|--|--|--|--|--|--
   长度   |2      |  2     |2       |2 |1 |1 |2         |8
   - 注：
     - 支持最大结点数:65535
     - 支持最大分组数:65535
     - 单个目录下支持最大文件分块数:65535
     - 支持最大的文件：4G
     - 支持最大的分块:逻辑上无限制
3. 文件块的储存格式：
   字段项 |magic number|version|文件总长度|文件名长度|文件名   |文件内容|校验位|....
   --|--|--|--|--|--|--
   长度   |8           |2      |4         |     2    |[0,65535]| [0,4G] |4    |....
   
4. 储存到具体结点时，优先储存到当前结点，若当前结点的磁盘限额已用完，则随机选择一台还有额度的结点进行储存
5. 机器的结点编号一旦设定不能修改，因为结点编号与数据绑定
6. 所有结点都是同质的，提供完全相同的功能
7. 每个结点定期同步其它所有结点的文件统计信息

#### 对外接口
##### http接口
- 添加文件
  - url: /file/upload
  - param:
    - group:组号
    - timestamp:时间戳
    - filename:文件名
  - method: POST
  - body: 文件数据
  - response
    ```
    {
       "code":0/失败编号,
       "message":succes/失败原因,
       "group":组号,
       "filename":文件名,
       "id":文件id号
    }
     ```
- 访问文件
  - url: /file/download
  - param: 
    - id:文件id
  - method: GET
  - response:图片文件(以提交时的文件名命名)
 
- 按天删除
  - url: /file/batchDelete
  - param:
    - group:分组号
    - day:日期(年月日)
  - method: DELETE
  - response
    ```
    {
    "code":0,
    "message":"success",
    }
    ``` 
- 统计数据 
  - url: /system/stat
  - method: GET
  - response
    ```
    {
    "code":0,
    "message":"success",
    "fileNum":1000000,
    "diskUsed":100000000,
    "nodes":[
        {
            "id":0,
            "alive":true/false,
            "status":0,
            "host":192.168.1.223,
            "port":8090,
            "fileNum":10000,
            "diskUsed":100000,
            "diskLimited":200000000
        }
    ],
    "groups":
        [
            {
            "id": 1,
            "fileNum":100000,
            "diskUsed":10000000,
            "nodes":[
                {
                    "id":0,
                    "fileNum":10000
                    "diskUsed":100000,
                    "diskAvailable":200000000,
                    "partitions":[
                        {
                            "date":"20190304"
                            "limitedBlockNum":100,
                            "fileNum":10000,
                            "diskUsed":123234,
                            "blocks":[
                                {
                                    "id":1,
                                    "fileNum":1232324,
                                    "diskUsed":123234,
                                }
                            ]
                            
                        }
                    ]
                }
                
            ]
            }
        ]
    }
    ``` 

#### 配置
```
[server]
ip:服务绑定的ip
port:服务绑定的端口
peers: ip1:port1:id1,ip2:port2:id2,ip3:port3:id3,ip4:port4:id4
data_path:数据目录
sync_period:从其它结点同步状态数据的周期(秒)

[log]
path:日志路径
level:日志级别
```
#### 内部rpc接口
- localStore(groupId int, fileName string, fileData []byte) (error, string)
  - 功能:储存文件到本地
  - 参数
    - groupId: 分组id
    - fileName:文件名
    - fileData:文件数据 
  - 返回
    - 出错信息
    - 文件id
- heartBeat() int 
  - 功能:心跳接口,返回当前结点的模式,可读写，只读，不可读不可写
  - 参数:无
  - 返回
    - 当前结点的模式:0-不可读且不可写, 1-仅可读, 2-仅可写, 3-可读可写
- localRetrieve(groupId int, partition string, blockId int, offset int64) (error, string, []byte)
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
- localStat() (error, StatInfo)
  - 功能: 统计本地的状态信息
  - 参数:无
  - 返回
    - 出错信息
    - 状态信息数据
