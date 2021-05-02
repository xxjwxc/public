package mydraw

import (
	"fmt"
	"testing"
)

func TestMytest(t *testing.T) {
	pen, b := OnGetPen("./luximr.ttf", 0, 0, 0, 255)
	if b {
		var hdc HDC
		hdc.SetBg("./src.png")
		pen.Dpi = 200
		pen.FontSize = 16
		pen.StartPoint.X = 150
		pen.StartPoint.Y = 78
		hdc.DrawText(pen, "哈哈")
		pen.FontSize = 12
		pen.StartPoint.X = 150
		pen.StartPoint.Y = 160
		hdc.DrawText(pen, "男")
		pen.StartPoint.X = 350
		pen.StartPoint.Y = 160
		hdc.DrawText(pen, "汉")
		pen.StartPoint.X = 150
		pen.StartPoint.Y = 240
		hdc.DrawText(pen, "1996")
		pen.StartPoint.X = 300
		pen.StartPoint.Y = 240
		hdc.DrawText(pen, "6")
		pen.StartPoint.X = 370
		pen.StartPoint.Y = 240
		hdc.DrawText(pen, "26")

		pen.StartPoint.X = 150
		pen.StartPoint.Y = 275

		str := []rune("北京市海淀区西北旺东路100号中关村科技园")
		//str := "北京市海淀区西北旺东路100号中关村科技园"

		for i := 0; i < len(str); i += 11 {
			var end = i + 11
			if end > len(str) {
				end = len(str)
			}
			pen.StartPoint.Y += 40
			tmp := str[i:end]
			hdc.DrawText(pen, string(tmp))
		}

		pen.StartPoint.X = 300
		pen.StartPoint.Y = 520
		hdc.DrawText(pen, "310666196606266666")

		b = hdc.Save("./out.png")
		fmt.Println(b)
	}
}
