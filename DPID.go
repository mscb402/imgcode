package main

import (
	"fmt"
	"imgcode"
	"os"
)

func main() {
	f, err := os.Open("test.png") //要解码的图片
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	err = DPTD(f, "out.out") //解码后文件的保存地址
	if err != nil {
		fmt.Println(err)
		return
	}
}
