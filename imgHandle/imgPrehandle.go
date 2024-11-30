package imgHandle

import (
	"image"
	"image/color"

	"github.com/disintegration/imaging"
)

// 获取背景色并将其设置为透明
func RemoveBackground(image *image.Image, threshold uint8) *image.NRGBA {
	img := imaging.Clone(*image)
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 获取边缘颜色样本
	backgroundColor := detectBackgroundColor(img)

	// 创建新的透明背景图像
	newImg := imaging.New(width, height, color.NRGBA{0, 0, 0, 0})

	// 处理每个像素
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := img.NRGBAAt(x, y)
			// 如果像素颜色接近背景色，则设置为透明
			if isColorSimilar(pixel, backgroundColor, threshold) {
				newImg.Set(x, y, color.NRGBA{0, 0, 0, 0})
			} else {
				newImg.Set(x, y, pixel)
			}
		}
	}

	return newImg
}

// 检测背景色（通过采样边缘像素）
func detectBackgroundColor(img *image.NRGBA) color.NRGBA {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// 采样点数量
	sampleSize := 10
	samples := make([]color.NRGBA, 0, sampleSize*4)

	// 采样四个边
	for i := 0; i < sampleSize; i++ {
		// 上边
		x := (width * i) / sampleSize
		samples = append(samples, img.NRGBAAt(x, 0))

		// 下边
		samples = append(samples, img.NRGBAAt(x, height-1))

		// 左边
		y := (height * i) / sampleSize
		samples = append(samples, img.NRGBAAt(0, y))

		// 右边
		samples = append(samples, img.NRGBAAt(width-1, y))
	}

	// 计算平均颜色
	var sumR, sumG, sumB, sumA uint32
	for _, c := range samples {
		sumR += uint32(c.R)
		sumG += uint32(c.G)
		sumB += uint32(c.B)
		sumA += uint32(c.A)
	}

	count := uint32(len(samples))
	return color.NRGBA{
		R: uint8(sumR / count),
		G: uint8(sumG / count),
		B: uint8(sumB / count),
		A: uint8(sumA / count),
	}
}

// 判断两个颜色是否相似
func isColorSimilar(c1, c2 color.NRGBA, threshold uint8) bool {
	diffR := absDiff(c1.R, c2.R)
	diffG := absDiff(c1.G, c2.G)
	diffB := absDiff(c1.B, c2.B)

	return diffR <= threshold &&
		diffG <= threshold &&
		diffB <= threshold
}

// 计算绝对差值
func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}
