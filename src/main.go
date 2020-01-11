package main

import (
    . "trival/server"
    . "trival/service"
    . "trival/utils"
    "os"
//    "log"
    "flag"
    "fmt"
    "os/signal"
)

const(
   PROGRAM_VERSION = "v0.0.1"
   DATABASE_VERSION = "1"
)

var signals chan os.Signal = make(chan os.Signal, 1)
func main(){
    var version bool 
    flag.BoolVar(&version, "v", false, "show version")
    flag.StringVar(&ConfPath, "c", "conf/conf.toml", "configure file path")
    flag.Parse()
    if version{
        fmt.Printf("program version:%s, database version:%s\n", 
                        PROGRAM_VERSION, DATABASE_VERSION);
        return
    }
	InitLog(Config().LogPath)
    go ServeRPC()
    go ServeHTTP()
	signal.Notify(signals, os.Interrupt)
    <-signals
}
