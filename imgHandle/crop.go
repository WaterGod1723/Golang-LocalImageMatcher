package imgHandle

import (
	"image"
	"image/color"
	"image/draw"
	"math"

	"github.com/disintegration/imaging"
)

func LoadImg(path string) (*image.Image, error) {
	img, err := imaging.Open(path)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

type EdgeDetectResult struct {
	MinX int
	MaxX int
	MinY int
	MaxY int
}

func EdgeDetect(img *image.Image) EdgeDetectResult {
	grayImg := imaging.Grayscale(*img)
	bounds := grayImg.Bounds()
	output := imaging.New(bounds.Dx(), bounds.Dy(), color.Gray{})

	var minX, maxX, minY, maxY int
	minX = bounds.Max.X
	maxX = bounds.Min.X
	minY = bounds.Max.Y
	maxY = bounds.Min.Y

	const threshold = 70

	// sobel算子
	x := [][]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	y := [][]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	// 处理边界问题，需要跳过图像边缘的像素
	for i := 1; i < bounds.Max.X-1; i++ {
		for j := 1; j < bounds.Max.Y-1; j++ {
			var sumX, sumY float64

			// 执行3x3卷积
			for k := -1; k <= 1; k++ {
				for l := -1; l <= 1; l++ {
					pixel := grayImg.At(i+k, j+l)
					gray, _, _, _ := pixel.RGBA()
					grayFloat := float64(gray >> 8) // 转换为0-255范围

					sumX += x[k+1][l+1] * grayFloat
					sumY += y[k+1][l+1] * grayFloat
				}
			}

			// 计算梯度幅值
			magnitude := math.Sqrt(sumX*sumX + sumY*sumY)
			// 归一化到0-255
			magnitude = math.Min(255, magnitude)

			if magnitude > threshold {
				if i < minX {
					minX = i
				}
				if i > maxX {
					maxX = i
				}
				if j < minY {
					minY = j
				}
				if j > maxY {
					maxY = j
				}
			}

			output.Set(i, j, color.Gray{uint8(magnitude)})
		}
	}
	maxX++
	maxY++
	return EdgeDetectResult{minX, maxX, minY, maxY}
}

func Crop(img *image.Image, result EdgeDetectResult) *image.Image {
	// 1. 先裁剪到边缘检测区域
	croppedImg := cropImage(img, result)
	// 2. 移除背景色（设置阈值为15）
	croppedImg = RemoveBackground(&croppedImg, 15)
	// 3. 确保最小尺寸
	scaledImg := scaleImage(croppedImg)

	// 4. 添加透明边框并居中
	finalImg := addBorder(scaledImg, color.NRGBA{0, 0, 0, 0})

	// 转换回 *image.NRGBA
	return &finalImg
}

func cropImage(img *image.Image, ed EdgeDetectResult) image.Image {
	rect := image.Rect(ed.MinX, ed.MinY, ed.MaxX, ed.MaxY)
	croppedImg := image.NewRGBA(rect)
	draw.Draw(croppedImg, rect, *img, image.Point{X: ed.MinX, Y: ed.MinY}, draw.Src)
	return croppedImg
}

func scaleImage(img image.Image) image.Image {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 确保至少有一边不小于50像素
	minSize := 50
	if width < minSize || height < minSize {
		scale := float64(minSize) / float64(width)
		if height < width {
			scale = float64(minSize) / float64(height)
		}

		newWidth := int(float64(width) * scale)
		newHeight := int(float64(height) * scale)

		return imaging.Resize(img, newWidth, newHeight, imaging.Lanczos)
	}

	return img
}

func addBorder(img image.Image, borderColor color.Color) image.Image {
	const MIN_PADDING = 5
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 取最大边作为最终尺寸，并添加padding
	finalSize := int(math.Max(float64(width), float64(height))) + (MIN_PADDING * 2)

	// 创建新的透明背景图像
	newImg := imaging.New(finalSize, finalSize, borderColor)

	// 计算居中位置
	x := (finalSize - width) / 2
	y := (finalSize - height) / 2

	// 将原图绘制到新图像的中心位置
	return imaging.Paste(newImg, img, image.Point{x, y})
}
