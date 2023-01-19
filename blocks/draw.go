package blocks

import (
	"fmt"
	"image"

	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var (
	scoreFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	scoreFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

const (
	horiLowerBound = 10
	vertLowerBound = 50
)

func DrawSingleBlock(r *ebiten.Image, block BlockType, x, y int) {
	if block == BlockTypeNone {
		return
	}

	op := &ebiten.DrawImageOptions{}
	op.ColorM = ebiten.ColorM{}
	x = x*BlockSideLen + horiLowerBound
	y = y*BlockSideLen + vertLowerBound
	op.GeoM.Translate(float64(x), float64(y))

	x0 := int(block) * BlockSideLen
	r.DrawImage(imageBlocks.SubImage(image.Rect(x0, 0, x0+BlockSideLen, BlockSideLen)).(*ebiten.Image), op)
}

func DrawGameOver(r *ebiten.Image, x, y int) {
	op := &ebiten.DrawImageOptions{}
	op.ColorM = ebiten.ColorM{}

	op.GeoM.Translate(float64(x), float64(y))
	r.DrawImage(imageGameover, op)
}

func DrawScore(r *ebiten.Image, score, x, y int) {
	text.Draw(r, fmt.Sprint("得分:", score), scoreFont, x, y, color.White)
}
