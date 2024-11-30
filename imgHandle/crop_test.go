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
	err = os.MkdirAll("./textImgs/output", 0755)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("清空output文件夹")

	filepath.Walk(dir+"/input", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		fmt.Println(path)
		img, err := imgHandle.LoadImg(path)
		if err != nil {
			log.Fatal(err)
		}
		newImg := imgHandle.Crop(img, imgHandle.EdgeDetect(img))
		if err != nil {
			log.Fatal(err)
		}
		pngData := new(bytes.Buffer)
		png.Encode(pngData, *newImg)
		err = os.WriteFile("./textImgs/output/"+filepath.Base(path), pngData.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})
}

func TestHash(t *testing.T) {
	// l,_ := goimagehash.PerceptionHash()
	// l.
}
