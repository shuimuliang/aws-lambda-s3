package main

import (
	"fmt"
	"github.com/disintegration/imaging"
)

const Resolution480p = 480

// 等比例缩放测试
func thumbnail(inputFileName string, width int, height int, outFilename string) {
	srcImage, err := imaging.Open(inputFileName)
	if err != nil {
		panic(err)
	}
	// srcImage = imaging.Sharpen(srcImage, 1.0)

	dstImage := imaging.Thumbnail(srcImage, width, height, imaging.Lanczos)
	imaging.Save(dstImage, outFilename)
}

func getResizeBoundWidth(widthIn int, heightIn int, defaultResolution int) (int, int) {
	var radio float32
	var widthOut, heightOut int

	radio = float32(widthIn) / float32(heightIn)
	widthOut = defaultResolution
	heightOut = int(float32(widthOut) / radio)

	// 偶数像素对齐
	if heightOut % 2 == 1 {
		heightOut -= 1
	}

	return widthOut, heightOut
}

func getResizeBoundHeight(widthIn int, heightIn int, defaultResolution int) (int, int) {
	var radio float32
	var widthOut, heightOut int

	radio = float32(heightIn) / float32(widthIn)
	heightOut = defaultResolution
	widthOut = int(float32(heightOut) / radio)

	// 偶数像素对齐
	if widthOut % 2 == 1 {
		widthOut -= 1
	}

	return widthOut, heightOut
}

func getResizeBound(widthIn int, heightIn int, defaultResolution int) (int, int) {

	if heightIn < widthIn {
		// 横屏
		return getResizeBoundHeight(widthIn, heightIn, defaultResolution)
	} else {
		// 竖屏
		return getResizeBoundWidth(widthIn, heightIn, defaultResolution)
	}

}

func gen480pCover(inputFileName string, outFilename string) {
	srcImage, err := imaging.Open(inputFileName)
	if err != nil{}

	srcBounds := srcImage.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()
	fmt.Println(srcW, srcH)

	// destW, destH := getResizeBound(srcW, srcH, Resolution480p)
	destW, destH := 480, 300

	dstImage := imaging.Thumbnail(srcImage, destW, destH, imaging.Lanczos)
	imaging.Save(dstImage, outFilename)
}

func foo() {
	inputFileName := "/tmp/enterprise.jpg"
	thumbnail(inputFileName, 164, 100, "/tmp/out164x100.jpg")
	thumbnail(inputFileName, 70, 70, "/tmp/out70x70.jpg")
	thumbnail(inputFileName, 96, 96, "/tmp/out96x96.jpg")
	thumbnail(inputFileName, 48, 48, "/tmp/out48x48.jpg")
	thumbnail(inputFileName, 24, 24, "/tmp/out24x24.jpg")
}

func foo2() {
	inputFileName := "/tmp/enterprise.jpg"
	outFilename := "/tmp/enterprise_cover.jpg"
	gen480pCover(inputFileName, outFilename)
}

func main() {
    foo()
	foo2()
}
