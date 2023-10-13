package main

import (
	"net"
	"time"
)

func isNetworkConnected() bool {
	host := "www.baidu.com"
	port := "80"
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
