package imgHandle

import (
	"image"

	"github.com/corona10/goimagehash"
)

const HASH_SIZE = 16

type Hashs struct {
	PerceptionHash []uint64
	AverageHash    []uint64
	DifferenceHash []uint64
}

// 计算图像的感知哈希
func CalculateImageHash(img *image.Image) (Hashs, error) {
	hash, err := goimagehash.ExtPerceptionHash(*img, HASH_SIZE, HASH_SIZE)
	if err != nil {
		return Hashs{}, err
	}
	averageHash, err := goimagehash.ExtAverageHash(*img, HASH_SIZE, HASH_SIZE)
	if err != nil {
		return Hashs{}, err
	}
	differenceHash, err := goimagehash.ExtDifferenceHash(*img, HASH_SIZE, HASH_SIZE)
	if err != nil {
		return Hashs{}, err
	}
	return Hashs{
		PerceptionHash: hash.GetHash(),
		AverageHash:    averageHash.GetHash(),
		DifferenceHash: differenceHash.GetHash(),
	}, nil
}

// 计算两个图像的相似度（0-1之间，1表示完全相同）
func CompareImages(img1, img2 *image.Image) (int, error) {
	hash1, err := goimagehash.ExtPerceptionHash(*img1, HASH_SIZE, HASH_SIZE)
	if err != nil {
		return 0, err
	}

	hash2, err := goimagehash.ExtPerceptionHash(*img2, HASH_SIZE, HASH_SIZE)
	if err != nil {
		return 0, err
	}

	distance, err := hash1.Distance(hash2)
	if err != nil {
		return 0, err
	}

	return distance, nil
}
