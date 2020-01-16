package types
//配置相关
type ConfigInfo struct{
	Log LogInfo
    Http HttpInfo
    Rpc RpcInfo
    Storage StorageInfo
}
type LogInfo struct{
	Level string `toml:"log_level"`	
    Path string  `toml:"log_path"`
}
type HttpInfo struct{
    IP string `toml:"ip"`
    Port int `toml:"port"`
    ReadTimeout int `toml:"read_timeout"`
    WriteTimeout int `toml:"write_timeout"`
}

type RpcInfo struct{
    IP string `toml:"ip"`
    Port int `toml:"port"`
    EtcdRegistry string `toml:"etcd_registry"`
    SyncPeriod int `toml:"sync_period"`
}

type StorageInfo struct{
	ThreadNum int `toml:"thread_num"`
    MaxBlockNum int `toml:"max_block_num"`
    MaxBlockSize int64 `toml:"max_block_size"`
    MaxDiskUsage int `toml:"max_disk_usage"`
    AdjustInterval uint `toml:"adjust_thread_interval"`
	DataPath string `toml:"data_path"`
	DataNodeId string	`toml:"data_node_id"`
}


