package font

import (
	"errors"
	"fonts2Image/fImages"
	"github.com/disintegration/imaging"
	"github.com/donnie4w/dom4g"
	"image"
	"image/color"
	"io"
	"os"
	"strconv"
	"sync"
)

const(
	Xc = 100
	Yc = 200
)

type TTGlyph struct {
	Name string
	xMin int
	xMax int
	yMin int
	yMax int
	tPoints [][]image.Point
}

func (t *TTGlyph) MappingFont(lColor color.RGBA, bColor color.RGBA) *fImages.FImages {
	img := fImages.NewFImages(image.Rect(t.xMin, t.yMin + Yc, t.xMax + Xc, t.yMax - Yc))
	img.SetBkg(bColor)//设置背景
	mappingFont(img, t.tPoints, lColor)
	return img
}

func (t *TTGlyph) MappingFontByImg(img *fImages.FImages, lColor color.RGBA) {
	mappingFont(img, t.tPoints, lColor)
}

type Font2Points struct {
	//字体
	font io.Reader
	//解析到的点
	points map[string]*TTGlyph

	fontLen int
}

func NewFont2Points(fontFilePath string) (*Font2Points, error) {
	fontFile, err := os.Open(fontFilePath)
	if err != nil {
		return nil, err
	}

	ret := &Font2Points{
		font:   fontFile,
		points: make(map[string]*TTGlyph),
		fontLen:0,
	}
	if err = ret.ParsePoints(); err != nil {
		return nil, err
	}
	return ret, nil
}

func (f2p *Font2Points) GetFontLength() int {
	return f2p.fontLen
}

//解析字体
func (f2p *Font2Points) ParsePoints() error {
	ele, err := dom4g.LoadByStream(f2p.font)
	if err != nil {
		return err
	}
	eles := ele.Node("glyf").Nodes("TTGlyph")
	f2p.fontLen = 0
	for _, em := range eles {
		f2p.fontLen += 1
		if name, ok := em.AttrValue("name"); ok {
			xMin, _ := em.AttrValue("xMin")
			xMax, _ := em.AttrValue("xMax")
			yMin, _ := em.AttrValue("yMin")
			yMax, _ := em.AttrValue("yMax")
			f2p.points[name] = &TTGlyph{}
			f2p.points[name].Name = name
			f2p.points[name].xMax, _ = strconv.Atoi(xMax)
			f2p.points[name].xMin, _ = strconv.Atoi(xMin)
			f2p.points[name].yMin, _ = strconv.Atoi(yMin)
			f2p.points[name].yMax, _ = strconv.Atoi(yMax)

			//坐标系关系 y取反
			f2p.points[name].yMax = -f2p.points[name].yMax
			f2p.points[name].yMin = -f2p.points[name].yMin
			f2p.points[name].tPoints = make([][]image.Point, 0, 8)
			for _, contour := range em.Nodes("contour") {
				p2 := make([]image.Point, 0, 20)
				for _, pt := range contour.Nodes("pt") {
					xValue, _ := pt.AttrValue("x")
					yValue, _ := pt.AttrValue("y")
					xN, _ := strconv.Atoi(xValue)
					yN, _ := strconv.Atoi(yValue)
					p2 = append(p2, image.Point{X:xN, Y:-yN})
				}
				f2p.points[name].tPoints = append(f2p.points[name].tPoints, p2)
			}
		}
	}

	return nil
}

func (f2p *Font2Points) GetPoints() map[string]*TTGlyph {
	return f2p.points
}

func (f2p *Font2Points) GetPoint(fCode string) (*TTGlyph, bool) {
	ret, ok := f2p.points[fCode]
	return ret, ok
}

//画出文件单个指定字体
func (f2p *Font2Points) MappingOneFont(fCode string, lColor color.RGBA, bColor color.RGBA) (*fImages.FImages, error) {
	if fontPoint, ok := f2p.GetPoint(fCode); ok {
		img := fontPoint.MappingFont(lColor, bColor)
		return img, nil
	}
	return nil, errors.New("没有找到对应的fCode")
}

//画出文件所有字体 同步的 可在回调函数return error来中断执行
func (f2p *Font2Points) MappingALLFont(lColor color.RGBA, bColor color.RGBA, cFunc func(img *fImages.FImages, glyph *TTGlyph) error ) error {
	for _, point := range f2p.GetPoints() {
		img := point.MappingFont(lColor, bColor)
		err := cFunc(img, point)
		if err != nil {
			return err
		}
	}
	return nil
}

//画出文件所有字体 开多个协程提升速度
func (f2p *Font2Points) MappingALLFontAsync(lColor color.RGBA, bColor color.RGBA, cFunc func(img *fImages.FImages, glyph *TTGlyph) ) {
	wait := sync.WaitGroup{}
	for _, point := range f2p.GetPoints() {
		wait.Add(1)
		go func(p *TTGlyph) {
			defer wait.Done()
			img := p.MappingFont(lColor, bColor)
			cFunc(img, p)
		}(point)
	}
	wait.Wait()
}

func mappingFont(img *fImages.FImages, points [][]image.Point, lColor color.RGBA) {
	var x1, y1, x2, y2 int
	for _, point := range points {
		pointLen := len(point)
		xMax := -999999
		for k, v := range point {
			//最后一个点与第一个点相连接闭合
			x1, y1 = v.X, v.Y
			//提取最大值
			if x1 > xMax {
				xMax = x1
			}
			if k >= pointLen - 1 {
				x2, y2 = point[0].X, point[0].Y
			} else {
				x2, y2 = point[k + 1].X, point[k + 1].Y
			}
			img.DrawLine(x1, y1, x2, y2 , lColor)
		}
		img.EdgeFill(point, xMax)
	}

	//这里因为像素填充有部分问题不太好处理 所以填充后缩略图片 ocr可以更精确的识别
	img.NRGBA = imaging.Resize(img, 40, 0, imaging.Lanczos)
}