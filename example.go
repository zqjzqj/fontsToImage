package main

import (
	"fonts2Image/fImages"
	"fonts2Image/font"
	"image/color"
	"log"
	"path"
)

func main() {
	f2p, err := font.NewFont2Points(`./0b355749.ttx`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("total fonts :", f2p.GetFontLength())
	f2p.MappingALLFontAsync(color.RGBA{0,0,0, 0}, color.RGBA{255,255,255,255}, func(img *fImages.FImages, glyph *font.TTGlyph) {
		img.SaveJPG(path.Join("./img", glyph.Name), nil)
	})

	/*//获取单个字体
	if tt, ok := f2p.GetPoint("unie5e8"); ok {
		img := tt.MappingFont(color.RGBA{0,0,0,0}, color.RGBA{255,255,255,255})
		img.SaveJPG("5e8", nil)
	}*/
}
