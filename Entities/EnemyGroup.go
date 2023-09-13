package Entities

import (
	"math"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

type EnemyGroup struct {
	sprite           pixel.Sprite
	hp               float64
	damage           float64
	enemies          []Enemy
	speed            float64
	win              *pixelgl.Window
	scale            float64
	character        *Character
	batch            *pixel.Batch
	hitbox           float64
	maxSpawnDistance float64
	minSpawnDistance float64
	maxDistance      float64
	movementType     EnemyMovementType
	maxSpawns        int
	name             string
	MarkedForDeath   bool
}

func (e *EnemyGroup) GetDamage() float64 {
	return e.damage
}

func (e *EnemyGroup) MarkForDeath() {
	e.MarkedForDeath = true
}

func (e *EnemyGroup) IsMarkedForDeath() bool {
	return e.MarkedForDeath
}

func (e *EnemyGroup) GetName() string {
	return e.name
}

func (e *EnemyGroup) GetCount() int {
	return len(e.enemies)
}

func (e *EnemyGroup) GetEnemies() []Enemy {
	return e.enemies
}

func (e *EnemyGroup) GetEnemiesPointer() []*Enemy {
	result := []*Enemy{}
	for i := range e.enemies {
		result = append(result, &e.enemies[i])
	}
	return result
}

func HandleEnemyCollisionsGeneric(enemis []*Enemy) {
	Collided := true
	if len(enemis) > 1 {
		for Collided {
			Collided = false
			for i := range enemis {
				for j := range enemis {
					if i != j {
						pos1 := enemis[i].GetWorldPos()
						pos2 := enemis[j].GetWorldPos()
						dir := pos1.To(pos2)
						CollisionRange := enemis[i].hitbox + enemis[j].hitbox
						if dir.Len() < CollisionRange {
							Collided = true
							RepelVector := dir.Unit().Scaled(dir.Len() - CollisionRange)
							enemis[i].Move(RepelVector)
							RepelVector = RepelVector.Scaled(-1)
							enemis[j].Move(RepelVector)
						}
					}
				}
			}
		}
	}
}

func (e *EnemyGroup) RemoveDead() int {
	count := 0
	for i := len(e.enemies) - 1; i >= 0; i-- {
		if e.enemies[i].hp <= 0 {
			count++
			if i+1 >= len(e.enemies) {
				e.enemies = e.enemies[:i]
			} else {
				e.enemies = append(e.enemies[:i], e.enemies[i+1:]...)
			}
		}
	}
	return count
}

func (e *EnemyGroup) Spawn() {
	if e.MarkedForDeath {
		return
	}
	//get a random position, 300-500 pixels away from the character, in a random direction
	if len(e.enemies) >= e.maxSpawns {
		return
	}
	rand.Seed(time.Now().UnixNano())
	randXi := int((rand.Float64()-1)*e.minSpawnDistance + e.maxSpawnDistance)
	randYi := int((rand.Float64()-1)*e.minSpawnDistance + e.maxSpawnDistance)
	if randXi%2 == 0 {
		randXi *= -1
	}
	if randYi%2 == 0 {
		randYi *= -1
	}
	randX := float64(randXi)
	randY := float64(randYi)
	pos := pixel.Vec{X: randX, Y: randY}.Add(e.character.GetWorldPos())
	enemy := NewEnemyBuilder().
		FromGroup(e).
		WithPos(pos).
		Build()

	//add the enemy to the list
	e.Add(*enemy)
}

func (e *EnemyGroup) Move(dir pixel.Vec) {
	for i := range e.enemies {
		e.enemies[i].Move(dir)
	}
}

func (e *EnemyGroup) IncreaseDifficulty(scale int) {
	//increase the hp damage and speed of the enemies, hp scales exponentially, damage scale linearly, speed scales logarithmically
	e.hp = e.hp*(float64(scale)*0.1) + 1
	e.damage = 1 + float64(scale)
	e.speed = 1 + 0.1*float64(scale)
}

func (e *EnemyGroup) Follow() {
	e.movementType.Func(e)
}

func (e *EnemyGroup) FollowDirect() {
	for i := range e.enemies {
		target := e.character.GetWorldPos()
		pos := e.enemies[i].GetWorldPos()
		//move the enemy towards the character
		direction := target.Sub(pos).Unit().Scaled(e.speed)
		e.enemies[i].Move(direction)
	}
}

func (e *EnemyGroup) FollowOrth() {
	for i := range e.enemies {
		target := e.character.GetWorldPos()
		pos := e.enemies[i].GetWorldPos()
		//move the enemy towards the character
		direction := target.Sub(pos)
		// Create a movement pattern that will move only orthogonally, and make sure they move on one axis only for at least one second before changing direction and duration is unique per enemy in the group
		if time.Since(e.enemies[i].lastDirChange).Seconds() > 1 {
			// Generate a random number between 0 and 3 to determine which direction to move in
			randNum := rand.Intn(2)
			switch randNum {
			case 0:
				direction.X = 0
			case 1:
				direction.Y = 0
			}
			e.enemies[i].lastDirChange = time.Now()
			e.enemies[i].lastDir = direction.Unit().Scaled(e.speed)
		}
		e.enemies[i].Move(e.enemies[i].lastDir)
	}
}

func (e *EnemyGroup) FollowDiag() {
	for i := range e.enemies {
		target := e.character.GetWorldPos()
		pos := e.enemies[i].GetWorldPos()
		//move the enemy towards the character
		direction := target.Sub(pos)
		// Create a movement pattern that will move only orthogonally, and make sure they move on one axis only for at least one second before changing direction and duration is unique per enemy in the group
		if time.Since(e.enemies[i].lastDirChange).Seconds() > 1 {
			// Generate a random number between 0 and 3 to determine which direction to move in
			randNum := rand.Intn(2)
			switch randNum {
			case 0:
				direction.X = direction.Y
			case 1:
				direction.X = -direction.Y
			}
			e.enemies[i].lastDirChange = time.Now()
			e.enemies[i].lastDir = direction.Unit().Scaled(e.speed)
		}
		e.enemies[i].Move(e.enemies[i].lastDir)
	}
}

func (e *EnemyGroup) FollowSpiral() {
	for i := range e.enemies {
		target := e.character.GetWorldPos()
		pos := e.enemies[i].GetWorldPos()
		if e.enemies[i].OffsetAngle == 0 {
			//randomly choose a direction to spiral in
			if rand.Intn(2) == 0 {
				e.enemies[i].OffsetAngle = math.Pi / 2.1
			} else {
				e.enemies[i].OffsetAngle = -math.Pi / 2.1
			}
		}
		direction := target.Sub(pos).Unit().Scaled(e.speed).Rotated(e.enemies[i].OffsetAngle)

		e.enemies[i].Move(direction)
	}
}

func (e *EnemyGroup) Add(enemy Enemy) {
	e.enemies = append(e.enemies, enemy)
}

func (e *EnemyGroup) Remove() int {
	startLen := len(e.enemies)
	for i, v := range e.enemies {
		if v.hp <= 0 {
			e.enemies = append(e.enemies[:i], e.enemies[i+1:]...)
		}
	}
	result := startLen - len(e.enemies)
	//if an enemy is more than 1000 pixels away from the character, remove it
	for i, v := range e.enemies {
		if v.GetWorldPos().To(e.character.GetWorldPos()).Len() > e.maxDistance {
			//check if i+1 is out of bounds
			if i+1 >= len(e.enemies) {
				e.enemies = e.enemies[:i]
			} else {
				e.enemies = append(e.enemies[:i], e.enemies[i+1:]...)
			}
		}
	}
	return result

}

func (e *EnemyGroup) Draw() {
	for _, v := range e.enemies {
		v.DrawAsBatch(e.batch, e.character)
	}
	e.batch.Draw(e.win)
	//for debugging, draw the position of the character as text
	// for _, v := range e.enemies {
	// 	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	// 	text := text.New(pixel.Vec{X: 0, Y: 0}, basicAtlas)
	// 	text.WriteString(fmt.Sprintf("Pos: %v", v.pos))
	// 	text.Draw(e.win, pixel.IM.Moved(v.pos.Sub(e.character.GetWorldPos()).Add(Character.Center)))
	// }
	e.batch.Clear()
}
