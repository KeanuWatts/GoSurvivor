package Entities

import (
	"time"

	"github.com/faiface/pixel"
)

type Enemy struct {
	sprite        pixel.Sprite
	pos           pixel.Vec
	hp            float64
	damage        float64
	scale         float64
	hitbox        float64
	lastDirChange time.Time
	lastDir       pixel.Vec
	MovementType  ProjectileMovementType
	OffsetAngle   float64
}

func (e *Enemy) GetDamage() float64 {
	return e.damage
}

func (e *Enemy) DrawAsBatch(batch *pixel.Batch, char *Character) {
	pos := pixel.IM.Scaled(pixel.ZV, e.scale).Moved(e.pos.Sub(char.GetWorldPos()).Add(char.screenSpaceCenter))
	e.sprite.Draw(batch, pos)

}

func (e *Enemy) Move(dir pixel.Vec) {
	e.pos = e.pos.Add(dir)
}

func (e *Enemy) IsKinnematic() bool {
	return true
}

func (e *Enemy) GetWorldPos() pixel.Vec {
	return e.pos
}

func (e *Enemy) SetWorldPos(pos pixel.Vec) {
	e.pos = pos
}

func (e *Enemy) GetSprite() pixel.Sprite {
	return e.sprite
}

func (e *Enemy) TakeDamage(damage float64) {
	e.hp -= damage
}

func (e *Enemy) IsColliding() bool {
	return true
}

func (e *Enemy) Draw() {
	// batch := pixel.NewBatch(&pixel.TrianglesData{}, e.sprite.Picture())
	// e.sprite.Draw(batch, pixel.IM.Moved(e.pos))
	// batch.Clear()
}

func (e *Enemy) GetScale() float64 {
	return e.scale
}
