package app

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/wasituf/cheddar/assets"
	"github.com/wasituf/cheddar/internal/game"
	ui "github.com/wasituf/cheddar/internal/ui/widgets"
	"golang.org/x/image/colornames"
)

type App struct {
	ui       *ebitenui.UI
	game     *game.Game
	pgnInput *widget.TextInput
	board    *ui.Board
}

func New() *App {
	a := &App{
		game: game.NewGame(),
	}

	font := DefaultFont()
	importButton := widget.NewButton(
		widget.ButtonOpts.TextLabel("Import"),

		widget.ButtonOpts.TextFace(&font),
		widget.ButtonOpts.TextColor(&widget.ButtonTextColor{
			Idle:    colornames.Gainsboro,
			Hover:   colornames.Gainsboro,
			Pressed: Mix(colornames.Gainsboro, colornames.Black, 0.4),
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSliceColor(colornames.Darkslategray),
			Hover: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Mediumseagreen, 0.4),
			),
			Disabled: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Gainsboro, 0.8),
			),
			Pressed: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Black, 0.4),
			),
			PressedHover: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Black, 0.4),
			),
		}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 48,
			}),
			widget.WidgetOpts.MinSize(20, 48),
		),
		widget.ButtonOpts.ClickedHandler(a.onImport),
	)
	a.pgnInput = widget.NewTextInput(
		widget.TextInputOpts.Face(&font),
		widget.TextInputOpts.Color(
			&widget.TextInputColor{
				Idle:     Mix(colornames.Black, colornames.White, 0.1),
				Caret:    Mix(colornames.Black, colornames.White, 0.15),
				Disabled: Mix(colornames.Black, colornames.White, 0.5),
			},
		),
		widget.TextInputOpts.Placeholder("Enter PGN here..."),
		widget.TextInputOpts.Image(
			&widget.TextInputImage{
				Idle: image.NewNineSliceColor(colornames.Gainsboro),
			},
		),
		widget.TextInputOpts.Padding(
			&widget.Insets{
				Top:    8,
				Left:   16,
				Right:  16,
				Bottom: 8,
			},
		),
		widget.TextInputOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  true,
			}),
			widget.WidgetOpts.MinSize(20, 48),
		),
	)
	root := widget.NewContainer(
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(colornames.Gainsboro),
		),
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Stretch(
					[]bool{true, false},
					[]bool{true},
				),
				widget.GridLayoutOpts.Padding(
					&widget.Insets{
						Top:    16,
						Left:   16,
						Right:  16,
						Bottom: 16,
					},
				),
				widget.GridLayoutOpts.Spacing(0, 8),
			),
		),
	)

	backMove := widget.NewButton(
		widget.ButtonOpts.TextLabel("Back"),

		widget.ButtonOpts.TextFace(&font),
		widget.ButtonOpts.TextColor(&widget.ButtonTextColor{
			Idle:    colornames.Gainsboro,
			Hover:   colornames.Gainsboro,
			Pressed: Mix(colornames.Gainsboro, colornames.Black, 0.4),
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSliceColor(colornames.Darkslategray),
			Hover: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Mediumseagreen, 0.4),
			),
			Disabled: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Gainsboro, 0.8),
			),
			Pressed: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Black, 0.4),
			),
			PressedHover: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Black, 0.4),
			),
		}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 48,
			}),
			widget.WidgetOpts.MinSize(100, 48),
		),
		widget.ButtonOpts.ClickedHandler(a.back),
	)
	forwardMove := widget.NewButton(
		widget.ButtonOpts.TextLabel("Forward"),

		widget.ButtonOpts.TextFace(&font),
		widget.ButtonOpts.TextColor(&widget.ButtonTextColor{
			Idle:    colornames.Gainsboro,
			Hover:   colornames.Gainsboro,
			Pressed: Mix(colornames.Gainsboro, colornames.Black, 0.4),
		}),
		widget.ButtonOpts.Image(&widget.ButtonImage{
			Idle: image.NewNineSliceColor(colornames.Darkslategray),
			Hover: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Mediumseagreen, 0.4),
			),
			Disabled: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Gainsboro, 0.8),
			),
			Pressed: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Black, 0.4),
			),
			PressedHover: image.NewNineSliceColor(
				Mix(colornames.Darkslategray, colornames.Black, 0.4),
			),
		}),
		widget.ButtonOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position:  widget.RowLayoutPositionCenter,
				Stretch:   true,
				MaxHeight: 48,
			}),
			widget.WidgetOpts.MinSize(100, 48),
		),
		widget.ButtonOpts.ClickedHandler(a.forward),
	)
	nav := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewGridLayout(
				widget.GridLayoutOpts.Columns(2),
				widget.GridLayoutOpts.Stretch(
					[]bool{true, true}, []bool{true},
				),
				widget.GridLayoutOpts.Spacing(4, 4),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.AnchorLayoutData{
				HorizontalPosition: widget.AnchorLayoutPositionCenter,
				VerticalPosition:   widget.AnchorLayoutPositionEnd,
				StretchHorizontal:  true,
			}),
		),
	)

	board := ui.NewBoard(func() *game.Game { return a.game }, font,
		widget.WidgetOpts.LayoutData(widget.GridLayoutData{
			HorizontalPosition: widget.GridLayoutPositionCenter,
			VerticalPosition:   widget.GridLayoutPositionCenter,
		}),
	)

	nav.AddChild(backMove)
	nav.AddChild(forwardMove)

	pgnImport := ui.NewFlexContainer(true, 0, 0, 4, colornames.Lightslategrey,
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.AnchorLayoutData{
					HorizontalPosition: widget.AnchorLayoutPositionCenter,
					VerticalPosition:   widget.AnchorLayoutPositionStart,
					StretchHorizontal:  true,
				},
			),
		),
	)
	pgnImport.AddChild(a.pgnInput)
	pgnImport.AddChild(importButton)

	sideBar := ui.NewAnchorContainer(8, 8, colornames.Lightgrey,
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(
				widget.GridLayoutData{},
			),
			widget.WidgetOpts.MinSize(300, 0),
		),
	)
	sideBar.AddChild(pgnImport)
	sideBar.AddChild(nav)

	root.AddChild(board, sideBar)

	a.ui = &ebitenui.UI{Container: root}
	a.board = board

	return a
}

func (a *App) Update() error {
	if (ebiten.IsKeyPressed(ebiten.KeyControlLeft) ||
		ebiten.IsKeyPressed(ebiten.KeyControlRight)) &&
		inpututil.IsKeyJustPressed(ebiten.KeyV) {
		text, err := clipboard.ReadAll()
		if err != nil {
			return fmt.Errorf("read clipboard: %w", err)
		}
		if text != "" {
			a.pgnInput.Insert(text)
		}
	}

	a.board.UpdateAnimation()
	a.ui.Update()
	return nil
}

func (a *App) Draw(screen *ebiten.Image) {
	a.ui.Draw(screen)
}

func (a *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func DefaultFont() text.Face {
	fName := "Manrope"
	fVariant := "Medium"
	fFileName := fmt.Sprintf("fonts/%s/%s-%s.ttf", strings.ToLower(fName), fName, fVariant)

	f, err := assets.FontsFS.ReadFile(fFileName)
	if err != nil {
		log.Fatal(err)
	}

	// s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	s, err := text.NewGoTextFaceSource(bytes.NewReader(f))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: s,
		Size:   16,
	}
}

func Mix(a, b color.Color, percent float64) color.Color {
	rgba := func(c color.Color) (r, g, b, a uint8) {
		r16, g16, b16, a16 := c.RGBA()
		return uint8(r16 >> 8), uint8(g16 >> 8), uint8(b16 >> 8), uint8(a16 >> 8)
	}
	lerp := func(x, y uint8) uint8 {
		return uint8(math.Round(float64(x) + percent*(float64(y)-float64(x))))
	}
	r1, g1, b1, a1 := rgba(a)
	r2, g2, b2, a2 := rgba(b)

	return color.RGBA{
		R: lerp(r1, r2),
		G: lerp(g1, g2),
		B: lerp(b1, b2),
		A: lerp(a1, a2),
	}
}

func (a *App) onImport(_ *widget.ButtonClickedEventArgs) {
	imported := game.NewGameFromPGN(a.pgnInput.GetText())
	a.game = imported
}

func (a *App) back(_ *widget.ButtonClickedEventArgs) {
	// Animation
	currentMove := a.game.GetLatestMove()
	if currentMove == nil {
		return
	}
	a.board.AddAnimation(&ui.MoveAnimation{
		Piece:    a.game.GetPieceFromMove(currentMove),
		From:     currentMove.S2(), // Start from the current move's square
		To:       currentMove.S1(), // End at the previous square
		Progress: 0,
	})

	a.game.State.GoBack()
}

func (a *App) forward(_ *widget.ButtonClickedEventArgs) {
	nextMove := a.game.GetNextMove()
	if nextMove == nil {
		return
	}
	a.board.AddAnimation(&ui.MoveAnimation{
		Piece:    a.game.GetPieceFromMove(nextMove),
		From:     nextMove.S1(), // Start from the current move's square
		To:       nextMove.S2(), // End at the next square
		Progress: 0,
	})

	a.game.State.GoForward()
}
