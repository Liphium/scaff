package main

import (
	"time"

	"github.com/Liphium/scaff"
	"github.com/hajimehoshi/ebiten/v2"
)

var scaleFactor = scaff.NewSignal[float64](0)

var _ scaff.WorldLayer = (*scaleLayer)(nil)

type scaleLayer struct{}

func (rl *scaleLayer) Update(_ *scaff.LayerContext) error {
	return nil
}

func (rl *scaleLayer) Draw(c *scaff.LayerContext, _ *scaff.Camera, _ *ebiten.Image) {
	const cycleDuration = 5 * time.Second

	cyclePosition := float64(c.Now.UnixNano()%int64(cycleDuration)) / float64(cycleDuration)
	cyclePosition *= 2
	if cyclePosition > 1 {
		cyclePosition = (2 - cyclePosition) / 2
	} else {
		cyclePosition /= 2
	}
	scaleFactor.Set(0.5 + cyclePosition*4)
}

func (rl *scaleLayer) Load() {}

func (rl *scaleLayer) Unload() {}
