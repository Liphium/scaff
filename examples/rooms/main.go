package rooms

import (
	"github.com/Liphium/scaff"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {

	g := scaff.NewGame()

	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
