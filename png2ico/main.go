package main

import (
	"image"
	"image/png"
	"os"
)

func main() {
	// 读取 PNG
	f, err := os.Open("favicon.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()

	// ICO 文件头 (6 bytes)
	ico := []byte{0, 0, 1, 0, 1, 0}
	// 图标目录条目 (16 bytes)
	ico = append(ico, byte(w&0xff), byte(h&0xff), 0, 0, 1, 0, 32, 0)

	// 先生成 BMP 数据 (BITMAPINFOHEADER + pixel data)
	bih := make([]byte, 40)
	// biSize
	bih[0] = 40
	// biWidth
	bih[4] = byte(w)
	bih[5] = byte(w >> 8)
	bih[6] = byte(w >> 16)
	bih[7] = byte(w >> 24)
	// biHeight = 2*h (ICO format requires double height)
	bih[8] = byte(2 * h)
	bih[9] = byte((2 * h) >> 8)
	bih[10] = byte((2 * h) >> 16)
	bih[11] = byte((2 * h) >> 24)
	// biPlanes
	bih[12] = 1
	// biBitCount
	bih[14] = 32
	// biSizeImage
	imgSize := uint32(w * h * 4)
	bih[20] = byte(imgSize)
	bih[21] = byte(imgSize >> 8)
	bih[22] = byte(imgSize >> 16)
	bih[23] = byte(imgSize >> 24)

	// 像素数据 (BGRA, bottom-up)
	pixels := make([]byte, w*h*4)
	for y := h - 1; y >= 0; y-- {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			idx := ((h-1-y)*w + x) * 4
			pixels[idx] = byte(b >> 8)
			pixels[idx+1] = byte(g >> 8)
			pixels[idx+2] = byte(r >> 8)
			pixels[idx+3] = byte(a >> 8)
		}
	}

	// AND mask (all zeros = fully visible)
	andMask := make([]byte, w*h/8)

	bmpData := append(bih, pixels...)
	bmpData = append(bmpData, andMask...)

	// 图标目录条目: offset 和 size
	dataSize := uint32(len(bmpData))
	ico = append(ico, byte(dataSize), byte(dataSize>>8), byte(dataSize>>16), byte(dataSize>>24))
	offset := uint32(6 + 16)
	ico = append(ico, byte(offset), byte(offset>>8), byte(offset>>16), byte(offset>>24))

	// 追加 BMP 数据
	ico = append(ico, bmpData...)

	// 写入 ICO
	out, err := os.Create("favicon.ico")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	out.Write(ico)

	// 确认 image 包被使用
	_ = image.NewRGBA
}
