package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

func main() {
	ConnectGDUT()
	startTime := time.Now()
	timeout := 2 * time.Second
	for {
		ssid, _ := getCurrentWiFiSSID()
		if ssid == "gdut" {
			break
		}
		duration := time.Since(startTime) / time.Millisecond
		fmt.Printf("running in %dms\n", duration)
		if duration > timeout {
			fmt.Println("程序超时")
			panic("连接gdut失败")
		}
	}
	//fmt.Scanf("h")
	username, password := getConfig()
	fmt.Println(username)
	fmt.Println(len(username))
	fmt.Println(password)
	fmt.Println(len(password))
	cookie := make(map[string]string)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 禁用自动重定向
			return http.ErrUseLastResponse
		},
	}
	//生成要访问的url
	url := "http://www.msftconnecttest.com/redirect"

	//提交请求
	request, _ := http.NewRequest("GET", url, nil)

	//增加header选项
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	request.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	request.Header.Add("Host", "www.msftconnecttest.com")
	request.Header.Add("Pragma", "no-cache")
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Cache-Control", "no-cache")
	request.Header.Add("Upgrade-Insecure-Requests", "1")
	//fmt.Println(request.Header)
	//fmt.Println("中文")
	//处理返回结果
	resp, _ := client.Do(request)
	location := resp.Header.Get("Location")
	// 正则表达式来匹配IP地址
	fmt.Println(location)
	ipRegex := `wlanuserip=([\d.]+)&wlanacname=&wlanacip=([\d.]+)`
	re := regexp.MustCompile(ipRegex)
	matches := re.FindStringSubmatch(location)
	cookie["wlanuserip"] = matches[1]
	cookie["wlanacip"] = matches[2]
	h := "http://10.0.3.2:801/eportal/portal/login?" +
		"callback=dr1003&" +
		"login_method=1&" +
		"user_account=%2C0%2C" + username + "&" +
		"user_password=" + password + "&" +
		"wlan_user_ip=" + cookie["wlanuserip"] + "&" +
		"wlan_user_ipv6=" +
		"&wlan_user_mac=000000000000&" +
		"wlan_ac_ip=" + cookie["wlanacip"] + "&" +
		"wlan_ac_name=&" +
		"jsVersion=4.1.3&" +
		"terminal_type=1&" +
		"lang=zh-cn&" +
		"v=8649&" +
		"lang=zh"
	newRequest, _ := http.NewRequest("GET", h, nil)
	newRequest.Header.Add("Accept", "*/*")
	newRequest.Header.Add("Accept-Encoding", "gzip, deflate")
	newRequest.Header.Add("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	newRequest.Header.Add("Cache-Control", "no-cache")
	newRequest.Header.Add("Connection", "keep-alive")
	newRequest.Header.Add("Host", "10.0.3.2:801")
	newRequest.Header.Add("Pragma", "no-cache")
	newRequest.Header.Add("Referer", "http://10.0.3.2/")
	newRequest.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.60")
	do, _ := client.Do(newRequest)
	defer do.Body.Close()
	all, _ := io.ReadAll(do.Body)
	fmt.Println(string(all))
	defer resp.Body.Close()
	if isWiFiLoggedIn() {
		fmt.Println("network connected !!!! Ciallo～(∠・ω< )⌒☆")
	} else {
		fmt.Println("network error qwq")
	}
}
