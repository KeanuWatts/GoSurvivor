package Entities

import (
	"math"
	"time"

	"github.com/faiface/pixel"
)

type Projectile struct {
	sprite       pixel.Sprite
	scale        float64
	dir          pixel.Vec
	pos          pixel.Vec
	Rotation     float64
	movementType ProjectileMovementType
	damage       float64
	lifeSpan     int64
	spawnTime    time.Time
	char         *Character
	LastMoveTime time.Time
	hitbox       float64
}

func (p *Projectile) DrawAsBatch(batch *pixel.Batch) {
	pos := pixel.IM.Scaled(pixel.ZV, p.scale).Moved(p.pos.Sub(p.char.GetWorldPos()).Add(p.char.screenSpaceCenter))
	p.sprite.Draw(batch, pos)
	//using the character, draw a line from the character to the projectile
}

func (p *Projectile) IsExpired() bool {
	return time.Since(p.spawnTime).Milliseconds() > p.lifeSpan
}

func (p *Projectile) GetSprite() pixel.Sprite {
	return p.sprite
}

func (p *Projectile) Move() {
	p.movementType.Func(p)
}

func (p *Projectile) linearMove() {
	p.pos = p.pos.Add(p.dir.Scaled(p.movementType.speeds[0]))
}

// 0 = speed
// 1 = amplitude
// 2 = frequency
func (p *Projectile) osscilatingMove() {
	SinValue := p.movementType.speeds[1] * math.Sin(p.movementType.speeds[2]*time.Since(p.spawnTime).Seconds()+float64(time.Now().Second()))

	//move the sineOffset in the direction of the projectile
	SinOffset := p.dir.Unit().Scaled(SinValue).Rotated(math.Pi / 2)
	//move the offset half its length backwards so that the projectile is centered on the sine wave
	SinOffset = SinOffset.Add(SinOffset.Scaled(-SinOffset.Len() / 2))
	p.pos = p.pos.Add(p.dir.Unit().Scaled(p.movementType.speeds[0]).Add(SinOffset))
}

func (p *Projectile) GetPos() pixel.Vec {
	return p.pos
}

func (p *Projectile) HandleCollisions(e []*Enemy) {
	for i := range e {
		dir := p.pos.To(e[i].GetWorldPos())
		CollisionRange := p.hitbox + e[i].hitbox
		if dir.Len() < CollisionRange {
			e[i].TakeDamage(p.damage)
			return
		}
	}
}
