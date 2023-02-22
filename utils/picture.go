package utils

import (
	"fmt"
	"image/color"
	"io"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/fogleman/gg"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

const fontPath = "DingTalk JinBuTi.ttf"

var font *truetype.Font

var once sync.Once

func Init() {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
	}
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
}

func String2Pic(text, picPath string) error {
	once.Do(Init)
	// 创建一个 800x600 的图像
	dc := gg.NewContext(800, 40*countNewLine(text))

	// 设置背景色为白色
	dc.SetColor(color.White)
	dc.Clear()

	// 设置字体大小和颜色
	face := truetype.NewFace(font, &truetype.Options{Size: 24})
	dc.SetFontFace(face)
	dc.SetColor(color.Black)

	// 在图像中绘制文本
	dc.DrawStringWrapped(text, 20, 50, 0, 0, 800, 1.5, gg.AlignLeft)

	// 将图像保存为 PNG 文件
	if err := dc.SavePNG(picPath); err != nil {
		return fmt.Errorf("error saving PNG: %w", err)
	}

	return nil
}

func String2PicWriter(text string, w io.Writer) error {
	once.Do(Init)
	// 创建一个 800x600 的图像
	dc := gg.NewContext(800, 40*countNewLine(text))

	// 设置背景色为白色
	dc.SetColor(color.White)
	dc.Clear()

	// 设置字体大小和颜色
	face := truetype.NewFace(font, &truetype.Options{Size: 24})
	dc.SetFontFace(face)
	dc.SetColor(color.Black)

	// 在图像中绘制文本
	dc.DrawStringWrapped(text, 20, 50, 0, 0, 800, 1.5, gg.AlignLeft)

	return dc.EncodePNG(w)
}

// 计算字符串中有多少个换行
func countNewLine(str string) int {
	return strings.Count(str, "\n")
}
