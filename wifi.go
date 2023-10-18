package main

import (
	"fmt"
	_ "github.com/alexbrainman/odbc"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os/exec"
	"strings"
)

func ConnectGDUT() {
	// 打开 Wi-Fi
	enableWiFi := exec.Command("netsh", "interface", "set", "interface", "name=\"Wi-Fi\"", "admin=enable")
	enableWiFi.Run()

	// 列出可用的 Wi-Fi 网络
	listNetworks := exec.Command("netsh", "wlan", "show", "network")
	output, _ := listNetworks.Output()
	//fmt.Println(output)
	//fmt.Println(networks)
	gbkOutput, _, _ := transform.Bytes(simplifiedchinese.GBK.NewDecoder(), output)

	networks := string(gbkOutput)
	fmt.Println(networks)
	// 查找要连接的 Wi-Fi 网络名称
	desiredNetwork := "gdut" // 将 "YourSSID" 替换为你要连接的 Wi-Fi 网络名称

	// 查找网络名称对应的配置文件名
	var profileName string
	networkLines := strings.Split(networks, "\n")
	for _, line := range networkLines {
		if strings.Contains(line, "SSID") && strings.Contains(line, desiredNetwork) {
			profileName = strings.TrimSpace(strings.TrimPrefix(line, "SSID"))
			break
		}
	}

	if profileName == "" {
		println("未找到指定的 Wi-Fi 网络")
		return
	} else {
		profileName = "gdut"
		fmt.Println(profileName)
	}

	// 连接到指定的 Wi-Fi 网络
	connectToNetwork := exec.Command("netsh", "wlan", "connect", "name="+profileName)
	connectToNetwork.Run()
}
func getCurrentWiFiSSID() (string, error) {
	// 使用Windows的netsh命令获取当前连接的WiFi的SSID
	cmd := exec.Command("netsh", "wlan", "show", "interfaces")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	// 将输出数据转换为字符串
	outputStr := string(output)

	// 在输出中查找SSID信息
	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		if strings.Contains(line, "SSID") {
			ssid := strings.TrimSpace(strings.Split(line, ":")[1])
			return ssid, nil
		}
	}

	// 如果没有找到SSID信息
	return "", fmt.Errorf("无法获取当前WiFi的SSID")
}
