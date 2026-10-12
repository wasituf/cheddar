package app

import (
	"bytes"
	"fmt"
	"image/color"
	"math"

	"github.com/atotto/clipboard"
	"github.com/ebitenui/ebitenui"
	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/wasituf/cheddar/internal/game"
	ui "github.com/wasituf/cheddar/internal/ui/widgets"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/gofont/goregular"
)

type App struct {
	ui       *ebitenui.UI
	game     *game.Game
	pgnInput *widget.TextInput
}

func New() *App {
	a := &App{
		game: game.NewGame(),
	}

	font := DefaultFont()
	button := widget.NewButton(
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
				Idle:  Mix(colornames.Black, colornames.White, 0.1),
				Caret: Mix(colornames.Black, colornames.White, 0.15),
			},
		),
		widget.TextInputOpts.Placeholder("Enter PGN here..."),
		widget.TextInputOpts.Image(
			&widget.TextInputImage{
				Idle: image.NewNineSliceColor(colornames.Beige),
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
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(
					widget.DirectionVertical,
				),
				widget.RowLayoutOpts.Padding(
					&widget.Insets{
						Top:    16,
						Left:   16,
						Right:  16,
						Bottom: 16,
					},
				),
				widget.RowLayoutOpts.Spacing(8),
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
			widget.WidgetOpts.MinSize(200, 48),
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
			widget.WidgetOpts.MinSize(200, 48),
		),
		widget.ButtonOpts.ClickedHandler(a.forward),
	)
	nav := widget.NewContainer(
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(
					widget.DirectionHorizontal,
				),
				widget.RowLayoutOpts.Spacing(8),
			),
		),
		widget.ContainerOpts.WidgetOpts(
			widget.WidgetOpts.LayoutData(widget.RowLayoutData{
				Position: widget.RowLayoutPositionCenter,
				Stretch:  false,
			}),
		),
	)
	nav.AddChild(backMove)
	nav.AddChild(forwardMove)
	root.AddChild(a.pgnInput)
	root.AddChild(button)
	root.AddChild(ui.NewBoard(func() *game.Game { return a.game }, font,
		widget.WidgetOpts.LayoutData(widget.RowLayoutData{
			Position: widget.RowLayoutPositionCenter,
			Stretch:  true,
		}),
	))
	root.AddChild(nav)

	a.ui = &ebitenui.UI{Container: root}

	return a
}

func (g *App) Update() error {
	if (ebiten.IsKeyPressed(ebiten.KeyControlLeft) ||
		ebiten.IsKeyPressed(ebiten.KeyControlRight)) &&
		inpututil.IsKeyJustPressed(ebiten.KeyV) {
		text, err := clipboard.ReadAll()
		if err != nil {
			return fmt.Errorf("read clipboard: %w", err)
		}
		if text != "" {
			g.pgnInput.Insert(text)
		}
	}

	g.ui.Update()
	return nil
}

func (g *App) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}

func (g *App) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func DefaultFont() text.Face {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}
	return &text.GoTextFace{
		Source: s,
		Size:   20,
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
	a.game.State.GoBack()
}

func (a *App) forward(_ *widget.ButtonClickedEventArgs) {
	a.game.State.GoForward()
}
