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
	*image.RGBA
}

func NewFImages(rect image.Rectangle) *FImages {
	return &FImages{image.NewRGBA(rect)}
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

//边界填充
//eColor  边界颜色
func (img *FImages) EdgeFill(xMax, xMin int, eColor color.Color) {
	//获取边界四个点
	rect := img.Bounds()
	for x := xMin; x <= xMax; x++ {
		for y := rect.Min.Y; y <= rect.Max.Y; y = y + 5 {
			pColor := img.At(x, y)
			//遇到边界
			if ColorIsEq(pColor, eColor) {
				if ColorIsEq(img.At(x + 1, y), eColor) {
					continue
				}
				//扫描线往右边填充
				n := x + 1
			//	TO:
				for x1 := n; x1 <= xMax; x1++ {
					//下一个对称点
					if ColorIsEq(img.At(x1, y), eColor) && ( !ColorIsEq(img.At(x1 - 1, y), eColor) && !ColorIsEq(img.At(x1 + 1, y), eColor) ) {
						for cx := x + 1; cx <= xMax; cx++ {
							img.Set(cx, y, img.GetColorFlip(cx, y))
						}
						for cx := x1; cx <= xMax; cx++ {
							img.Set(cx, y, img.GetColorFlip(cx, y))
						}
						break
					}
				}
				//return
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