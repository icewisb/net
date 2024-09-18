package main

import (
        "fmt"
        "net"
        "net/http"
	"io/ioutil"

        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/app"
        "fyne.io/fyne/v2/container"
        "fyne.io/fyne/v2/widget"
)



// 获取公网 IP
func getIP() string {
    // 发送 GET 请求
    response, err := http.Get("https://ident.me")
    if err != nil {
        return ""
    }
    defer response.Body.Close() // 确保在函数结束时关闭响应体

    // 读取响应体
    body, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return ""
    }

    // 输出响应文本
    return string(body)
}

// 获取 MAC 地址和 IP 地址
func getMAC() (map[string]string, error) {
        dct := make(map[string]string)

        interfaces, err := net.Interfaces()
        if err != nil {
                return nil, err
        }

        for _, iface := range interfaces {
                if iface.Flags&net.FlagLoopback != 0 {
                        continue // 跳过回环接口
                }

                addrs, err := iface.Addrs()
                if err != nil {
                        return nil, err
                }

                for _, addr := range addrs {
                        if ipNet, ok := addr.(*net.IPNet); ok {
                                if ipNet.IP.To4() != nil { // IPv4 地址
                                        dct[iface.Name] = fmt.Sprintf("MAC: %s, IP: %s", iface.HardwareAddr.String(), ipNet.IP.String())
                                }
                        }
                }
        }
        return dct, nil
}

func main() {

        // 获取 IP 和 MAC 地址
        ip := getIP()

        macInfo, err := getMAC()
        if err != nil {
                macInfo = map[string]string{"Error": "Error getting MAC addresses"}
        }

        // 创建显示信息的文本
        infoText := "出口 IP: " + ip + "\n\nMAC Addresses:\n"
        for iface, info := range macInfo {
                infoText += fmt.Sprintf("%s: %s\n", iface, info)
        }

        myApp := app.New()
        myWindow := myApp.NewWindow("Network Info")
        // 创建文本组件
        // text := widget.NewLabel(infoText)
	textArea := widget.NewMultiLineEntry()
        textArea.SetText(infoText)

        // 创建窗口内容
        myWindow.SetContent(container.NewVBox(textArea))
	myWindow.Resize(fyne.NewSize(400, 300))
        myWindow.ShowAndRun()
}
