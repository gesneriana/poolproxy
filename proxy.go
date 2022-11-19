package main

import (
	"net"
	"strings"

	"github.com/bjdgyc/slog"
)

// 监听连接
func StartProxy(connPool *ConnPool, addr string) {
	n := "unix"
	if strings.Contains(addr, ":") {
		n = "tcp"
	}

	l, err := net.Listen(n, addr)
	if err != nil {
		slog.Fatal("listen error:", err)
		return
	}
	defer l.Close()

	for {
		local, err := l.Accept()
		if err != nil {
			slog.Warn("accept error:", err)
			continue
		}

		go HandlerData(connPool, local)
	}

}

// 数据交换方法
func HandlerData(connPool *ConnPool, local net.Conn) {

	slog.Debug("remote addr:", local.RemoteAddr())

	conn, err := connPool.Get()
	if err != nil {
		local.Close()
		slog.Error("pool get error:", err)
		return
	}

	forceClose := conn.SwapData(local)
	local.Close()
	connPool.Put(conn, forceClose)

}
