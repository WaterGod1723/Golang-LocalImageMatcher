package imgHandle_test

import (
	"bytes"
	"fmt"
	"image/png"
	"imgSearcher/imgHandle"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func TestCrop(t *testing.T) {
	// 切换工作目录
	err := os.Chdir("D:/golang/Golang-LocalImageMatcher")
	if err != nil {
		log.Fatal(err)
	}
	// 指定要遍历的目录
	dir := "./textImgs"
	// 清空output文件夹
	err = os.RemoveAll("./textImgs/output")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("清空output文件夹")
	imageFiles, err := imgHandle.GetAllImageFiles(dir)
	if err != nil {
		log.Fatal(err)
	}

	// 遍历所有图片文件
	for _, file := range imageFiles {
		// 读取图片文件
		data, err := os.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}

		// 图片边缘检测，输出检测后的图片
		pngData := new(bytes.Buffer)
		resImg, _ := imgHandle.DetectEdges(data)
		png.Encode(pngData, resImg)

		// 创建output文件夹
		err = os.MkdirAll("./textImgs/output", 0755)
		if err != nil {
			log.Fatal(err)
		}

		// 将处理后的数据保存为png文件
		err = os.WriteFile("./textImgs/output/"+filepath.Base(file), pngData.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
