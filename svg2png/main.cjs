const sharp = require('sharp');
const fs = require('fs/promises');
const path = require('path');

const MIN_SIZE = 100;  // 最小尺寸
const PADDING = 4;     // 边框大小

async function processImage(svgBuffer) {
    // 获取图像信息
    const metadata = await sharp(svgBuffer).metadata();
    
    // 计算新的尺寸，确保至少一边为100像素
    let scale = 1;
    if (metadata.width < MIN_SIZE && metadata.height < MIN_SIZE) {
        scale = Math.max(
            MIN_SIZE / metadata.width,
            MIN_SIZE / metadata.height
        );
    }
    
    const newWidth = Math.round(metadata.width * scale);
    const newHeight = Math.round(metadata.height * scale);
    
    // 添加padding并转换为PNG
    return await sharp(svgBuffer)
        .resize(newWidth, newHeight)
        .extend({
            top: PADDING,
            bottom: PADDING,
            left: PADDING,
            right: PADDING,
            background: { r: 0, g: 0, b: 0, alpha: 0 }
        })
        .png()
        .toBuffer();
}

async function convertSVGtoPNG(inputPath, outputPath) {
    try {
        const svgBuffer = await fs.readFile(inputPath);
        const pngBuffer = await processImage(svgBuffer);
        await fs.writeFile(outputPath, pngBuffer);
        console.log(`Converted: ${inputPath} -> ${outputPath}`);
    } catch (err) {
        console.error(`Error converting ${inputPath}:`, err);
    }
}

async function processDirectory(inputDir, outputDir) {
    try {
        // 递归获取所有文件
        async function* getFiles(dir) {
            const items = await fs.readdir(dir, { withFileTypes: true });
            for (const item of items) {
                const fullPath = path.join(dir, item.name);
                if (item.isDirectory()) {
                    yield* getFiles(fullPath);
                } else if (item.isFile() && item.name.toLowerCase().endsWith('.svg')) {
                    yield fullPath;
                }
            }
        }

        // 处理每个SVG文件
        for await (const filePath of getFiles(inputDir)) {
            // 保持相对路径结构
            const relativePath = path.relative(inputDir, filePath);
            const outputPath = path.join(
                outputDir,
                relativePath.replace(/\.svg$/i, '.png')
            );

            // 确保输出目录存在
            await fs.mkdir(path.dirname(outputPath), { recursive: true });
            
            // 转换文件
            await convertSVGtoPNG(filePath, outputPath);
        }

        console.log('All conversions completed!');
    } catch (err) {
        console.error('Error processing directory:', err);
    }
}

// 命令行参数处理
const args = process.argv.slice(2);
if (args.length !== 2) {
    console.log('Usage: node convert.js <input_dir> <output_dir>');
    process.exit(1);
}

const [inputDir, outputDir] = args;
processDirectory(inputDir, outputDir);