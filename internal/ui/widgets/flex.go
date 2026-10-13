package ui

import (
	"image/color"

	"github.com/ebitenui/ebitenui/image"
	"github.com/ebitenui/ebitenui/widget"
)

type FlexContainer struct {
	*widget.Container
}

func NewFlexContainer(
	col bool,
	px, py, spacing int,
	bg color.Color,
	opts ...widget.ContainerOpt,
) *FlexContainer {
	var direction widget.Direction
	if col {
		direction = widget.DirectionVertical
	} else {
		direction = widget.DirectionHorizontal
	}

	finalOpts := append([]widget.ContainerOpt{
		widget.ContainerOpts.BackgroundImage(
			image.NewNineSliceColor(bg),
		),
		widget.ContainerOpts.Layout(
			widget.NewRowLayout(
				widget.RowLayoutOpts.Direction(
					direction,
				),
				widget.RowLayoutOpts.Padding(
					&widget.Insets{
						Top:    py,
						Left:   px,
						Right:  px,
						Bottom: py,
					},
				),
				widget.RowLayoutOpts.Spacing(spacing),
			),
		),
	}, opts...)
	return &FlexContainer{
		Container: widget.NewContainer(finalOpts...),
	}
}
