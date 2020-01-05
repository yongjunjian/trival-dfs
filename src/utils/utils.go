package utils
import (
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
    "github.com/BurntSushi/toml"
	"os"
    "flag"
    . "trival/types"
)
//exported
var (
   ConfPath string
)

//日志相关
type MultiWriter struct {
	filew io.Writer
	stdw  io.Writer
}

func (mw *MultiWriter) Write(p []byte) (n int, err error) {
	n1, err := mw.filew.Write(p)
	if err != nil {
		return n1, err
	}
	n2, err := mw.stdw.Write(p)
	if err != nil {
		return n2, err
	}
	return n1, err
}

func InitLog(logPath string){
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.SetOutput(&MultiWriter{
		filew: &lumberjack.Logger{
			Filename: logPath,
		},
		stdw: os.Stderr})
}

var config *ConfigInfo
func Config() *ConfigInfo{
    if config == nil{
        config = &ConfigInfo{}
        if !flag.Parsed(){
            log.Fatalf("flag not parsed!") 
        }
        _, err := toml.DecodeFile(ConfPath, config)
        if err != nil{
            log.Fatalf("load conf failed:%v", err) 
        }
    }
    return config
}
