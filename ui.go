package main

import (
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

func init() {
	//设置中文字体:解决中文乱码问题
	fontPaths := findfont.List()
	for _, path := range fontPaths {
		if strings.Contains(path, "STKAITI.TTF") || strings.Contains(path, "msyh.ttf") || strings.Contains(path, "simhei.ttf") || strings.Contains(path, "simsun.ttc") || strings.Contains(path, "simkai.ttf") {
			os.Setenv("FYNE_FONT", path)
			os.Setenv("FYNE_FONT_MONOSPACE", path)
			break
		}
	}
}

// -ldflags "-s -w -H=windowsgui"
func startUI() {
	APP = app.New()
	APP.Settings().SetTheme(theme.LightTheme())
	WINDOW = APP.NewWindow("GDUT WIFI")
	WINDOW.Resize(fyne.NewSize(600, 400))
	whiteRectangle := canvas.NewRectangle(color.White)
	_ = getUsernamePassword()
	usernameEntry := widget.NewEntry()
	passwordEntry := widget.NewPasswordEntry()
	if USERNAME != "" && PASSWORD != "" {
		usernameEntry.SetText(USERNAME)
		passwordEntry.SetText(PASSWORD)
	} else {
		passwordEntry.SetPlaceHolder("Enter password...")
		usernameEntry.SetPlaceHolder("Enter username...")
	}
	text := canvas.NewText("current status:                         ", color.Black)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}
	userInformation := container.NewGridWithColumns(1, whiteRectangle, usernameEntry, passwordEntry, widget.NewButton("Save", func() {
		save(usernameEntry.Text, passwordEntry.Text)
	}), whiteRectangle)
	connectStatus := canvas.NewCircle(isConnected())
	go func() {
		for _ = range time.Tick(time.Second * 2) {
			connectStatus.FillColor = isConnected()
			connectStatus.Refresh()
		}
	}()
	txtResults := widget.NewTextGrid()
	txtResults.ShowLineNumbers = true
	txtResults.SetText("中文connect log ...")
	cntScrolling := container.NewScroll(txtResults)
	// 创建一个较小的网格布局
	connectStatusLayout := container.New(layout.NewGridLayoutWithColumns(1), whiteRectangle, whiteRectangle, whiteRectangle, whiteRectangle, connectStatus, whiteRectangle, whiteRectangle, whiteRectangle, whiteRectangle)
	connectButton := widget.NewButton("connect To Gdut", func() {
		doConnect(txtResults, cntScrolling)
	})
	buttonLayout := container.New(layout.NewGridLayout(1), connectButton)
	// 将连接状态的圆形放置在网格布局中
	userStatus1 := container.New(layout.NewGridLayout(2), text, connectStatusLayout)
	userStatus2 := container.New(layout.NewGridLayout(2), buttonLayout, cntScrolling)
	userStatus := container.New(layout.NewGridLayout(1), userStatus1, userStatus2)
	hyperlink := widget.NewHyperlink("开发不易,求给作者一个star", nil)
	hyperlink.SetURLFromString(`https://gitee.com/qiu_168/campus-network-public`)
	statusHelp := widget.NewLabel("tips:\n\n" +
		"1.`current status` 是绿色表示当前可联网,红色表示当前不可联网\n" +
		"2.`username,password` 是校园网账号和密码\n" +
		"3.连接超时的情况下 1.查看自己有没有打开wifi，2.重试\n" +
		"4.在连接网络的日志显示模块中，对中文的适配性不高\n" +
		"5.如果当天已连接过，再连接日志会显示连接错误，但没有影响\n")
	tabs := container.NewAppTabs(
		container.NewTabItem("connect wifi", userStatus),
		container.NewTabItem("user information", userInformation),
		container.NewTabItem("help", container.NewGridWithColumns(1, hyperlink, statusHelp)),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	WINDOW.SetContent(tabs)
	WINDOW.SetMaster()
	WINDOW.ShowAndRun()
}

func doConnect(txtResults *widget.TextGrid, cntScrolling *container.Scroll) {
	go DoConnectGdut()
	txtResults.SetText("")
	cntScrolling.Refresh()
	go func() {
		for s := range CHANNEL {
			txtResults.SetText(strings.TrimPrefix(txtResults.Text()+"\n"+s, "\n"))
			cntScrolling.Refresh()
			cntScrolling.ScrollToBottom()
		}
	}()
}

func isConnected() color.Color {
	if isWiFiLoggedIn() {
		return color.RGBA{R: 0, G: 255, B: 0, A: 255}
	}
	return color.RGBA{R: 255, G: 0, B: 0, A: 255}
}

func save(usernameText string, passwordText string) {
	err := createFileIfNotExistAndWrite("./config.txt", usernameText, passwordText)
	if err != nil {
		dialog.NewInformation("error", "文件创建或写入失败", WINDOW).Show()
		return
	}
	USERNAME = usernameText
	PASSWORD = passwordText
	text := canvas.NewText("Save successfully", color.Black)
	var notice *widget.PopUp
	button := widget.NewButton("close", func() {
		notice.Hide()
	})
	notice = widget.NewModalPopUp(container.New(layout.NewGridLayout(1), text, button), WINDOW.Canvas())
	notice.Show()
}
