package main

import (
	"encoding/json"
	"image"
	"imgSearcher/imgHandle"
	"log"
	"math/bits"
	"os"
	"path/filepath"
	"sync"

	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
)

func main() {
	service := NewImageHashService("./static/images")
	if err := service.Initialize(); err != nil {
		log.Fatal("Failed to initialize service:", err)
	}

	r := setupRouter(service)
	r.Run(":8080")
}

type ImageHashService struct {
	HashList  map[string]uint64 // 文件名 -> hash值
	mutex     sync.RWMutex
	hashFile  string // hash值持久化文件路径
	imagePath string // 图片目录
}

func NewImageHashService(imagePath string) *ImageHashService {
	return &ImageHashService{
		HashList:  make(map[string]uint64),
		hashFile:  "image_hashes.json",
		imagePath: imagePath,
	}
}

func (s *ImageHashService) Initialize() error {
	// 尝试加载已存在的hash值
	if err := s.loadHashes(); err != nil {
		return err
	}

	// 处理目录中的图片
	return filepath.Walk(s.imagePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && isImageFile(path) {
			if _, exists := s.HashList[info.Name()]; !exists {
				if err := s.processAndStoreImage(path, info.Name()); err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func isImageFile(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".png" || ext == ".jpg" || ext == ".jpeg" || ext == ".bmp" || ext == ".gif"
}

func (s *ImageHashService) processAndStoreImage(path, filename string) error {
	// 打开并处理图片
	img, err := imaging.Open(path)
	if err != nil {
		return err
	}

	// 裁剪处理
	processed := imgHandle.Crop(&img, imgHandle.EdgeDetect(&img))

	// 计算hash
	hash, err := imgHandle.CalculateImageHash(processed)
	if err != nil {
		return err
	}

	log.Println("filename:", filename, "hash:", hash)

	// 存储hash
	s.mutex.Lock()
	s.HashList[filename] = hash
	s.mutex.Unlock()

	// 保存到文件
	return s.saveHashes()
}

func (s *ImageHashService) loadHashes() error {
	data, err := os.ReadFile(s.hashFile)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()
	return json.Unmarshal(data, &s.HashList)
}

func (s *ImageHashService) saveHashes() error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	data, err := json.MarshalIndent(s.HashList, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.hashFile, data, 0644)
}

type Match struct {
	Filename   string `json:"filename"`
	Similarity int    `json:"similarity"`
}

func setupRouter(service *ImageHashService) *gin.Engine {
	r := gin.Default()

	// 提供静态文件服务
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// 首页
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 图片上传和匹配
	r.POST("/match", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(400, gin.H{"error": "No image uploaded"})
			return
		}

		// 打开上传的文件
		src, err := file.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer src.Close()

		// 解码图片
		img, _, err := image.Decode(src)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to decode image"})
			return
		}

		// 处理图片
		processed := imgHandle.Crop(&img, imgHandle.EdgeDetect(&img))

		// 计算hash
		hash, err := imgHandle.CalculateImageHash(processed)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to calculate hash"})
			return
		}

		// 查找匹配
		matches := make([]Match, 0)
		service.mutex.RLock()
		for filename, storedHash := range service.HashList {
			similarity := calculateSimilarity(hash, storedHash)
			if similarity <= 100 { // 设置相似度阈值
				matches = append(matches, Match{
					Filename:   filename,
					Similarity: similarity,
				})
			}
		}
		service.mutex.RUnlock()

		c.JSON(200, matches)
	})

	return r
}

func calculateSimilarity(hash, storedHash uint64) int {
	distance := bits.OnesCount64(hash ^ storedHash)
	return distance
}
