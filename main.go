package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func main() {
	filePath := "C:\\Windows\\SysWOW64\\Macromed\\Flash\\Flash.ocx"
	fmt.Println("读取文件：", filePath)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件错误，请检查文件是否存在")
		return
	}
	targetBytes := []byte{0x00, 0x00, 0x40, 0x46, 0x3E, 0x6F, 0x77, 0x42}
	fmt.Println("查找目标代码...")
	if bytes.Contains(content, targetBytes) {
		fmt.Println("找到了目标代码，尝试修改...")
		fmt.Println("修改文件权限...")
		c := exec.Command("cmd", "/C", "TAKEOWN", "/F", filePath, "/a")
		if err := c.Run(); err != nil {
			fmt.Println("修改文件权限失败，请用管理员权限运行")
			return
		}
		c2 := exec.Command("cmd", "/C", "ICACLS", filePath, "/grant", "Administrators:F")
		if err := c2.Run(); err != nil {
			fmt.Println("修改文件权限失败，请用管理员权限运行")
			return
		}
		targetIndex := bytes.Index(content, targetBytes)
		modifyIndex := targetIndex + 7
		content[modifyIndex] = 0x72
		err := ioutil.WriteFile(filePath, content, os.ModePerm)
		if err == nil {
			fmt.Println("修改成功，flash已被激活")
		} else {
			fmt.Println("修改失败，原因：")
			fmt.Println(err)
		}
	} else {
		fmt.Println("找不到目标代码，修改失败")
	}

	fmt.Print("\n按enter键退出...")
	fmt.Scanln()
}
