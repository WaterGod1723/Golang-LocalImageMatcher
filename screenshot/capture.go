package screenshot

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"
	"unsafe"

	"github.com/atotto/clipboard"
	"github.com/lxn/win"
)

func WatchClipboard() {
	// 定义保存图片的路径
	filePath := "clipboard_image.png"

	// 监听剪贴板变化
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	lastClipboardContent := ""

	for range ticker.C {
		// 打开剪贴板
		if !win.OpenClipboard(0) {
			log.Printf("Failed to open clipboard")
			continue
		}
		defer win.CloseClipboard()

		// 检查剪贴板中是否有图像
		hBitmap := win.GetClipboardData(win.CF_BITMAP)
		if hBitmap == 0 {
			// 如果没有图像，检查是否有文本内容
			text, err := clipboard.ReadAll()
			if err != nil {
				log.Printf("Failed to read clipboard text: %v", err)
				continue
			}
			if text != lastClipboardContent {
				lastClipboardContent = text
				log.Printf("Clipboard content is text: %s", text)
			}
			continue
		}

		// 获取图像的设备上下文
		hDC := win.GetDC(0)
		if hDC == 0 {
			log.Printf("Failed to get device context")
			continue
		}
		defer win.ReleaseDC(0, hDC)

		memDC := win.CreateCompatibleDC(hDC)
		if memDC == 0 {
			log.Printf("Failed to create compatible DC")
			continue
		}
		defer win.DeleteDC(memDC)

		oldBitmap := win.SelectObject(memDC, win.HGDIOBJ(hBitmap))
		if oldBitmap == 0 {
			log.Printf("Failed to select object into DC")
			continue
		}
		defer win.SelectObject(memDC, oldBitmap)

		// 获取图像的尺寸
		bitmap := win.BITMAP{}
		if win.GetObject(win.HGDIOBJ(hBitmap), uintptr(unsafe.Sizeof(bitmap)), unsafe.Pointer(&bitmap)) == 0 {
			log.Printf("Failed to get bitmap object")
			continue
		}

		// 创建一个内存缓冲区来存储图像数据
		bmpBytes := make([]byte, int(bitmap.BmWidthBytes)*int(bitmap.BmHeight))

		// 获取图像数据
		bi := win.BITMAPINFO{
			BmiHeader: win.BITMAPINFOHEADER{
				BiSize:          uint32(unsafe.Sizeof(win.BITMAPINFOHEADER{})),
				BiWidth:         bitmap.BmWidth,
				BiHeight:        -bitmap.BmHeight, // 负值表示自上而下的扫描线顺序
				BiPlanes:        1,
				BiBitCount:      32,
				BiCompression:   win.BI_RGB,
				BiSizeImage:     uint32(bitmap.BmWidthBytes * int32(bitmap.BmHeight)),
				BiXPelsPerMeter: 0,
				BiYPelsPerMeter: 0,
				BiClrUsed:       0,
				BiClrImportant:  0,
			},
		}

		if win.GetDIBits(memDC, win.HBITMAP(hBitmap), 0, uint32(bitmap.BmHeight), &bmpBytes[0], &bi, win.DIB_RGB_COLORS) == 0 {
			log.Printf("Failed to get bitmap bits")
			continue
		}

		// 创建一个 RGBA 图像
		img := image.NewRGBA(image.Rect(0, 0, int(bitmap.BmWidth), int(bitmap.BmHeight)))
		for y := 0; y < int(bitmap.BmHeight); y++ {
			for x := 0; x < int(bitmap.BmWidth); x++ {
				offset := y*int(bitmap.BmWidthBytes) + x*4
				img.Set(x, y, color.RGBA{R: bmpBytes[offset], G: bmpBytes[offset+1], B: bmpBytes[offset+2], A: bmpBytes[offset+3]})
			}
		}

		// 保存图片
		file, err := os.Create(filePath)
		if err != nil {
			log.Printf("Failed to create file: %v", err)
			continue
		}
		defer file.Close()

		if err := png.Encode(file, img); err != nil {
			log.Printf("Failed to encode image: %v", err)
			continue
		}

		fmt.Println("Image saved to", filePath)
	}
}
