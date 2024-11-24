package imgHandle

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// 计算两个64位整数的汉明距离
func hammingDistance(x, y uint64) int {
	var dist int
	var val uint64 = x ^ y
	for val != 0 {
		dist++
		val &= val - 1
	}
	return dist
}

// 计算感知哈希
func perceptionHash(image *mat.Dense) uint64 {
	// 将64x64矩阵缩小为8x8
	smaller := mat.NewDense(8, 8, nil)
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			var sum float64
			for k := 0; k < 8; k++ {
				for l := 0; l < 8; l++ {
					sum += image.At(i*8+k, j*8+l)
				}
			}
			smaller.Set(i, j, sum/64)
		}
	}

	// 计算DCT（这里简化处理，实际需要更复杂的算法）
	var dct mat.Dense
	for u := 0; u < 8; u++ {
		for v := 0; v < 8; v++ {
			var sum float64
			for x := 0; x < 8; x++ {
				for y := 0; y < 8; y++ {
					cu := 1.0
					cv := 1.0
					if u == 0 {
						cu = 1 / math.Sqrt2
					}
					if v == 0 {
						cv = 1 / math.Sqrt2
					}
					sum += cu * cv * smaller.At(x, y) * math.Cos((float64(u)*(2*float64(x)+1)*math.Pi)/16) * math.Cos((float64(v)*(2*float64(y)+1)*math.Pi)/16)
				}
			}
			dct.Set(u, v, sum)
		}
	}

	// 取DCT系数并计算平均值
	var sum float64
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			sum += dct.At(i, j)
		}
	}
	avg := sum / 64

	// 生成哈希
	var hash uint64
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			bitIndex := i*8 + j
			if dct.At(i, j) > avg {
				hash |= 1 << bitIndex
			}
		}
	}

	return hash
}

// 计算两个图片的相似度
func ImageSimilarity(img1, img2 *mat.Dense) float64 {
	hash1 := perceptionHash(img1)
	hash2 := perceptionHash(img2)
	dist := hammingDistance(hash1, hash2)
	return 1 - float64(dist)/64.0
}
