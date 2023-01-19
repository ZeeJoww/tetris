package blocks

import (
	"github.com/hajimehoshi/ebiten/v2"
)

const (
	BoardHeight = 16
	BoardWidth  = 10
)

type Board struct {
	fixedBlocks [][]BlockType
	flushCount  int
}

func GetNewBoard() Board {
	board := Board{fixedBlocks: make([][]BlockType, BoardHeight)}
	line := make([]BlockType, BoardWidth)
	for j := 0; j < BoardWidth; j++ {
		line[j] = BlockTypeNone
	}
	for i := 0; i < BoardHeight; i++ {
		board.fixedBlocks[i] = make([]BlockType, BoardWidth)
		copy(board.fixedBlocks[i], line)
	}
	return board
}

func validCoord(x, y int) bool {
	return 0 <= x && x < BoardWidth && 0 <= y && y < BoardHeight
}

func (board Board) IsValid(block Block, x, y int) bool {
	blockShape := block.GetShape()
	n := len(blockShape)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if blockShape[i][j] {
				if !validCoord(x+j, y+i) || board.fixedBlocks[y+i][x+j] != BlockTypeNone {
					return false
				}
			}
		}
	}
	return true
}

func (board Board) Put(block Block, x, y int) (int, bool) {
	blockShape := block.GetShape()
	n := len(blockShape)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if blockShape[i][j] {
				board.fixedBlocks[y+i][x+j] = block.blockType
			}
		}
	}
	res := 0
	quicker := false
	for i := 0; i < n; i++ {
		if board.Flushable(y + i) {
			board.FlushLine(y + i)
			board.flushCount += 1
			res += 1
			if board.flushCount%5 == 0 {
				quicker = true
			}
		}
	}
	return res * (res + 1) * 5, quicker
}

func (board Board) Flushable(line int) bool {
	if line >= BoardHeight {
		return false
	}
	for j := 0; j < BoardWidth; j++ {
		if board.fixedBlocks[line][j] == BlockTypeNone {
			return false
		}
	}
	return true
}

func (board Board) FlushLine(line int) {
	for row := line; row > 0; row-- {
		copy(board.fixedBlocks[row], board.fixedBlocks[row-1])
	}
	for j := 0; j < BoardWidth; j++ {
		board.fixedBlocks[0][j] = BlockTypeNone
	}

}

func (board Board) DrawBlocks(r *ebiten.Image, block Block, x, y int) {
	for i, row := range block.GetShape() {
		for j, b := range row {
			if b {
				DrawSingleBlock(r, block.blockType, x+j, y+i)
			}
		}
	}
	for i, row := range board.fixedBlocks {
		for j, b := range row {
			if b != BlockTypeNone {
				DrawSingleBlock(r, b, j, i)
			}
		}
	}
}
