package cvnode

import (
	"image"
	"math"

	"github.com/Liphium/scaff"
	"github.com/Liphium/scaff/paint"
	"github.com/Liphium/scaff/scaffcv"
	"github.com/Liphium/scaff/scath"
	"github.com/hajimehoshi/ebiten/v2"
)

type TilePosition struct {
	X, Y int
}

func TilePos(x int, y int) TilePosition {
	return TilePosition{X: x, Y: y}
}

type TilemapStore struct {
	// Offset from the center
	Offset scath.Vec

	// Width for the tile
	TileWidth int

	// Height for the tile
	TileHeight int

	tilePositions map[TilePosition]int
	changedTiles  []TilePosition
}

func NewTilemapStore() *TilemapStore {
	return &TilemapStore{
		tilePositions: map[TilePosition]int{},
		changedTiles:  []TilePosition{},
	}
}

// Set a tile at a specific position on the tilemap (you can use -1 to clear if nothing should be there)
func (t *TilemapStore) Tile(tile TilePosition, tileId int) {
	t.changedTiles = append(t.changedTiles, tile)
	t.tilePositions[tile] = tileId
}

func (t *TilemapStore) clearChanged() {
	t.changedTiles = []TilePosition{}
}

// Get the tile at a specific position (-1 if it's not there)
func (t *TilemapStore) TileAt(tile TilePosition) int {
	if tile, ok := t.tilePositions[tile]; ok {
		return tile
	}
	return -1
}

func (t TilemapStore) WorldToTile(pos scath.Vec) TilePosition {
	relative := pos.Sub(t.Offset)
	return TilePosition{
		X: int(relative.X) / t.TileWidth,
		Y: int(relative.Y) / t.TileHeight,
	}
}

func (t TilemapStore) TileToWorldCenter(pos TilePosition) scath.Vec {
	return scath.Vec{
		X: float64(pos.X*t.TileWidth) + float64(t.TileWidth)/2 + t.Offset.X,
		Y: float64(pos.Y*t.TileHeight) + float64(t.TileHeight)/2 + t.Offset.Y,
	}
}

type TilemapProps struct {

	// Link to the tileset
	Tileset string

	// All tiles available (in the image, based on tile width and height)
	Tiles []TilePosition

	// The size of the cached chunks of the tilemap
	ChunkSize int

	*TilemapStore
	scaffcv.AcceptNoChild
}

// Get the top left position (chunk position) for a
func (t TilemapProps) ChunkAt(tile TilePosition) TilePosition {
	return TilePosition{
		X: tile.X - (tile.X % t.ChunkSize),
		Y: tile.Y - (tile.Y % t.ChunkSize),
	}
}

func Tilemap(create func(t *scaff.Tracker, props *TilemapProps)) scaffcv.NodeBuilder {
	// Top left pixel position -> image for the chunk
	chunks := map[TilePosition]*ebiten.Image{}

	return scaffcv.Standard(scaffcv.StandardCreate[TilemapProps]{
		ID: "tilemap",
		DefaultProps: TilemapProps{
			ChunkSize: 10,
			TilemapStore: &TilemapStore{
				TileWidth:     16,
				TileHeight:    16,
				tilePositions: map[TilePosition]int{},
				changedTiles:  []TilePosition{},
			},
		},
		PropsCreator: create,
		Create: func(methods *scaffcv.StandardMethods[TilemapProps]) {
			methods.Position = func(node *scaffcv.StandardNode[TilemapProps]) scath.Vec {
				width := node.Props().TileWidth * node.Props().ChunkSize
				height := node.Props().TileHeight * node.Props().ChunkSize

				minX, minY := 0.0, 0.0
				for pos := range chunks {
					realPos := node.Props().Offset.Add(scath.Vec{
						X: float64(pos.X * width),
						Y: float64(pos.Y * height),
					})

					minX = math.Min(minX, realPos.X)
					minY = math.Min(minY, realPos.Y)
				}
				return scath.Vec{X: minX, Y: minY}
			}
			methods.Size = func(node *scaffcv.StandardNode[TilemapProps]) scath.Vec {
				width := node.Props().TileWidth * node.Props().ChunkSize
				height := node.Props().TileHeight * node.Props().ChunkSize

				maxX, maxY := 0.0, 0.0
				for pos := range chunks {
					realPos := node.Props().Offset.Add(scath.Vec{
						X: float64(pos.X * width),
						Y: float64(pos.Y * height),
					})

					maxX = math.Max(maxX, realPos.X+float64(width))
					maxY = math.Max(maxY, realPos.Y+float64(height))
				}
				return scath.Vec{X: maxX, Y: maxY}
			}

			methods.OnPropsChanged = func(node *scaffcv.StandardNode[TilemapProps]) {

				// Find chunks to redraw and create on demand
				chunksToRedraw := map[TilePosition]struct{}{}
				for _, position := range node.Props().changedTiles {
					chunk := node.Props().ChunkAt(position)

					if chunks[chunk] == nil {
						chunks[chunk] = ebiten.NewImage(node.Props().ChunkSize*node.Props().TileWidth, node.Props().ChunkSize*node.Props().TileHeight)
					}

					chunksToRedraw[chunk] = struct{}{}
				}

				// Get the tileset
				tileset, err := node.Context().AssetManager().GetImage(node.Props().Tileset)
				if err != nil {
					log.Error("couldn't load tileset image", "path", node.Props().Tileset, "error", err)
					return
				}

				// Redraw all the chunks
				for chunk := range chunksToRedraw {
					chunkImage := chunks[chunk]
					chunkImage.Clear()

					for x := 0; x < node.Props().ChunkSize; x++ {
						for y := 0; y < node.Props().ChunkSize; y++ {
							tilePos := TilePosition{
								X: chunk.X + x,
								Y: chunk.Y + y,
							}
							tile, ok := node.Props().tilePositions[tilePos]
							if ok && tile != -1 {
								imageTile := node.Props().Tiles[tile]
								tileImage := tileset.SubImage(image.Rectangle{
									Min: image.Pt(imageTile.X, imageTile.Y),
									Max: image.Pt(imageTile.X+node.Props().TileWidth, imageTile.Y+node.Props().TileHeight),
								})

								op := &ebiten.DrawImageOptions{
									Filter: ebiten.FilterPixelated,
								}
								op.GeoM.Translate(float64(x*node.Props().TileWidth), float64(y*node.Props().TileHeight))
								chunkImage.DrawImage(tileImage.(*ebiten.Image), op)
							}
						}
					}
				}

				node.Props().TilemapStore.clearChanged()
			}

			methods.OnDraw = func(node *scaffcv.StandardNode[TilemapProps], c *scaff.Context, painter paint.Painter) {
				for chunkPos, chunk := range chunks {
					op := &ebiten.DrawImageOptions{
						Filter: ebiten.FilterPixelated,
					}
					width := node.Props().TileWidth
					height := node.Props().TileHeight

					position := node.Props().Offset.Add(scath.Vec{
						X: float64(chunkPos.X * width),
						Y: float64(chunkPos.Y * height),
					})
					op.GeoM.Translate(position.X, position.Y)

					painter.DrawRaw(chunk, op)
				}
			}
		},
	})
}
