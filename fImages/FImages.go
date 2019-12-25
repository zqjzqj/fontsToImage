package fImages

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"
)

type FImages struct {
	*image.NRGBA
}

func NewFImages(rect image.Rectangle) *FImages {
	return &FImages{image.NewNRGBA(rect)}
}

func (img *FImages) SetBkg(bColor color.RGBA) {
	for x := img.Rect.Min.X; x < img.Rect.Max.X; x++ {
		for y := img.Rect.Min.Y; y < img.Rect.Max.Y; y++ {
			img.Set(x, y, bColor)
		}
	}
}

func sign(x int) int {
	if x > 0 {
		return 1
	}
	return -1
}

//颜色比较
func ColorIsEq(color1 color.Color, color2 color.Color) bool {
	r, g, b, a := color1.RGBA()
	r2, g2, b2, a2 := color2.RGBA()
	if r == r2 && g == g2 && b == b2 && a == a2 {
		return true
	}
	return false
}

func (img *FImages) EdgeFill(points []image.Point, rightX int) {
	var x0, y0, x1, y1 float64
	pLen := len(points)
	for k, v := range  points {
		//首尾相连
		if k + 1 > pLen - 1 {
			x0, y0, x1, y1 = float64(v.X), float64(v.Y), float64(points[0].X), float64(points[0].Y)
		} else {
			x0, y0, x1, y1 = float64(v.X), float64(v.Y), float64(points[k + 1].X), float64(points[k + 1].Y)
		}
		//如果线段是水平线或者是右边边界，不做处理
		if ( y0 == y1 ) || ( rightX == int(x1) && rightX == int(x0) ) {
			continue
		} else {
			dx, dy := (x1 - x0), (y1 - y0)

			//扫描线的范围
			eps1 := int(math.Abs(dy))

			xIncre := dx / float64(eps1)
			yIncre := dy / float64(eps1)
			for k := 0; k <= eps1; k++ {
				p, q := int(x0 + 0.5), int(y0 + 0.5)
				for j := p; j <= rightX; j++ {
					img.Set(j, q, img.GetColorFlip(j, q))
				}
				x0 += xIncre;
				y0 += yIncre;
			}
		}
	}
}

//获取相反色
func (img *FImages) GetColorFlip(x, y int) color.Color {
	r, g, b, a := img.At(x, y).RGBA()
	r_uint8 := uint8(r >> 8)
	g_uint8 := uint8(g >> 8)
	b_uint8 := uint8(b >> 8)
	a_uint8 := uint8(a >> 8)
	r_uint8 = 255 - r_uint8
	g_uint8 = 255 - g_uint8
	b_uint8 = 255 - b_uint8
	return color.RGBA{r_uint8, g_uint8, b_uint8, a_uint8}
}

func (img *FImages) DrawLine(x1, y1, x2, y2 int, c color.Color) {
	dx, dy, flag := int(math.Abs(float64(x2-x1))),
		int(math.Abs(float64(y2-y1))),
		false
	if dy > dx {
		flag = true
		x1, y1 = y1, x1
		x2, y2 = y2, x2
		dx, dy = dy, dx
	}
	ix, iy := sign(x2-x1), sign(y2-y1)
	n2dy := dy * 2
	n2dydx := (dy - dx) * 2
	d := n2dy - dx
	for x1 != x2 {
		if d < 0 {
			d += n2dy
		} else {
			y1 += iy
			d += n2dydx
		}
		if flag {
			img.Set(y1, x1, c)
		} else {
			img.Set(x1, y1, c)
		}
		x1 += ix
	}
}

func (img *FImages) SaveJPG (fName string, o *jpeg.Options) error {
	file, err := os.Create(fName + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	return jpeg.Encode(file, img, o)
}