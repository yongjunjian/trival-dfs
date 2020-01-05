package types
//配置相关
type ConfigInfo struct{
    DataPath string `toml:"data_path"`
    LogPath string  `toml:"log_path"`
    Http HttpInfo
    Rpc RpcInfo
    Storage StorageInfo
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
    MaxBlockNum int `toml:"max_block_num"`
    MaxBlockSize int `toml:"max_block_size"`
    MaxDiskUsage int `toml:"max_disk_usage"`
}


