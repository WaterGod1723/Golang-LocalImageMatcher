package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"imgSearcher/imgHandle"
	"imgSearcher/screenshot"
	"log"
	"os"
)

func main() {
	fmt.Println("Starting...")
	screenshot.WatchClipboard(func(img *image.RGBA) {
		fmt.Println("Image copied to clipboard")
		// 图片边缘检测，输出检测后的图片
		pngData := new(bytes.Buffer)
		resImg, _ := imgHandle.DetectEdges(img)
		png.Encode(pngData, resImg)

		// 将处理后的数据保存为png文件
		err := os.WriteFile("./aaaaaaaaaaaa.png", pngData.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	})
}
