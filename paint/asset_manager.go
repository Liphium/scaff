package paint

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// TODO for this thing:
// - Cache dropping when not used for a few frames
// - Font loading

type AssetManager struct {
	fs     fs.FS
	images sync.Map // path -> *ebiten.Image
}

// NewAssetManager creates a new asset manager for loading images and more based on a file system. You can, for example, use this with go:embed (which you should probably).
func NewAssetManager(fs fs.FS) *AssetManager {
	return &AssetManager{
		fs: fs,
	}
}

func (am *AssetManager) GetImage(path string) (*ebiten.Image, error) {
	if img, ok := am.images.Load(path); ok {
		return img.(*ebiten.Image), nil
	}

	f, err := am.fs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	eImg := ebiten.NewImageFromImage(img)
	am.images.Store(path, eImg)
	return eImg, nil
}

func (am *AssetManager) Clear() {
	am.images.Range(func(key, value interface{}) bool {
		am.images.Delete(key)
		return true
	})
}
