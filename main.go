package main

import (
	"bytes"
	"fmt"
	"github.com/mattn/go-colorable"
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
	str :=
		` ---------------
| Flash On v0.3 |
 ---------------

`
	fmt.Print(str)

	fmt.Printf("！！！运行前请关闭浏览器！！！\n\n")
	fmt.Printf("%d位操作系统\n", bit)
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		PrintRed("\n读取文件失败，请检查文件是否存在\n")
		return
	}
	targetBytes := []byte{0x00, 0x00, 0x40, 0x46, 0x3E, 0x6F, 0x77}
	if bytes.Contains(content, targetBytes) {
		targetIndex := bytes.Index(content, targetBytes)
		modifyIndex := targetIndex + 7
		if content[modifyIndex] == 0x72 {
			PrintGreen("\n已被激活的Flash，不需要操作\n")
		} else {
			fmt.Println("开始尝试激活...")
			c := exec.Command("cmd", "/C", "TAKEOWN", "/F", filePath, "/a")
			if err := c.Run(); err != nil {
				PrintRed("获取文件权限失败，请用管理员权限运行！\n")
				return
			}
			c2 := exec.Command("cmd", "/C", "ICACLS", filePath, "/grant", "Administrators:F")
			if err := c2.Run(); err != nil {
				PrintRed("获取文件权限失败，请用管理员权限运行！\n")
				return
			}
			content[modifyIndex] = 0x72
			err := ioutil.WriteFile(filePath, content, os.ModePerm)
			if err == nil {
				PrintGreen("修改成功，flash已被激活\n")
			} else {
				PrintRed("修改失败，原因：\n")
				fmt.Println(err)
			}
		}
	} else {
		PrintRed("找不到目标代码，修改失败")
	}
	fmt.Print("\n按enter键退出...")
	fmt.Scanln()
}

func PrintRed(str string) {
	PrintColor(str, "31")
}

func PrintGreen(str string) {
	PrintColor(str, "32")
}

func PrintColor(str string, color string) {
	stdout := colorable.NewColorable(os.Stdout)
	_, _ = fmt.Fprintf(stdout, "\x1b["+color+"m%s\x1b[0m", str)
}
