package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var (
	CHANNEL  = make(chan string, 30)
	PASSWORD string
	USERNAME string
	APP      fyne.App
	WINDOW   fyne.Window
)

func main() {
	//DoConnectGdut()
	startUI()
}

func DoConnectGdut() {
	defer func() {
		if err := recover(); err != nil {
			CHANNEL <- "连接发生错误"
			fmt.Printf("连接发生错误 %v", err)
		}
	}()
	ConnectGDUT()
	waitForConnected()
	if getUsernamePassword() != nil {
		return
	}
	CHANNEL <- "username:" + USERNAME
	CHANNEL <- "len(username):" + strconv.Itoa(len(USERNAME))
	fmt.Println("username:" + USERNAME)
	fmt.Println("len(username):" + strconv.Itoa(len(USERNAME)))
	client := newClient()
	userIp, acIp := sendPreRequest(client)
	sendConnectRequest(userIp, acIp, client)
	if isWiFiLoggedIn() {
		CHANNEL <- "network connected !!!! Ciallo～(∠・ω< )⌒☆"
		fmt.Println("network connected !!!! Ciallo～(∠・ω< )⌒☆")
	} else {
		CHANNEL <- "network error qwq"
		fmt.Println("network error qwq")
	}
}

func getUsernamePassword() error {
	var err error
	if USERNAME == "" || PASSWORD == "" {
		USERNAME, PASSWORD, err = getConfig()
		if err != nil {
			dialog.NewInformation("error", "获取username，password失败\n请先配置", WINDOW).Show()
			return err
		}
		return nil
	}
	return nil
}
func newClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 禁用自动重定向
			return http.ErrUseLastResponse
		},
	}
}

func sendPreRequest(client *http.Client) (string, string) {
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

	//处理返回结果
	resp, _ := client.Do(request)
	defer resp.Body.Close()
	location := resp.Header.Get("Location")
	// 正则表达式来匹配IP地址
	ipRegex := `wlanuserip=([\d.]+)&wlanacname=&wlanacip=([\d.]+)`
	re := regexp.MustCompile(ipRegex)
	matches := re.FindStringSubmatch(location)
	CHANNEL <- "finish preRequest"
	return matches[1], matches[2]
}

func sendConnectRequest(userIp, acIp string, client *http.Client) {
	h := "http://10.0.3.2:801/eportal/portal/login?" +
		"callback=dr1003&" +
		"login_method=1&" +
		"user_account=%2C0%2C" + USERNAME + "&" +
		"user_password=" + PASSWORD + "&" +
		"wlan_user_ip=" + userIp + "&" +
		"wlan_user_ipv6=" +
		"&wlan_user_mac=000000000000&" +
		"wlan_ac_ip=" + acIp + "&" +
		"wlan_ac_name=&" +
		"jsVersion=4.1.3&" +
		"terminal_type=1&" +
		"lang=zh-cn&" +
		"v=8649&" +
		"lang=zh"
	newRequest, _ := http.NewRequest("GET", h, nil)
	newRequest.Header.Add("Accept", "*/*")
	newRequest.Header.Add("Accept-Encoding", "")
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
	CHANNEL <- string(all)
}

func waitForConnected() {
	startTime := time.Now()
	timeout := 3 * time.Second
	for {
		ssid, _ := getCurrentWiFiSSID()
		if ssid == "gdut" {
			break
		}
		duration := time.Since(startTime)
		info := fmt.Sprintf("running in %s", duration.String())
		fmt.Println(info)
		CHANNEL <- info
		if duration > timeout {
			fmt.Println("程序超时")
			CHANNEL <- "程序超时"
			panic("连接gdut失败")
		}
	}
}
