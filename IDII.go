package main

import (
	"fmt"
	"imgcode"
	"math"
	"os"
)

func main() {
	f, err := os.Open("test.exe") //要编码的文件
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	fileData := make([]byte, 1024*1024*10) //1Kb=1024byte,1Mb=1024Kb 默认最大读取10M的文件进行编码
	count, err := f.Read(fileData)
	xy := int(math.Sqrt(float64(count/3))) + 1
	err = imgcode.IDII(xy, xy, "test.png", fileData[:count]) //要输出的png图片
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("成功!")
}
