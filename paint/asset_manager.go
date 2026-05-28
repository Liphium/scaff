package paint

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	// Default font face
	defaultFontFace *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Error("could not load default font", "err", err)
	}
	defaultFontFace = s
}

// TODO for this thing:
// - Cache dropping when not used for a few frames
// - Font loading

type AssetManager struct {
	fs     fs.FS
	images sync.Map // path -> *ebiten.Image
	fonts  sync.Map // path -> *text.GoTextFaceSource
}

// NewAssetManager creates a new asset manager for loading images and more based on a file system. You can, for example, use this with go:embed (which you should probably).
func NewAssetManager(fs fs.FS) *AssetManager {
	return &AssetManager{
		fs: fs,
	}
}

// GetImage loads an image from a certain path in the asset file system. It supports caching as well.
//
// Automatic cleanup of the cache is planned.
func (am *AssetManager) GetImage(path string) (*ebiten.Image, error) {
	if img, ok := am.images.Load(path); ok {
		return img.(*ebiten.Image), nil
	}

	if am.fs == nil {
		return nil, errors.New("no file system specified")
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

// GetFont loads a font from the asset file system and caches it.
//
// Automatic cleanup of the cache is planned.
func (am *AssetManager) GetFont(path string) (*text.GoTextFaceSource, error) {
	if font, ok := am.fonts.Load(path); ok {
		return font.(*text.GoTextFaceSource), nil
	}

	if am.fs == nil {
		return defaultFontFace, nil
	}
	f, err := am.fs.Open(path)
	if err != nil {
		return defaultFontFace, nil
	}
	defer f.Close()

	font, err := text.NewGoTextFaceSource(f)
	if err != nil {
		return nil, err
	}
	am.fonts.Store(path, font)
	return font, nil
}

func (am *AssetManager) Clear() {
	am.images.Range(func(key, value any) bool {
		am.images.Delete(key)
		return true
	})
}
