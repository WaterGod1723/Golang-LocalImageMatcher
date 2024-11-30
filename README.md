# Image Search Tool / 图像搜索工具

[English](#english) | [中文](#中文)

## English

### Problem Solved
As frontend developers, we often face these challenges:
- Need to find the exact name of an icon from a large icon library
- Different team members use different naming conventions
- Multiple duplicate icons exist in the codebase
- Time wasted searching through hundreds of icons manually

This tool helps you quickly find icons and their filenames by:
- Simply uploading a screenshot or pasting from clipboard
- Finding similar icons even if they're not exact matches
- Showing all variations of the same icon with their filenames

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

### 解决的问题
作为前端开发人员，我们经常面临这些挑战：
- 需要从大量图标库中找到目标图标的确切名称
- 团队成员使用不同的命名规范
- 代码库中存在多个重复的图标
- 手动搜索数百个图标浪费时间

这个工具通过以下方式帮助你快速找到图标及其文件名：
- 简单地上传截图或从剪贴板粘贴（Ctrl+V）
- 即使不是完全匹配也能找到相似图标
- 显示同一图标的所有变体及其文件名

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