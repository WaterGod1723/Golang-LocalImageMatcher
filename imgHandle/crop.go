package imgHandle

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"

	"github.com/disintegration/imaging"
	"gonum.org/v1/gonum/mat"
)

var count = 0

func BytesImageDetect(data []byte) (*image.NRGBA, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return DetectEdges(img)
}

func DetectEdges(img image.Image) (*image.NRGBA, error) {
	matData, err := ImageToDense(img)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	// 利用库计算图像边缘
	edges := mat.NewDense(matData.RawMatrix().Rows, matData.RawMatrix().Cols, nil)
	edges.Apply(func(i, j int, v float64) float64 {
		// 使用 Sobel 算法计算边缘
		dx := 0.0
		dy := 0.0
		if j+1 < matData.RawMatrix().Cols {
			dx += matData.At(i, j+1)
		}
		if j-1 >= 0 {
			dx -= matData.At(i, j-1)
		}
		if i+1 < matData.RawMatrix().Rows {
			dy += matData.At(i+1, j)
		}
		if i-1 >= 0 {
			dy -= matData.At(i-1, j)
		}
		if math.Sqrt(dx*dx+dy*dy) > 100 {
			return 255
		} else {
			return 0
		}
	}, matData)
	// 从上下左右四个方向向边缘靠拢，裁剪掉多余部分
	maxY0, maxX0 := edges.Dims()
	maxX, maxY := 0, 0
	minX, minY := maxX0, maxY0
	MAX := 1.0
	for y := 0; y < maxY0; y++ {
		for x := 0; x < maxX0; x++ {
			if edges.At(y, x) > MAX {
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	// 保存边缘图像
	bounds := image.Rect(0, 0, maxX-minX, maxY-minY)
	edgeImage := image.NewNRGBA(bounds)
	for y := 0; y < bounds.Dy(); y++ {
		for x := 0; x < bounds.Dx(); x++ {
			value := edges.At(minY+y, minX+x)
			// 将 float64 值转换为 uint8 值
			grayValue := uint8(math.Min(255, value))
			edgeImage.Set(x, y, color.NRGBA{R: grayValue, G: grayValue, B: grayValue, A: 255})
		}
	}

	// file, _ := os.Create(fmt.Sprintf("./textImgs/edge_%d.png", count))
	// png.Encode(file, edgeImage)
	// defer file.Close()

	// 根据min，max对边缘进行裁剪
	// newImg := imaging.Crop(img, image.Rect(minX, minY, maxX+1, maxY+1))
	newImg := imaging.Resize(edgeImage, 64, 64, imaging.Lanczos)
	return newImg, nil
}

type ImageForDense interface {
	Bounds() image.Rectangle
	At(int, int) color.Color
}

func ImageToDense(img ImageForDense) (*mat.Dense, error) {
	var gray ImageForDense
	if s, ok := img.(*image.NRGBA); ok {
		gray = s
	} else if s, ok := img.(image.Image); ok {
		gray = imaging.Grayscale(s)
	} else {
		return nil, errors.New("不支持的图片类型")
	}
	// 将灰度图转换为矩阵
	metaData := mat.NewDense(gray.Bounds().Dy(), gray.Bounds().Dx(), nil)
	for y := 0; y < gray.Bounds().Dy(); y++ {
		for x := 0; x < gray.Bounds().Dx(); x++ {
			r, _, _, _ := gray.At(x, y).RGBA()
			metaData.Set(y, x, float64(r))
		}
	}
	return metaData, nil
}

func DrawEdges(img image.Image, minX, minY, maxX, maxY int) []byte {
	// 复制原图
	newImg := imaging.Crop(img, image.Rect(minX, minY, maxX+1, maxY+1))

	fmt.Println("裁剪坐标：", minX, minY, maxX, maxY)

	buf := new(bytes.Buffer)
	err := png.Encode(buf, newImg)
	if err != nil {
		log.Fatal(err)
	}
	return buf.Bytes()
}

func GetAllImageFiles(dir string) ([]string, error) {
	var imageFiles []string

	// 遍历目录下的所有文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断是否为图片文件
		if IsImageFile(path) {
			imageFiles = append(imageFiles, path)
		}
		return nil
	})
	fmt.Println("遍历文件：", imageFiles)
	return imageFiles, err
}

func IsImageFile(filePath string) bool {
	// 根据文件扩展名判断是否为图片文件
	ext := filepath.Ext(filePath)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	default:
		return false
	}
}
