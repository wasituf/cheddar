package ui

import (
	"bytes"
	"fmt"
	"image/color"
	"log"

	chess "github.com/corentings/chess/v2"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/sagelyone/ebitsvg"
	"github.com/wasituf/cheddar/assets"
	"github.com/wasituf/cheddar/internal/game"
)

type MoveAnimation struct {
	Piece    chess.Piece
	From, To chess.Square
	Progress float64 // 0-1
}

type Board struct {
	*widget.Widget
	game        func() *game.Game
	face        text.Face
	pieceSVGs   map[string]*ebitsvg.SVG
	pieceImages [64]*ebitsvg.Image
	pieceKeys   [64]string
	animation   *MoveAnimation
}

func NewBoard(
	currentGame func() *game.Game,
	face text.Face,
	opts ...widget.WidgetOpt,
) *Board {
	return &Board{
		Widget:    widget.NewWidget(opts...),
		game:      currentGame,
		face:      face,
		pieceSVGs: make(map[string]*ebitsvg.SVG),
	}
}

func (b *Board) GetWidget() *widget.Widget {
	return b.Widget
}

func (b *Board) PreferredSize() (int, int) {
	return 0, 0
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
			sq := chess.Square((7-row)*8 + col)
			piece := position.Piece(sq)
			if piece == chess.NoPiece {
				continue
			}
			if b.animation != nil && sq == b.animation.To {
				continue
			}

			fileName := fmt.Sprintf(
				"pieces/tatiana/%s%s.svg",
				pieceTypeMap[piece.Type()],
				piece.Color().String(),
			)
			square := row*8 + col

			pieceX := float64(x + size*0)
			pieceY := float64(y + size*0)
			pieceSize := float64(size * 1)

			b.drawPiece(screen, square, fileName, pieceX, pieceY, pieceSize)
		}
	}

	if a := b.animation; a != nil {
		// File: a=0 ... h=7. Screen row: rank 8=0 ... rank 1=7.
		fromCol := int(a.From) % 8
		fromRow := 7 - int(a.From)/8
		toCol := int(a.To) % 8
		toRow := 7 - int(a.To)/8

		fromX := float64(x0) + float64(fromCol)*float64(size)
		fromY := float64(y0) + float64(fromRow)*float64(size)
		toX := float64(x0) + float64(toCol)*float64(size)
		toY := float64(y0) + float64(toRow)*float64(size)

		p := a.Progress
		t := EaseInOut(p, 5)
		x := fromX + (toX-fromX)*t
		y := fromY + (toY-fromY)*t

		fileName := fmt.Sprintf(
			"pieces/tatiana/%s%s.svg",
			pieceTypeMap[a.Piece.Type()],
			a.Piece.Color().String(),
		)

		// Use the destination's screen index as the image-cache slot.
		// That square is skipped in the normal rendering loop.
		square := toRow*8 + toCol

		b.drawPiece(screen, square, fileName, x, y, float64(size))
	}
}

func (b *Board) drawPiece(screen *ebiten.Image, square int, fileName string, x, y, size float64) {
	if b.pieceKeys[square] != fileName {
		svg, ok := b.pieceSVGs[fileName]
		if !ok {
			f, err := assets.PiecesFS.ReadFile(fileName)
			if err != nil {
				log.Fatalf("failed to read piece asset %q: %v", fileName, err)
			}
			svg, err = ebitsvg.Parse(bytes.NewReader(f))
			if err != nil {
				log.Fatalf("failed to parse piece asset %q: %v", fileName, err)
			}
			b.pieceSVGs[fileName] = svg
		}
		b.pieceImages[square] = ebitsvg.NewImage(svg)
		b.pieceKeys[square] = fileName
	}
	pieceImg := b.pieceImages[square]

	pieceImg.Prepare(
		x,
		y,
		size, size,
		nil,
	)

	pieceImg.Draw(
		screen,
		x,
		y,
		size,
		size,
		nil,
	)
}

func (b *Board) AddAnimation(a *MoveAnimation) {
	b.animation = a
}

func (b *Board) UpdateAnimation() {
	if b.animation == nil {
		return
	}

	const durationSeconds = 0.125
	step := 1.0 / (durationSeconds * float64(ebiten.TPS()))

	b.animation.Progress += step

	if b.animation.Progress >= 1 {
		b.animation = nil
	}
}

var pieceTypeMap = map[chess.PieceType]string{
	chess.King:        "K",
	chess.Queen:       "Q",
	chess.Rook:        "R",
	chess.Bishop:      "B",
	chess.Knight:      "N",
	chess.Pawn:        "P",
	chess.NoPieceType: "",
}
