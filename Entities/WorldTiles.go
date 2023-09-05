package Entities

import (
	"Go_Survivor/Util"
	"math"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type WorldTilesBuilder struct {
	sprite    *pixel.Sprite
	size      pixel.Vec
	win       *pixelgl.Window
	center    pixel.Vec
	character *Character
}

func NewWorldTilesBuilder() *WorldTilesBuilder {
	return &WorldTilesBuilder{}
}

func (b *WorldTilesBuilder) WithSpriteFile(sprite string) *WorldTilesBuilder {
	b.sprite = Util.LoadSprite(sprite)
	return b
}

func (b *WorldTilesBuilder) WithSize(size pixel.Vec) *WorldTilesBuilder {
	b.size = size
	return b
}

func (b *WorldTilesBuilder) WithWin(win *pixelgl.Window) *WorldTilesBuilder {
	b.win = win
	return b
}

func (b *WorldTilesBuilder) WithCenter(center pixel.Vec) *WorldTilesBuilder {
	b.center = center
	return b
}

func (b *WorldTilesBuilder) WithCharacter(char *Character) *WorldTilesBuilder {
	b.character = char
	return b
}

func (b *WorldTilesBuilder) Build() *WorldTiles {
	return &WorldTiles{
		sprite:    *b.sprite,
		size:      b.sprite.Frame().Size(),
		win:       b.win,
		center:    b.center,
		batch:     pixel.NewBatch(&pixel.TrianglesData{}, b.sprite.Picture()),
		character: b.character,
	}
}

type WorldTiles struct {
	sprite    pixel.Sprite
	size      pixel.Vec
	win       *pixelgl.Window
	center    pixel.Vec
	character *Character
	batch     *pixel.Batch
}

func (w *WorldTiles) GetSprite() pixel.Sprite {
	return w.sprite
}

func (w *WorldTiles) GetSize() pixel.Vec {
	return w.size
}

func (w *WorldTiles) Draw() {
	offset := w.character.GetWorldPos()
	offset.X = math.Mod(offset.X, w.size.X)
	offset.Y = math.Mod(offset.Y, w.size.Y)
	for i := 0; i < 9; i++ {
		XPos := float64((i%3))*w.size.X - w.size.X
		YPos := float64((i/3))*w.size.Y - w.size.Y
		w.sprite.Draw(w.batch, pixel.IM.Moved(w.center).Moved(pixel.Vec{X: XPos, Y: YPos}).Moved(offset.Rotated(math.Pi)))
	}
	w.batch.Draw(w.win)
	w.batch.Clear()
}

func (w *WorldTiles) DrawWithOffset(offset pixel.Vec) {
	//mod the x and y pos by the size of the tile
}

func (w *WorldTiles) IsKinnematic() bool {
	return false
}

func (w *WorldTiles) GetWorldPos() pixel.Vec {
	return pixel.Vec{}
}

func (w *WorldTiles) IsColliding() bool {
	return false
}
