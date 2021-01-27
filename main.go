package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

/**
 * 此代码只针对windows系统
 */
func main() {
	// 获取操作系统位数
	bit := 32 << (^uint(0) >> 63)
	var filePath string
	if bit == 64 {
		filePath = "C:\\Windows\\SysWOW64\\Macromed\\Flash\\Flash.ocx"
	} else {
		filePath = "C:\\Windows\\System32\\Macromed\\Flash\\Flash.ocx"
	}
	fmt.Println("Flash On v0.3")
	fmt.Println("！！！运行前请关闭浏览器！！！")
	fmt.Printf("%d位操作系统\n", bit)
	fmt.Printf("读取文件：%s\n", filePath)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("读取文件失败，请检查文件是否存在")
		return
	}
	targetBytes := []byte{0x00, 0x00, 0x40, 0x46, 0x3E, 0x6F, 0x77}
	fmt.Println("开始查找目标代码...")
	if bytes.Contains(content, targetBytes) {
		targetIndex := bytes.Index(content, targetBytes)
		modifyIndex := targetIndex + 7
		if content[modifyIndex] == 0x72 {
			fmt.Println("已被激活的Flash，不需要操作")
		} else {
			fmt.Println("查找到了目标代码")
			fmt.Println("开始尝试修改...")
			fmt.Println("获取文件权限...")
			c := exec.Command("cmd", "/C", "TAKEOWN", "/F", filePath, "/a")
			if err := c.Run(); err != nil {
				fmt.Println("修改文件权限失败，请用管理员权限运行！")
				return
			}
			c2 := exec.Command("cmd", "/C", "ICACLS", filePath, "/grant", "Administrators:F")
			if err := c2.Run(); err != nil {
				fmt.Println("修改文件权限失败，请用管理员权限运行！")
				return
			}
			content[modifyIndex] = 0x72
			err := ioutil.WriteFile(filePath, content, os.ModePerm)
			if err == nil {
				fmt.Println("修改成功，flash已被激活")
			} else {
				fmt.Println("修改失败，原因：")
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println("找不到目标代码，修改失败")
	}
	fmt.Print("\n按enter键退出...")
	fmt.Scanln()
}
