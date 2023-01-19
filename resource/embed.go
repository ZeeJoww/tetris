package resource

import (
	_ "embed"
)

var (
	//go:embed blocks.png
	Blocks_png []byte
	//go:embed gameover.png
	Gameover_png []byte
)
