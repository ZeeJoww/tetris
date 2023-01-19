package main

import (
	"log"
	"os"

	"github.com/ZeeJoww/tetris/blocks"
	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 520
	screenHeight = 860
	sampleRate   = 48000
)

type Game struct {
	x, y                   int
	score                  int
	isPressingKeyArrowDown bool
	gameover               bool
	speed                  float64
	frameCount             float64
	block                  blocks.Block
	board                  blocks.Board
	audioContext           *audio.Context
	audioPlayer            *audio.Player
}

func (g *Game) resetBlock() {
	g.block = blocks.GetNewBlock()
	g.x = 4
	g.y = 0
}

func (g *Game) initGame() {
	g.score = 0
	g.speed = 1.0
	g.frameCount = 0.0
	g.gameover = false
	g.board = blocks.GetNewBoard()
	g.resetBlock()
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.isPressingKeyArrowDown = true
	} else if inpututil.IsKeyJustReleased(ebiten.KeyArrowDown) {
		g.isPressingKeyArrowDown = false
	}

	if g.gameover {
		if g.audioPlayer.IsPlaying() {
			g.audioPlayer.Pause()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
			g.initGame()
		}
		return nil
	}

	if !g.audioPlayer.IsPlaying() {
		g.audioPlayer.Rewind()
		g.audioPlayer.Play()
	}

	// for testing
	// if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
	// 	g.block = blocks.GetNewBlock()
	// }

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		if g.board.IsValid(blocks.RotateCW(g.block), g.x, g.y) {
			g.block = blocks.RotateCW(g.block)
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		if g.board.IsValid(g.block, g.x-1, g.y) {
			g.x -= 1
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		if g.board.IsValid(g.block, g.x+1, g.y) {
			g.x += 1
		}
	}

	if g.isPressingKeyArrowDown {
		g.frameCount += 9 * g.speed
	}
	g.frameCount += g.speed

	cut := 30.0
	for g.frameCount > cut {
		g.frameCount -= cut
		g.moveDown()
	}

	return nil
}

func (g *Game) moveDown() {
	if g.board.IsValid(g.block, g.x, g.y+1) {
		g.y += 1
	} else {
		add, speedup := g.board.Put(g.block, g.x, g.y)
		g.score += add
		if speedup {
			g.speed += 0.1
		}
		g.resetBlock()
		if !g.board.IsValid(g.block, g.x, g.y) {
			g.gameover = true
		}
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.gameover {
		blocks.DrawGameOver(screen, (screenWidth-330)/2, (screenHeight-80)/2)
	} else {
		g.board.DrawBlocks(screen, g.block, g.x, g.y)
	}
	blocks.DrawScore(screen, g.score, 10, 30)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	g, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Hello, Tetris!")
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func NewGame() (*Game, error) {
	g := &Game{
		board: blocks.GetNewBoard(),
	}
	g.initGame()

	g.audioContext = audio.NewContext(sampleRate)

	f, err := os.Open("resource/ragtime.wav")
	if err != nil {
		return nil, err
	}

	d, err := wav.DecodeWithoutResampling(f)
	if err != nil {
		return nil, err
	}

	g.audioPlayer, err = g.audioContext.NewPlayer(d)
	if err != nil {
		return nil, err
	}

	return g, nil
}
