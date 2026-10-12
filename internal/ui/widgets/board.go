package ui

import (
	"image/color"
	"strings"

	chess "github.com/corentings/chess/v2"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/wasituf/cheddar/internal/game"
)

// Board displays the current position with White at the bottom.
type Board struct {
	*widget.Widget
	game func() *game.Game
	face text.Face
}

// NewBoard reads the game on each render, including after a game is replaced.
func NewBoard(currentGame func() *game.Game, face text.Face, opts ...widget.WidgetOpt) *Board {
	return &Board{
		Widget: widget.NewWidget(opts...),
		game:   currentGame,
		face:   face,
	}
}

func (b *Board) GetWidget() *widget.Widget {
	return b.Widget
}

func (b *Board) PreferredSize() (int, int) {
	return 500, 500
}

func (b *Board) Validate() {}

func (b *Board) Render(screen *ebiten.Image) {
	b.Widget.Render(screen)

	size := float32(min(b.Rect.Dx(), b.Rect.Dy())) / 8
	if size <= 0 {
		return
	}
	x0 := float32(b.Rect.Min.X) + (float32(b.Rect.Dx())-size*8)/2
	y0 := float32(b.Rect.Min.Y) + (float32(b.Rect.Dy())-size*8)/2
	position := b.game().Position().Board()

	for row := range 8 {
		for col := range 8 {
			x := x0 + float32(col)*size
			y := y0 + float32(row)*size
			squareColor := color.RGBA{240, 217, 181, 255}
			if (row+col)%2 != 0 {
				squareColor = color.RGBA{181, 136, 99, 255}
			}
			vector.FillRect(screen, x, y, size, size, squareColor, false)

			// Chess square indices start at a1; screen rows start at rank 8.
			piece := position.Piece(chess.Square((7-row)*8 + col))
			if piece == chess.NoPiece {
				continue
			}

			pieceColor, labelColor := color.Black, color.White
			if piece.Color() == chess.White {
				pieceColor, labelColor = color.White, color.Black
			}
			vector.FillCircle(screen, x+size/2, y+size/2, size*0.35, pieceColor, true)
			opts := &text.DrawOptions{}
			opts.PrimaryAlign = text.AlignCenter
			opts.SecondaryAlign = text.AlignCenter
			opts.GeoM.Translate(float64(x+size/2), float64(y+size/2))
			opts.ColorScale.ScaleWithColor(labelColor)
			text.Draw(screen, strings.ToUpper(piece.Type().String()), b.face, opts)
		}
	}
}
