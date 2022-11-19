package main

import (
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"github.com/bjdgyc/slog"
)

var (
	connPool *ConnPool
)

func main() {

	cfile := flag.String("c", "./config.toml", "配置文件")
	flag.Parse()

	config := LoadConfig(*cfile)

	if config.Logfile != "" {
		slog.SetLogfile(config.Logfile)
	}

	for _, opt := range config.Options {
		connPool = NewConnPool(opt)
		go StartProxy(connPool, opt.Addr)
		fmt.Println(connPool)
	}

	go func() {
		http.ListenAndServe("127.0.0.1:8090", nil)
	}()

	//connPool.Close()
	InitSignal()

}
