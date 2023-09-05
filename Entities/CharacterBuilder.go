package Entities

import (
	"Go_Survivor/Util"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type CharacterBuilder struct {
	sprite      *pixel.Sprite
	scale       float64
	hp          float64
	win         *pixelgl.Window
	speed       float64
	attackSpeed float64

	AvailableUpgrades []Upgrade
}

func NewCharacterBuilder() *CharacterBuilder {
	return &CharacterBuilder{}
}

func (cb *CharacterBuilder) WithSpriteFile(file string) *CharacterBuilder {
	cb.sprite = Util.LoadSprite(file)
	return cb
}

func (cb *CharacterBuilder) WithSpeed(speed float64) *CharacterBuilder {
	cb.speed = speed
	return cb
}

func (cb *CharacterBuilder) WithHP(hp float64) *CharacterBuilder {
	cb.hp = hp
	return cb
}

func (cb *CharacterBuilder) WithWin(win *pixelgl.Window) *CharacterBuilder {
	cb.win = win
	return cb
}

func (cb *CharacterBuilder) WithScale(scale float64) *CharacterBuilder {
	cb.scale = scale
	return cb
}

func (cb *CharacterBuilder) WithAttackSpeed(attackSpeed float64) *CharacterBuilder {
	cb.attackSpeed = attackSpeed
	return cb
}

func (cb *CharacterBuilder) Build() *Character {
	Center := cb.win.Bounds().Center()
	char := &Character{
		sprite:            *cb.sprite,
		hp:                cb.hp,
		maxHP:             int(cb.hp),
		worldPos:          Center,
		win:               cb.win,
		iFrame:            time.Now(),
		screenSpaceCenter: Center,
		hitbox:            cb.sprite.Frame().W() / 4 * cb.scale,
		speed:             cb.speed,
		scale:             cb.scale,
		attackSpeed:       cb.attackSpeed,
		damage:            10,
		regen:             0,
		lastRegenTime:     time.Now(),
	}
	char.AvailableUpgrades = GetDefaultUpgrades(char)
	return char
}
