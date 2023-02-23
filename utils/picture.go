package utils

import (
	"fmt"
	"image"
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

func String2PicWriterWithTailPicture(text, fontPath string, w io.Writer, tailPicture image.Image) error {
	dc := string2PicDCWithTailPicture(text, fontPath, tailPicture)

	return dc.EncodePNG(w)
}

func String2PicFile(text, fontPath, picPath string) error {
	dc := string2PicDC(text, fontPath)
	// 将图像保存为 PNG 文件
	if err := dc.SavePNG(picPath); err != nil {
		return fmt.Errorf("error saving PNG: %w", err)
	}

	return nil
}

func String2PicFileWithTailPicture(text, fontPath, picPath string, tailPicture image.Image) error {
	dc := string2PicDCWithTailPicture(text, fontPath, tailPicture)
	// 将图像保存为 PNG 文件
	if err := dc.SavePNG(picPath); err != nil {
		return fmt.Errorf("error saving PNG: %w", err)
	}

	return nil
}

func string2PicDC(text, fontPath string) *gg.Context {
	once.Do(func() {
		Init(fontPath)
	})
	// 把连续的换行符替换成一个换行符，因为转图片的时候gg库会自动去换行，导致后面的图片会被挤到下面很多行
	text = strings.ReplaceAll(text, "\n\n", "\n")
	// 简单换两次吧，用来处理连续三个换行的情况，四个以上的就不管了
	text = strings.ReplaceAll(text, "\n\n", "\n")
	fontsize, lineSpacing := 24.0, 1.5
	startBlankHeight := 50.0 // 最前面的空白的高度
	width := 800
	// 创建一个 800x600 的图像
	dc := gg.NewContext(width, int(fontsize*lineSpacing*float64(countNewLine(text))+startBlankHeight))

	// 设置背景色为白色
	dc.SetColor(color.White)
	dc.Clear()

	// 设置字体大小和颜色
	face := truetype.NewFace(font, &truetype.Options{Size: fontsize})
	dc.SetFontFace(face)
	dc.SetColor(color.Black)

	// 在图像中绘制文本
	dc.DrawStringWrapped(text, 20, startBlankHeight, 0, 0, float64(width), lineSpacing, gg.AlignLeft)

	return dc
}

func string2PicDCWithTailPicture(text, fontPath string, tailPicture image.Image) *gg.Context {
	once.Do(func() {
		Init(fontPath)
	})
	// 把连续的换行符替换成一个换行符，因为转图片的时候gg库会自动去换行，导致后面的图片会被挤到下面很多行
	text = strings.ReplaceAll(text, "\n\n", "\n")
	// 简单换两次吧，用来处理连续三个换行的情况，四个以上的就不管了
	text = strings.ReplaceAll(text, "\n\n", "\n")
	fontsize, lineSpacing := 24.0, 1.5
	startBlankHeight := 50.0 // 最前面的空白的高度
	width := 800
	// 创建一个 800x600 的图像
	dc := gg.NewContext(width, int(fontsize*lineSpacing*float64(countNewLine(text))+float64(tailPicture.Bounds().Dy())+startBlankHeight))

	// 设置背景色为白色
	dc.SetColor(color.White)
	dc.Clear()

	// 设置字体大小和颜色
	face := truetype.NewFace(font, &truetype.Options{Size: fontsize})
	dc.SetFontFace(face)
	dc.SetColor(color.Black)

	// 在图像中绘制文本
	dc.DrawStringWrapped(text, 20, startBlankHeight, 0, 0, float64(width), lineSpacing, gg.AlignLeft)
	dc.DrawImage(tailPicture, 0, int(startBlankHeight+fontsize*lineSpacing*float64(countNewLine(text)))+5)

	return dc
}

// 计算字符串中有多少个换行
func countNewLine(str string) int {
	return strings.Count(str, "\n")
}
