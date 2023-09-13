package Entities

import (
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type ProjectileGroup struct {
	projectiles   []Projectile
	sprite        pixel.Sprite
	scale         float64
	movementType  ProjectileMovementType
	damage        float64
	char          *Character
	win           *pixelgl.Window
	lifeSpan      int64
	batch         *pixel.Batch
	hitbox        float64
	lastSpawn     time.Time
	attacksPerSec float64
	name          string
}

func (pg *ProjectileGroup) SpawnProjectil(dir pixel.Vec) {
	pg.projectiles = append(pg.projectiles,
		*NewProjectileBuilder().
			FromGroup(pg).
			WithPos(pg.char.GetWorldPos()).
			WithDir(dir).
			Build())
}

func (pg *ProjectileGroup) TryShoot(char *Character) {
	if time.Since(pg.lastSpawn).Seconds() > pg.attacksPerSec/char.attackSpeed {
		pg.projectiles = append(pg.projectiles,
			*NewProjectileBuilder().
				FromGroup(pg).
				WithPos(char.worldPos).
				WithDir(char.lastMoveDir).
				Build())
		pg.lastSpawn = time.Now()
	}
}

func (pg *ProjectileGroup) ProlongLifeSpan() {
	for i := range pg.projectiles {
		pg.projectiles[i].spawnTime = pg.projectiles[i].spawnTime.Add(time.Since(pg.projectiles[i].lastMoveTime))
	}
}

func (pg *ProjectileGroup) GetProjectiles() []Projectile {
	return pg.projectiles
}

func (pg *ProjectileGroup) Draw() {
	for _, v := range pg.projectiles {
		v.DrawAsBatch(pg.batch)
	}
	pg.batch.Draw(pg.win)
	pg.batch.Clear()
}

func (pg *ProjectileGroup) Move() {
	for i := range pg.projectiles {
		pg.projectiles[i].Move()
	}
}

func (pg *ProjectileGroup) RemoveExpired() {
	if len(pg.projectiles) == 0 {
		return
	}
	for i := 0 + len(pg.projectiles) - 1; i >= 0; i-- {
		if pg.projectiles[i].IsExpired() {
			if i+1 >= len(pg.projectiles)-1 {
				pg.projectiles = pg.projectiles[:i]
			} else {
				pg.projectiles = append(pg.projectiles[:i], pg.projectiles[i+1:]...)
			}
		}
	}
}

func (pg *ProjectileGroup) HandleCollisions(e []*Enemy) {
	for i := range pg.projectiles {
		pg.projectiles[i].HandleCollisions(e)
	}
}

func (pg *ProjectileGroup) GetHitbox() float64 {
	return pg.hitbox
}

func (pg *ProjectileGroup) GetDamage() float64 {
	return pg.damage
}

func (pg *ProjectileGroup) GetName() string {
	return pg.name
}

func (pg *ProjectileGroup) GetSprite() pixel.Sprite {
	return pg.sprite
}

func (pg *ProjectileGroup) GetScale() float64 {
	return pg.scale
}

func (pg *ProjectileGroup) GetChar() *Character {
	return pg.char
}

func (pg *ProjectileGroup) GetLifeSpan() int64 {
	return pg.lifeSpan
}

func (pg *ProjectileGroup) GetMovementType() ProjectileMovementType {
	return pg.movementType
}

func (pg *ProjectileGroup) GetAttackSpeed() float64 {
	return pg.attacksPerSec
}

func (pg *ProjectileGroup) SetAttackSpeed(attackSpeed float64) {
	pg.attacksPerSec = attackSpeed
}

func (pg *ProjectileGroup) SetDamage(damage float64) {
	pg.damage = damage
}

func (pg *ProjectileGroup) SetLifeSpan(lifeSpan int64) {
	pg.lifeSpan = lifeSpan
}

func (pg *ProjectileGroup) SetMovementType(movementType ProjectileMovementType) {
	pg.movementType = movementType
}

func (pg *ProjectileGroup) SetChar(char *Character) {
	pg.char = char
}

func (pg *ProjectileGroup) SetSprite(sprite pixel.Sprite) {
	pg.sprite = sprite
}

func (pg *ProjectileGroup) SetScale(scale float64) {
	//adjust the hitbox to the new scale
	pg.hitbox *= scale / pg.scale
	pg.scale = scale
}

func (pg *ProjectileGroup) GetMovementTypeSpeed() float64 {
	return pg.movementType.speeds[0]
}

func (pg *ProjectileGroup) SetMovementTypeSpeed(speed float64) {
	pg.movementType.speeds[0] = speed
}
