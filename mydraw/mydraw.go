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
	"public/mylog"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
)

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

//获取画笔
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

func (this *HDC) SetBg(imagePath string) bool {
	file, _ := os.Open(imagePath)
	defer file.Close()
	//var err error
	img, _, err := image.Decode(file)
	if err != nil {
		fmt.Println("err = ", err)
		return false
	}

	this.Rgba = image.NewRGBA(img.Bounds())
	draw.Draw(this.Rgba, this.Rgba.Bounds(), img, image.ZP, draw.Src)
	return true
}

func (this *HDC) GetBgSize() (w, h int) {
	b := this.Rgba.Bounds()
	w = b.Max.X
	h = b.Max.Y
	return
}

//图片上画文字
func (this *HDC) DrawText(pen Pen, text string) bool {
	if this.Rgba == nil {
		return false
	}

	c := freetype.NewContext()
	c.SetDPI(pen.Dpi)
	c.SetFont(pen.Font)
	c.SetFontSize(pen.FontSize)
	c.SetClip(this.Rgba.Bounds())
	c.SetDst(this.Rgba)
	//c.SetSrc(image.NewUniform(color.RGBA{255, 255, 255, 255}))
	c.SetSrc(pen.Color)

	// Draw the text.
	pt := freetype.Pt(pen.StartPoint.X, pen.StartPoint.Y+int(c.PointToFixed(pen.FontSize)>>6))
	for _, s := range strings.Split(text, "\r\n") {
		_, err := c.DrawString(s, pt)
		if err != nil {
			mylog.Println("c.DrawString(%s) error(%v)", s, err)
			return false
		}
		pt.Y += c.PointToFixed(pen.FontSize * 1.5)
	}
	return false
}

//保存图片
func (this *HDC) Save(imagePath string) bool {
	output, err := os.OpenFile(imagePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		mylog.Error(err)
		return false
	}
	defer output.Close()

	if strings.HasSuffix(imagePath, ".png") || strings.HasSuffix(imagePath, ".PNG") {
		err = png.Encode(output, this.Rgba)
	} else {
		err = jpeg.Encode(output, this.Rgba, nil)
	}
	if err != nil {

		mylog.Println("image encode error(%v)", err)
		//mylog.Error(err)
		return false
	}
	return true
}
