package mydraw

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"github.com/xxjwxc/public/mylog"
)

// 获取画笔
func OnGetPen(fontPath string, R, G, B, A uint8) (pen Pen, b bool) {
	b = false
	pen.Color = image.NewUniform(color.RGBA{R: R, G: G, B: B, A: A})
	pen.Dpi = 72
	pen.FontSize = 10
	pen.StartPoint = image.Point{0, 0}
	// 读字体数据
	fontBytes, err := ioutil.ReadFile(fontPath)
	if err != nil {
		log.Println(err)
		return
	}
	pen.Font, err = freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println(err)
		return
	}
	b = true
	return
}

func GetImg(imagePath string) (img image.Image, err error) {
	file, _ := os.Open(imagePath)
	defer file.Close()
	//var err error
	img, _, err = image.Decode(file)
	if err != nil {
		fmt.Println("err = ", err)
		return nil, err
	}

	return img, nil
}

// Resize 设置图片高宽
func Resize(img image.Image, width, height uint) image.Image {
	// dx := img.Bounds().Dx()
	// dy := img.Bounds().Dy()
	return resize.Resize(width, height, img, resize.Lanczos3)
}

type Pen struct {
	FontSize   float64
	Dpi        float64
	Font       *truetype.Font
	StartPoint image.Point
	Color      *image.Uniform
}

type HDC struct {
	//Bg   image.Image
	Rgba *image.RGBA
}

func (h *HDC) SetBg(imagePath string) bool {
	file, _ := os.Open(imagePath)
	defer file.Close()
	//var err error
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("err = ", err)
		return false
	}

	h.Rgba = image.NewRGBA(img.Bounds())
	draw.Draw(h.Rgba, h.Rgba.Bounds(), img, image.ZP, draw.Src)
	return true
}

func (h *HDC) GetBgSize() (w, _h int) {
	b := h.Rgba.Bounds()
	w = b.Max.X
	_h = b.Max.Y
	return
}

// 图片上画文字
func (h *HDC) DrawText(pen Pen, text string) bool {
	if h.Rgba == nil {
		return false
	}

	c := freetype.NewContext()
	c.SetDPI(pen.Dpi)
	c.SetFont(pen.Font)
	c.SetFontSize(pen.FontSize)
	c.SetClip(h.Rgba.Bounds())
	c.SetDst(h.Rgba)
	//c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255}))
	c.SetSrc(pen.Color)

	// Draw the text.
	pt := freetype.Pt(pen.StartPoint.X, pen.StartPoint.Y+int(c.PointToFixed(pen.FontSize)>>6))
	for _, s := range strings.Split(text, "\r\n") {
		_, err := c.DrawString(s, pt)
		if err != nil {
			mylog.Infof("c.DrawString(%s) error(%v)", s, err)
			return false
		}
		pt.Y += c.PointToFixed(pen.FontSize * 1.5)
	}
	return false
}

// 图片上画图片
func (h *HDC) DrawImg(img image.Image, point image.Point) {
	draw.Draw(h.Rgba, h.Rgba.Bounds().Add(point),
		img,
		img.Bounds().Min,
		draw.Over)
}

// 保存图片
func (h *HDC) Save(imagePath string) bool {
	output, err := os.OpenFile(imagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		mylog.Error(err)
		return false
	}
	defer output.Close()

	if strings.HasSuffix(imagePath, ".png") || strings.HasSuffix(imagePath, ".PNG") {
		err = png.Encode(output, h.Rgba)
	} else {
		err = jpeg.Encode(output, h.Rgba, nil)
	}
	if err != nil {
		mylog.Infof("image encode error(%v)", err)
		return false
	}
	return true
}
