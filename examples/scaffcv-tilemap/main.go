package main

import (
	"embed"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scaffcv/cvnode"
	"github.com/Liphium/scaff/scath"
	sutil "github.com/Liphium/scaff/util"
	"github.com/hajimehoshi/ebiten/v2"
)

// Embed the assets file system for getting the images.
//
//go:embed assets/*
var assetsFS embed.FS

var log = sutil.NewLogger("app")

func main() {
	ebiten.SetWindowSize(900, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	g := scaff.NewGame()

	assetManager := paint.NewAssetManager(assetsFS)
	tree := scaff.NewSceneTree("scaffui-image-sample", assetManager)

	tree.Mount(func(t *scaff.Tracker, props *scaff.RootProps) {
		props.Child(0, scaffcv.Canvas(func(t *scaff.Tracker, props *scaffcv.CanvasProps) {
			props.Position = scath.Zero

			props.Child(cvnode.Stack(func(t *scaff.Tracker, props *cvnode.StackProps) {
				props.Child(0, cvnode.Tilemap(func(t *scaff.Tracker, props *cvnode.TilemapProps) {
					props.Tileset = "assets/microslop.png"

					props.TileWidth = 16
					props.TileHeight = 16

					// Register all the tiles
					const (
						TileRed = iota
						TileBlue
						TileGreen
						TileYellow
					)
					props.Tiles = []cvnode.TilePosition{
						cvnode.TilePos(0, 0),   // Red
						cvnode.TilePos(16, 0),  // Blue
						cvnode.TilePos(0, 16),  // Green
						cvnode.TilePos(16, 16), // Yellow
					}

					props.Tile(cvnode.TilePos(0, 0), TileRed)
					props.Tile(cvnode.TilePos(1, 0), TileBlue)
					props.Tile(cvnode.TilePos(0, 1), TileGreen)
					props.Tile(cvnode.TilePos(1, 1), TileYellow)
				}))
			}))
		}))
	})

	g.Goto(tree)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
