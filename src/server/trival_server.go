package server

import (
    "net/http"
    "time"
    . "trival/utils"
    "fmt"
    "log"
)



type Server struct {
}

func NewServer() *Server {
    return &Server{
	}
}

func (this *Server) Upload(w http.ResponseWriter, r *http.Request){

}

func (this *Server) Download(w http.ResponseWriter, r *http.Request){

}

func (this *Server) BatchDelete(w http.ResponseWriter, r *http.Request){

}
func (this *Server) Stat(w http.ResponseWriter, r *http.Request){

}

func ServeHTTP(){
    server := NewServer()
    http.HandleFunc("/file/upload", server.Upload)
    http.HandleFunc("/file/download", server.Download)
    http.HandleFunc("/file/batchDelete", server.BatchDelete)
    http.HandleFunc("/system/stat", server.Stat)

    addr := fmt.Sprintf("%s:%d", Config().Http.IP, Config().Http.Port)
	srv := &http.Server{
        Addr: addr,
		ReadTimeout:time.Duration(Config().Http.ReadTimeout) * time.Second,
		WriteTimeout:time.Duration(Config().Http.WriteTimeout) * time.Second,
	}
	err := srv.ListenAndServe()
    if err != nil {
        log.Fatal("start http service failed",err)
    }else{
        log.Printf("start http service success, listen on:%s", addr)
    }
}
