package imgHandle

import (
	"image"

	"github.com/corona10/goimagehash"
)

// 计算图像的感知哈希
func CalculateImageHash(img *image.Image) (uint64, error) {
	hash, err := goimagehash.PerceptionHash(*img)
	if err != nil {
		return 0, err
	}
	return hash.GetHash(), nil
}

// 计算两个图像的相似度（0-1之间，1表示完全相同）
func CompareImages(img1, img2 *image.Image) (int, error) {
	hash1, err := goimagehash.PerceptionHash(*img1)
	if err != nil {
		return 0, err
	}

	hash2, err := goimagehash.PerceptionHash(*img2)
	if err != nil {
		return 0, err
	}

	distance, err := hash1.Distance(hash2)
	if err != nil {
		return 0, err
	}

	return distance, nil
}
