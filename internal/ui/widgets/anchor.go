package ui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type AnchorContainer struct {
	*widget.Container
}

func NewAnchorContainer(px, py int, bg color.Color, opts ...widget.ContainerOpt) *AnchorContainer {
	finalOpts := append([]widget.ContainerOpt{
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(bg),
		),
		widget.ContainerOpts.Layout(
			widget.NewAnchorLayout(
				widget.AnchorLayoutOpts.Padding(
					&widget.Insets{
						Top:    py,
						Left:   px,
						Right:  px,
						Bottom: py,
					},
				),
			),
		),
	}, opts...)
	return &AnchorContainer{
		Container: widget.NewContainer(finalOpts...),
	}
}
