package imgcode

import (
	"bytes"
	"encoding/binary"
	"errors"
	"image"
	"image/color"
	"image/png"
	"io"
	"os"
)

var point uint32 = 0 //供SeparateData读取的全局变量，用于当前数组的下标储存

func IDII(x int, y int, imgName string, data []byte) error { //Integrate data into images
	if x <= 0 && y <= 0 {
		return errors.New("x,y必须大于0")
	}
	if len(imgName) <= 0 {
		return errors.New("图片文件地址必须指定")
	}
	if x*y*3 < len(data) {
		return errors.New("图片大小不够")
	}

	bs := make([]byte, 6) //每个像素储存3个直接，而uint32有4个字节，所以申请6个字节，存放到2个像素内
	binary.BigEndian.PutUint32(bs, uint32(len(data)))
	//fmt.Println("BS:", bs)
	point = 0
	dataTemp := append(bs, data...) //把data长度追加到data数组的前面，存到临时变量内
	data = dataTemp
	rect := image.Rect(0, 0, x, y)
	img := image.NewRGBA(rect)

	for j := 0; j < y; j++ {
		for i := 0; i < x; i++ {

			pixData := SeparateData(data)
			/** 坐标轴
				-----------------------------> X (i)
				|(0,0) (1,0) (2,0) (3,0) (4,0)
				|(0,1) (1,1) (2,1) (3,1) (4,1)
				|(0,2) (1,2) (2,2) (3,2) (4,2)
				|(0,3) (1,3) (2,3) (3,3) (4,3)
				|(0,4) (1,4) (2,4) (3,4) (4,4)
				V
				  Y (j)
			**/
			img.SetRGBA(i, j, color.RGBA{pixData[0], pixData[1], pixData[2], 255})

		}

	}
	//fmt.Println(point)
	file, err := os.Create(imgName)
	defer file.Close()
	if err != nil {
		return err
	}
	err = png.Encode(file, img)
	if err != nil {
		return err
	}
	return nil
}

func DPTD(reader io.Reader, filename string) error { //Decode pixels to data
	img, err := png.Decode(reader)
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	if err != nil {
		return err
	}
	r, g, b, _ := img.At(0, 0).RGBA()
	r2, g2, b2, _ := img.At(1, 0).RGBA()
	//下面这里是直接读出最前面的2个像素的颜色值，然后输出成uint32数据类型的文件大小到size变量
	size := binary.BigEndian.Uint32([]byte{uint8(r),
		uint8(g),
		uint8(b),
		uint8(r2),
		uint8(g2),
		uint8(b2)})
	if size <= 0 {
		return errors.New("文件大小错误")
	}
	var w *os.File

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		//文件不存在
		w, err = os.Create(filename)
		if err != nil {
			return err
		}
	} else {
		return errors.New("文件已经存在！")
	}

	w.Truncate(int64(size))
	defer w.Close()
	var FilePointer uint32 = 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if x == 0 && y == 0 {
				//nothing
				continue
			}
			if x == 1 && y == 0 {
				//nothing
				continue
			}
			dr, dg, db, _ := img.At(x, y).RGBA()
			byteData := []byte{uint8(dr), uint8(dg), uint8(db)}
			if FilePointer < size {
				if (FilePointer + 3) > size {
					switch size - FilePointer {
					case 2:
						byteData = byteData[:2]
					case 1:
						byteData = byteData[:1]
					}
				}
				_, err = w.WriteAt(byteData, int64(FilePointer))
				if err != nil {
					return err
				}
				FilePointer += 3
			} else {

				break
			}

		}
	}

	return nil
}

func SeparateData(data []byte) []uint8 {
	dataLen := uint32(len(data))
	if dataLen > (point + 3) {
		var data1, data2, data3 uint8
		/**
			这里主要是把data数组以uint8类型读出来，保存到对应的data1~3中
		**/
		b_buf := bytes.NewBuffer([]byte{data[point]})
		binary.Read(b_buf, binary.BigEndian, &data1)
		b_buf = bytes.NewBuffer([]byte{data[point+1]})
		binary.Read(b_buf, binary.BigEndian, &data2)
		b_buf = bytes.NewBuffer([]byte{data[point+2]})
		binary.Read(b_buf, binary.BigEndian, &data3)

		point += 3
		return []uint8{data1, data2, data3}
	} else if dataLen > point {

		var retData []uint8
		for i := point; i < dataLen; i++ {
			var data1 uint8
			b_buf := bytes.NewBuffer([]byte{data[i]})
			binary.Read(b_buf, binary.BigEndian, &data1)
			retData = append(retData, data1)

		}
		newRetData := make([]uint8, 3-len(retData))
		retData = append(retData, newRetData...)
		point = dataLen
		return retData
	} else {
		return []uint8{0x0, 0x0, 0x0}
	}
	return []uint8{0x0, 0x0, 0x0}
}
