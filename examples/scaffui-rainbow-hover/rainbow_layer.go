package main

import (
	"image/color"
	"math"
	"time"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/scaffui"
	"github.com/hajimehoshi/ebiten/v2"
)

var rainbow = scaffui.NewSignal(color.RGBA{255, 255, 255, 255})

var _ scaff.WorldLayer = (*rainbowLayer)(nil)

type rainbowLayer struct{}

func (rl *rainbowLayer) Update(_ *scaff.LayerContext) error {
	return nil
}

func (rl *rainbowLayer) Draw(c *scaff.LayerContext, _ *scaff.Camera, _ *ebiten.Image) {
	const cycleDuration = 10 * time.Second

	cyclePosition := float64(c.Now.UnixNano()%int64(cycleDuration)) / float64(cycleDuration)
	rainbow.Set(hsvToRGBA(cyclePosition, 1, 1))
}

func (rl *rainbowLayer) Load() {}

func (rl *rainbowLayer) Unload() {}

func hsvToRGBA(h, s, v float64) color.RGBA {
	h = math.Mod(h, 1)
	if h < 0 {
		h += 1
	}

	segment := h * 6
	chroma := v * s
	x := chroma * (1 - math.Abs(math.Mod(segment, 2)-1))
	m := v - chroma

	var r, g, b float64
	switch {
	case segment < 1:
		r, g, b = chroma, x, 0
	case segment < 2:
		r, g, b = x, chroma, 0
	case segment < 3:
		r, g, b = 0, chroma, x
	case segment < 4:
		r, g, b = 0, x, chroma
	case segment < 5:
		r, g, b = x, 0, chroma
	default:
		r, g, b = chroma, 0, x
	}

	return color.RGBA{
		R: uint8((r + m) * 255),
		G: uint8((g + m) * 255),
		B: uint8((b + m) * 255),
		A: 255,
	}
}
