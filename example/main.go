package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/kettek/polydecomp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func main() {
	e := &Example{}
	e.Init()
	ebiten.SetWindowSize(1280, 720)
	if err := ebiten.RunGame(e); err != nil {
		log.Fatal(err)
	}
}

type Example struct {
	poly       polydecomp.Polygon[float64]
	polys      []polydecomp.Polygon[float64]
	showDecomp bool
}

var (
	mplusNormalFont font.Face
)

func (e *Example) Init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	e.poly = polydecomp.Polygon[float64]{
		{-100, 100},
		{-100, 0},
		{100, 0},
		{100, 100},
		{80, 50},
		{40, 20},
		{50, 10},
	}

	// Make that polygon CCW.
	e.poly.CCW()

	e.polys = e.poly.Decompose()
}

func (e *Example) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		e.showDecomp = !e.showDecomp
	}
	return nil
}

func (e *Example) Draw(screen *ebiten.Image) {
	text.Draw(screen, "press 'd' to toggle decomposition", mplusNormalFont, 0, 24, color.White)
	if !e.showDecomp {
		DrawPolygon(screen, 150, 100, e.poly, color.RGBA{255, 255, 255, 255})
	} else {
		for i, p := range e.polys {
			if i == 0 {
				DrawPolygon(screen, 150, 100, p, color.RGBA{255, 0, 0, 255})
			} else if i == 1 {
				DrawPolygon(screen, 150, 100, p, color.RGBA{0, 255, 0, 255})
			} else {
				DrawPolygon(screen, 150, 100, p, color.RGBA{0, 0, 255, 255})
			}
		}
	}
}

func (e *Example) Layout(oW, oH int) (sW, sH int) {
	return oW, oH
}

func DrawPolygon(screen *ebiten.Image, x, y float64, poly polydecomp.Polygon[float64], c color.RGBA) {
	for i, p := range poly {
		if i+1 < len(poly) {
			ebitenutil.DrawLine(screen, p[0]+x, p[1]+y, poly[i+1][0]+x, poly[i+1][1]+y, c)
		} else {
			c.A = 128
			ebitenutil.DrawLine(screen, p[0]+x, p[1]+y, poly[0][0]+x, poly[0][1]+y, c)
		}
		if i == 0 {
			ebitenutil.DrawRect(screen, p[0]-2+x, p[1]-2+y, 4, 4, color.RGBA{0, 255, 0, 255})
		} else if i == len(poly)-1 {
			ebitenutil.DrawRect(screen, p[0]-2+x, p[1]-2+y, 4, 4, color.RGBA{255, 0, 0, 255})
		} else {
			ebitenutil.DrawRect(screen, p[0]-2+x, p[1]-2+y, 4, 4, color.RGBA{255, 255, 255, 255})
		}
	}
}
