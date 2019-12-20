package main

import (
	"fonts2Image/font"
	"image/color"
	"log"
)

func main() {
	f2p, err := font.NewFont2Points(`D:\dev\DzCrawl\fonts\ttx\dishname.ttx`)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("total fonts :", f2p.GetFontLength())
	/*f2p.MappingALLFont(color.RGBA{0,0,0, 0}, color.RGBA{255,255,255,255}, func(img *fImages.FImages, glyph *font.TTGlyph) error {
		go func(img *fImages.FImages, glyph *font.TTGlyph) {
			log.Println("max x:", img.Bounds().Max.X)
			img.SaveJPG(path.Join("./img", glyph.Name), nil)
			log.Println(glyph.Name, " is OK")
		}(img, glyph)

	})*/

		if tt, ok := f2p.GetPoint("unie426"); ok {
			img := tt.MappingFont(color.RGBA{0,0,0,0}, color.RGBA{255,255,255,255})
			img.SaveJPG("5e8", nil)
		}
}
