package main

import (
	"net"
	"net/http"
	"time"
)

func isWiFiLoggedIn() bool {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	resp, err := client.Get("https://www.baidu.com")
	if err != nil {
		return false // 无法访问，可能未登录
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true // 访问成功，认为已登录
	}

	return false
}
func isNetworkConnected() bool {
	host := "www.baidu.com"
	port := "443"
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
