package Entities

import (
	"Go_Survivor/Util"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type MovementTypeBuilder struct {
	Name   string
	speeds []float64
}

func NewMovementTypeBuilder() *MovementTypeBuilder {
	return &MovementTypeBuilder{}
}

func (b *MovementTypeBuilder) WithName(name string) *MovementTypeBuilder {
	b.Name = name
	return b
}

func (b *MovementTypeBuilder) WithSpeeds(speeds []float64) *MovementTypeBuilder {
	b.speeds = speeds
	return b
}

func (b *MovementTypeBuilder) Build() *ProjectileMovementType {
	return NewProjectileMovementType(b.Name, b.speeds)
}

type ProjectileGroupBuilder struct {
	win           *pixelgl.Window
	spriteFile    string
	scale         float64
	movement      ProjectileMovementType
	damage        float64
	char          *Character
	lifeSpan      int64
	attacksPerSec float64
	name          string
}

func NewProjectileGroupBuilder() *ProjectileGroupBuilder {
	return &ProjectileGroupBuilder{}
}

func (b *ProjectileGroupBuilder) WithWin(win *pixelgl.Window) *ProjectileGroupBuilder {
	b.win = win
	return b
}

func (b *ProjectileGroupBuilder) WithSpriteFile(sprite string) *ProjectileGroupBuilder {
	b.spriteFile = sprite
	return b
}

func (b *ProjectileGroupBuilder) WithScale(scale float64) *ProjectileGroupBuilder {
	b.scale = scale
	return b
}

func (b *ProjectileGroupBuilder) WithAttackPerSec(aps float64) *ProjectileGroupBuilder {
	b.attacksPerSec = aps
	return b
}

func (b *ProjectileGroupBuilder) WithMovementType(movement ProjectileMovementType) *ProjectileGroupBuilder {
	b.movement = movement
	return b
}

func (b *ProjectileGroupBuilder) WithDamage(damage float64) *ProjectileGroupBuilder {
	b.damage = damage
	return b
}

func (b *ProjectileGroupBuilder) WithChar(char *Character) *ProjectileGroupBuilder {
	b.char = char
	return b
}

func (b *ProjectileGroupBuilder) WithLifeSpan(lifeSpan int64) *ProjectileGroupBuilder {
	b.lifeSpan = lifeSpan
	return b
}

func (b *ProjectileGroupBuilder) WithName(name string) *ProjectileGroupBuilder {
	b.name = name
	return b
}

func (b *ProjectileGroupBuilder) Build() *ProjectileGroup {
	sprite := Util.LoadSprite(b.spriteFile)
	return &ProjectileGroup{
		projectiles:   []Projectile{},
		sprite:        *sprite,
		scale:         b.scale,
		movementType:  b.movement,
		damage:        b.damage,
		char:          b.char,
		win:           b.win,
		lifeSpan:      b.lifeSpan,
		batch:         pixel.NewBatch(&pixel.TrianglesData{}, sprite.Picture()),
		hitbox:        sprite.Frame().W() / 4 * b.scale,
		lastSpawn:     time.Now(),
		attacksPerSec: b.attacksPerSec,
		name:          b.name,
	}
}

type ProjectileBuilder struct {
	sprite       *pixel.Sprite
	scale        float64
	dir          pixel.Vec
	pos          pixel.Vec
	movementType ProjectileMovementType
	damage       float64
	lifeSpan     int64
	hitbox       float64
	char         *Character
}

func NewProjectileBuilder() *ProjectileBuilder {
	return &ProjectileBuilder{}
}

func (b *ProjectileBuilder) FromGroup(pg *ProjectileGroup) *ProjectileBuilder {
	b.sprite = &pg.sprite
	b.scale = pg.scale
	b.movementType = pg.movementType
	b.damage = pg.damage
	b.char = pg.char
	b.lifeSpan = pg.lifeSpan
	b.hitbox = pg.hitbox
	return b
}

func (b *ProjectileBuilder) WithSprite(sprite *pixel.Sprite) *ProjectileBuilder {
	b.sprite = sprite
	return b
}

func (b *ProjectileBuilder) WithScale(scale float64) *ProjectileBuilder {
	b.scale = scale
	return b
}

func (b *ProjectileBuilder) WithDir(dir pixel.Vec) *ProjectileBuilder {
	b.dir = dir
	return b
}

func (b *ProjectileBuilder) WithPos(pos pixel.Vec) *ProjectileBuilder {
	b.pos = pos
	return b
}

func (b *ProjectileBuilder) WithMovementType(movementType ProjectileMovementType) *ProjectileBuilder {
	b.movementType = movementType
	return b
}

func (b *ProjectileBuilder) WithDamage(damage float64) *ProjectileBuilder {
	b.damage = damage
	return b
}

func (b *ProjectileBuilder) WithLifeSpan(lifeSpan int64) *ProjectileBuilder {
	b.lifeSpan = lifeSpan
	return b
}

func (b *ProjectileBuilder) WithHitbox(hitbox float64) *ProjectileBuilder {
	b.hitbox = hitbox
	return b
}

func (b *ProjectileBuilder) WithChar(char *Character) *ProjectileBuilder {
	b.char = char
	return b
}

func (b *ProjectileBuilder) Build() *Projectile {
	return &Projectile{
		sprite:       *b.sprite,
		scale:        b.scale,
		dir:          b.dir,
		Rotation:     b.dir.Angle(),
		pos:          b.pos,
		movementType: b.movementType,
		damage:       b.damage,
		spawnTime:    time.Now(),
		lifeSpan:     b.lifeSpan,
		char:         b.char,
		LastMoveTime: time.Now(),
		hitbox:       b.hitbox,
	}
}
