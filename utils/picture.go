package utils

import (
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

var font *truetype.Font

var once sync.Once

// Init is used to load font face
func Init(fontPath string) {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatal(err)
	}
	font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatal(err)
	}
}

func String2PicWriter(text, fontPath string, w io.Writer) error {
	dc := string2PicDC(text, fontPath)

	return dc.EncodePNG(w)
}

//func String2PicFile(text, fontPath, picPath string) error {
//	dc := string2PicDC(text, fontPath)
//	// 将图像保存为 PNG 文件
//	if err := dc.SavePNG(picPath); err != nil {
//		return fmt.Errorf("error saving PNG: %w", err)
//	}
//
//	return nil
//}

func string2PicDC(text, fontPath string) *gg.Context {
	once.Do(func() {
		Init(fontPath)
	})
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

	return dc
}

// 计算字符串中有多少个换行
func countNewLine(str string) int {
	return strings.Count(str, "\n")
}
