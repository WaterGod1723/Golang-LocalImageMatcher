# Image Search Tool / 图像搜索工具

[English](#english) | [中文](#中文)

## English

### Features
- Local image reverse search
- Similar image search
- Support for image library management (e.g., searching frontend icons by image)

### Prerequisites
1. Go environment
2. Node.js environment (for SVG to PNG conversion)

### Usage
1. (Optional) Image preprocessing: Convert SVG images to PNG format

bash

node convert.js <input_dir> <output_dir>


The converted PNG files will:
- Have transparent background
- Include 4px padding
- Maintain aspect ratio
- Have minimum dimension of 100px

2. Copy your image library to `static/img/library`
3. Run `main.go`
4. Access `http://localhost:8080`
5. When updating the image library, manually clear the `image_hashes.json` cache file and repeat from step 1

### Directory Structure
.
├── static/
│ └── images/ # Image library directory
├── templates/
│ └── index.html # Web interface template
└── image_hashes.json # Hash cache file

---

## 中文

### 功能特性
- 本地图片反向搜索
- 以图找图
- 支持图片素材管理（如前端icon图片搜图，以实现以图片查名字）

### 环境准备
1. Go环境
2. Node.js环境（用于svg转png）

### 使用方式
1. （如果有svg图片）图片预处理：将svg图片转换为png图片

bash
node convert.js <输入目录> <输出目录>

转换后的PNG文件将：
- 使用透明背景
- 包含4像素内边距
- 保持原始比例
- 确保最小尺寸为100像素

2. 将图库复制到 static/img/库中
3. 运行 main.go
4. 访问 http://localhost:8080
5. 如果图库有刷新需要手动清除 image_hashes.json 缓存文件并重新从步骤1开始执行

### 目录结构
.
├── static/
│ └── images/ # 图库目录
├── templates/
│ └── index.html # Web界面模板
└── image_hashes.json # 哈希缓存文件