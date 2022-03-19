package main

import (
	"bufio"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	deviceMap := make(map[string]string)
	var displayDeviceFunc = func() {
		deviceMap = make(map[string]string)
		cmd := exec.Command("powercfg", "/devicequery", "wake_armed")
		out, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		}
		out, _ = simplifiedchinese.GBK.NewDecoder().Bytes(out)
		if strings.HasPrefix(string(out), "无") {
			fmt.Println("没有能唤醒休眠的设备")
			return
		}
		fmt.Println("请选择要禁用的设备序号：")
		deviceArr := strings.Split(string(out), "\n")
		for i, device := range deviceArr {
			if len(strings.TrimSpace(device)) > 0 {
				index := strconv.Itoa(i + 1)
				fmt.Println(index + ":" + device)
				deviceMap[index] = strings.Replace(device, "\r", "", -1)
			}
		}
	}
	displayDeviceFunc()
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		line := input.Text()
		if deviceMap[line] != "" {
			cmd := exec.Command("powercfg", "/devicedisablewake", deviceMap[line])
			stderr, _ := cmd.StderrPipe()
			err := cmd.Start()
			var failFunc = func(errMsg string) {
				fmt.Println("禁用失败：" + deviceMap[line])
				fmt.Println("原因：" + errMsg)
				displayDeviceFunc()
			}
			if err != nil {
				failFunc(err.Error())
				continue
			}
			slurp, _ := io.ReadAll(stderr)
			b, _ := simplifiedchinese.GBK.NewDecoder().Bytes(slurp)
			stderrMsg := string(b)
			if len(stderrMsg) > 0 {
				failFunc(stderrMsg)
				continue
			}
			fmt.Println("已禁用该设备唤醒功能：" + deviceMap[line])
			deviceMap[line] = ""
		} else {
			fmt.Println("选项不存在，请重新输入序号")
		}
		displayDeviceFunc()
	}

}
