package Entities

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type EnemyGroupBuilder struct {
	sprite           *pixel.Sprite
	hp               float64
	damage           float64
	speed            float64
	scale            float64
	win              *pixelgl.Window
	character        *Character
	maxSpawnDistance float64
	minSpawnDistance float64
	maxDistance      float64
	movementType     EnemyMovementType
	maxSpawns        int
	name             string
}

func NewEnemyGroupBuilder() *EnemyGroupBuilder {
	return &EnemyGroupBuilder{
		maxDistance:      1000,
		maxSpawnDistance: 500,
		minSpawnDistance: 300,
	}
}

func (b *EnemyGroupBuilder) WithMovementType(movementType EnemyMovementType) *EnemyGroupBuilder {
	b.movementType = movementType
	return b
}

func (b *EnemyGroupBuilder) WithName(name string) *EnemyGroupBuilder {
	b.name = name
	return b
}

func (b *EnemyGroupBuilder) WithMaxSpawns(maxSpawns int) *EnemyGroupBuilder {
	b.maxSpawns = maxSpawns
	return b
}

func (b *EnemyGroupBuilder) WithMaxDistance(maxDistance float64) *EnemyGroupBuilder {
	b.maxDistance = maxDistance
	return b
}

func (b *EnemyGroupBuilder) WithSprite(sprite *pixel.Sprite) *EnemyGroupBuilder {
	b.sprite = sprite
	return b
}

func (b *EnemyGroupBuilder) WithHP(hp float64) *EnemyGroupBuilder {
	b.hp = hp
	return b
}

func (b *EnemyGroupBuilder) WithDamage(damage float64) *EnemyGroupBuilder {
	b.damage = damage
	return b
}

func (b *EnemyGroupBuilder) WithSpeed(speed float64) *EnemyGroupBuilder {
	b.speed = speed
	return b
}

func (b *EnemyGroupBuilder) WithScale(scale float64) *EnemyGroupBuilder {
	b.scale = scale
	return b
}

func (b *EnemyGroupBuilder) WithWin(win *pixelgl.Window) *EnemyGroupBuilder {
	b.win = win
	return b
}

func (b *EnemyGroupBuilder) WithCharacter(char *Character) *EnemyGroupBuilder {
	b.character = char
	return b
}

func (b *EnemyGroupBuilder) WithMaxSpawnDistance(maxDistance float64) *EnemyGroupBuilder {
	b.maxSpawnDistance = maxDistance
	return b
}

func (b *EnemyGroupBuilder) WithMinSpawnDistance(minDistance float64) *EnemyGroupBuilder {
	b.minSpawnDistance = minDistance
	return b
}

func (b *EnemyGroupBuilder) Build() *EnemyGroup {
	return &EnemyGroup{
		sprite:           *b.sprite,
		hp:               b.hp,
		damage:           b.damage,
		enemies:          []Enemy{},
		win:              b.win,
		character:        b.character,
		scale:            b.scale,
		speed:            b.speed,
		batch:            pixel.NewBatch(&pixel.TrianglesData{}, b.sprite.Picture()),
		hitbox:           b.sprite.Frame().W() / 4 * b.scale,
		maxSpawnDistance: b.maxSpawnDistance,
		minSpawnDistance: b.minSpawnDistance,
		maxDistance:      b.maxDistance,
		movementType:     b.movementType,
		maxSpawns:        b.maxSpawns,
		name:             b.name,
		MarkedForDeath:   false,
	}

}

type EnemyBuilder struct {
	sprite pixel.Sprite
	pos    pixel.Vec
	hp     float64
	damage float64
	scale  float64
	batch  *pixel.Batch
	hitbox float64
}

func NewEnemyBuilder() *EnemyBuilder {
	return &EnemyBuilder{}
}

func (eb *EnemyBuilder) FromGroup(group *EnemyGroup) *EnemyBuilder {
	eb.sprite = group.sprite
	eb.hp = group.hp
	eb.damage = group.damage
	eb.scale = group.scale
	eb.batch = group.batch
	eb.hitbox = group.hitbox
	return eb
}

func (eb *EnemyBuilder) WithSprite(sprite pixel.Sprite) *EnemyBuilder {
	eb.sprite = sprite
	return eb
}

func (eb *EnemyBuilder) WithPos(pos pixel.Vec) *EnemyBuilder {
	eb.pos = pos
	return eb
}

func (eb *EnemyBuilder) WithHP(hp float64) *EnemyBuilder {
	eb.hp = hp
	return eb
}

func (eb *EnemyBuilder) WithDamage(damage float64) *EnemyBuilder {
	eb.damage = damage
	return eb
}

func (eb *EnemyBuilder) WithScale(scale float64) *EnemyBuilder {
	eb.scale = scale
	return eb
}

func (eb *EnemyBuilder) WithBatch(batch *pixel.Batch) *EnemyBuilder {
	eb.batch = batch
	return eb
}

func (eb *EnemyBuilder) WithHitbox(hitbox float64) *EnemyBuilder {
	eb.hitbox = hitbox
	return eb
}

func (eb *EnemyBuilder) Build() *Enemy {
	return &Enemy{
		sprite: eb.sprite,
		pos:    eb.pos,
		hp:     eb.hp,
		damage: eb.damage,
		scale:  eb.scale,
		hitbox: eb.hitbox,
	}
}
