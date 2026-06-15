package engine

type Background interface {
	Init()

	SetTPS(tps int)

	OnUpdate(func())

	OnDraw(func(texture Texture))

	LoadTexture(path string) (Texture, error)
}
