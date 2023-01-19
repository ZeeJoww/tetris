package blocks

import (
	"bytes"
	"image"
	_ "image/png"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	rb "github.com/ZeeJoww/tetris/resource"
)

const (
	BlockSideLen = 50
)

type BlockType int

const (
	BlockType0 BlockType = iota // L
	BlockType1                  // J
	BlockType2                  // Z
	BlockType3                  // S
	BlockType4                  // T
	BlockType5                  // #
	BlockType6                  // -
	BlockTypeUb
	BlockTypeNone BlockType = -1
)

type BlockAngle int

const (
	Angle0 BlockAngle = iota
	Angle90
	Angle180
	Angle270
)

type Block struct {
	blockType  BlockType
	blockAngle BlockAngle
}

func RotateCW(block Block) Block {
	newblock := block
	newblock.blockAngle = block.blockAngle.rotateCW()
	return newblock
}

func GetNewBlock() Block {
	return Block{
		blockType:  BlockType(rand.Intn(int(BlockTypeUb))),
		blockAngle: BlockAngle(rand.Intn(4)),
	}
}

func (angle BlockAngle) rotateCW() BlockAngle {
	return (angle + 3) % 4
}

var imageBlocks, imageGameover *ebiten.Image

func init() {
	rand.Seed(time.Now().Unix())
	img, _, err := image.Decode(bytes.NewReader(rb.Blocks_png))
	if err != nil {
		panic(err)
	}
	imageBlocks = ebiten.NewImageFromImage(img)
	img, _, err = image.Decode(bytes.NewReader(rb.Gameover_png))
	if err != nil {
		panic(err)
	}
	imageGameover = ebiten.NewImageFromImage(img)
}

var getShape map[BlockType]map[BlockAngle][][]bool

func (block Block) GetShape() [][]bool {
	return getShape[block.blockType][block.blockAngle]
}

func copyMatrix(mat [][]bool) [][]bool {
	res := make([][]bool, len(mat))
	for i := range res {
		res[i] = make([]bool, len(mat[i]))
		copy(res[i], mat[i])
	}
	return res
}

func rotate90(mat [][]bool) [][]bool {
	res := copyMatrix(mat)
	n := len(mat)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			res[i][j] = mat[j][n-i-1]
		}
	}
	var left, top int
	for j := 0; j < n; j++ {
		for i := 0; i < n; i++ {
			if res[i][j] {
				left = j
				j = n
				break
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if res[i][j] {
				top = i
				i = n
				break
			}
		}
	}
	ret := copyMatrix(res)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			ret[i][j] = false
		}
	}
	for i := top; i < n; i++ {
		for j := left; j < n; j++ {
			ret[i-top][j-left] = res[i][j]
		}
	}
	return ret
}

func init() {
	const (
		T = true
		F = false
	)
	var BaseShape = map[BlockType][][]bool{
		BlockType0: {
			{F, F, T},
			{T, T, T},
			{F, F, F},
		}, // L
		BlockType1: {
			{T, F, F},
			{T, T, T},
			{F, F, F},
		}, // J
		BlockType2: {
			{T, T, F},
			{F, T, T},
			{F, F, F},
		}, // Z
		BlockType3: {
			{F, T, T},
			{T, T, F},
			{F, F, F},
		}, // S
		BlockType4: {
			{T, T, T},
			{F, T, F},
			{F, F, F},
		}, // T
		BlockType5: {
			{T, T},
			{T, T},
		}, // #
		BlockType6: {
			{T, T, T, T},
			{F, F, F, F},
			{F, F, F, F},
			{F, F, F, F},
		}, // -
	}

	getShape = map[BlockType]map[BlockAngle][][]bool{}
	for block := BlockType0; block < BlockTypeUb; block++ {
		getShape[block] = map[BlockAngle][][]bool{}
		for _, angle := range []BlockAngle{Angle0, Angle90, Angle180, Angle270} {
			getShape[block][angle] = copyMatrix(BaseShape[block])
			BaseShape[block] = rotate90(BaseShape[block])
		}
	}
}
