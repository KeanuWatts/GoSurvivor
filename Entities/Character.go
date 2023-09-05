package Entities

import (
	"Go_Survivor/Util"
	"math/rand"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

type Character struct {
	sprite            pixel.Sprite
	maxHP             int
	scale             float64
	hp                float64
	worldPos          pixel.Vec
	win               *pixelgl.Window
	iFrame            time.Time
	lastMoveDir       pixel.Vec
	screenSpaceCenter pixel.Vec
	hitbox            float64
	projectiles       []ProjectileGroup
	enemies           []EnemyGroup
	speed             float64
	attackSpeed       float64
	damage            float64
	score             int
	regen             float64
	lastRegenTime     time.Time
	upgrades          []Upgrade
	AvailableUpgrades []*Upgrade
}

func (c *Character) RemoveEnemiesWithNam(name string) {
	for i := len(c.enemies) - 1; i >= 0; i-- {
		if Contains(c.enemies[i].GetName(), name) {
			c.enemies = append(c.enemies[:i], c.enemies[i+1:]...)
		}
	}
}

func Contains(basestr, searchstr string) bool {
	//seee if th searchstr exsists in the basestr
	//if it does return true
	if len(basestr) < len(searchstr) {
		return false
	}
	if basestr == searchstr {
		return true
	}
	for i := 0; i < len(basestr)-len(searchstr); i++ {
		if basestr[i:i+len(searchstr)] == searchstr {
			return true
		}
	}
	return false
}

func (c *Character) AddEnemyGroup(group EnemyGroup) {
	c.enemies = append(c.enemies, group)
}

func (c *Character) AddProjectileGroup(group ProjectileGroup) {
	c.projectiles = append(c.projectiles, group)
}

func (c *Character) Draw() {
	for i := range c.enemies {
		c.enemies[i].Draw()
	}
	for i := range c.projectiles {
		c.projectiles[i].Draw()
	}
	c.sprite.Draw(c.win, pixel.IM.Moved(c.screenSpaceCenter).Scaled(c.screenSpaceCenter, c.scale))

	//XP bar
	start, end := GetPrevAndNextTriangleNumber(c.score)
	darkBlue := colornames.Darkblue
	darkBlue.A = 128
	barCenter := pixel.Vec{X: c.screenSpaceCenter.X, Y: 5}
	barSize := pixel.Vec{X: c.screenSpaceCenter.X * 2, Y: 10}
	Util.DrawHollowRect(c.win, barCenter, barSize, darkBlue, 1, true)

	lightBlue := colornames.Lightblue
	lightBlue.A = 200
	offset := float64(c.score - start)
	max := float64(end - start)
	barCenter = pixel.Vec{X: c.screenSpaceCenter.X * offset / max, Y: 5}
	barSize = pixel.Vec{X: c.screenSpaceCenter.X * offset / max * 2, Y: 10}
	Util.DrawHollowRect(c.win, barCenter, barSize, lightBlue, 1, true)

	//HP bar
	white := colornames.White
	barSize = pixel.Vec{X: c.screenSpaceCenter.X / 2 * (float64(c.maxHP) / 100), Y: 20}
	barCenter = pixel.Vec{X: (c.screenSpaceCenter.X / 2 * (float64(c.maxHP) / 100) / 2) + 20, Y: c.screenSpaceCenter.Y*2 - 30}
	Util.DrawHollowRect(c.win, barCenter, barSize, white, 1, true)
	red := colornames.Red
	barSize = pixel.Vec{X: c.screenSpaceCenter.X / 2 * ((c.hp / 100) / (float64(c.maxHP) / 100)), Y: 20}
	barCenter = pixel.Vec{X: (c.screenSpaceCenter.X / 2 * (c.hp / float64(c.maxHP)) / 2) + 20, Y: c.screenSpaceCenter.Y*2 - 30}
	Util.DrawHollowRect(c.win, barCenter, barSize, red, 1, true)

}

func GetPrevAndNextTriangleNumber(n int) (int, int) {
	if n == 1 {
		return 0, 1
	}
	prev := 0
	next := 1
	for i := 0; i < n-1; i++ {
		if n >= prev && n < next {
			return prev, next
		}
		temp := next
		next = prev + next
		prev = temp
	}
	return prev, next
}

func (c *Character) HandleCollisions() bool {
	fullEnemyList := []*Enemy{}
	for i := 0; i < len(c.enemies); i++ {
		fullEnemyList = append(fullEnemyList, c.enemies[i].GetEnemiesPointer()...)
	}
	HandleEnemyCollisionsGeneric(fullEnemyList)

	for i := 0; i < len(c.projectiles); i++ {
		for j := 0; j < len(c.enemies); j++ {
			c.projectiles[i].HandleCollisions(fullEnemyList)
		}
	}

	if len(c.enemies) == 0 {
		return false
	}

	if !c.iFrame.Before(time.Now()) {
		return false
	}
	for i := 0; i < len(c.enemies); i++ {
		for j := 0; j < len(c.enemies[i].enemies); j++ {
			dir := c.worldPos.To(c.enemies[i].enemies[j].GetWorldPos())
			CollisionRange := c.hitbox + c.enemies[i].enemies[j].hitbox
			if dir.Len() < CollisionRange {
				c.hp -= c.enemies[i].damage
				c.iFrame = time.Now().Add(time.Millisecond * 500)
				c.enemies[i].enemies[j].Move(dir.Unit().Scaled(dir.Len() - CollisionRange))
				return true
			}
		}
	}
	return false
}

func (c *Character) GetWorldPos() pixel.Vec {
	return c.worldPos
}

func (c *Character) GetSprite() pixel.Sprite {
	return c.sprite
}

func (c *Character) GetScale() float64 {
	return c.scale
}

func (c *Character) SetPos(pos pixel.Vec) {
	c.worldPos = pos
}

func (c *Character) GetRegen() float64 {
	return c.regen
}

func (c *Character) SetRegen(regen float64) {
	c.regen = regen
}

func (c *Character) IncrementScore(inc int) {
	c.score += inc
}

func (c *Character) Move(dir pixel.Vec) {
	for i := range c.projectiles {
		c.projectiles[i].Move()
	}
	for i := range c.enemies {
		// c.enemies[i].Move(dir.Scaled(c.speed).Scaled(-1))
		c.enemies[i].Follow()
	}
	c.worldPos = c.worldPos.Add(dir.Scaled(c.speed))
	if dir.Len() > 0 {
		c.lastMoveDir = dir.Unit()
	}
}

func (c *Character) GetHP() float64 {
	return c.hp
}

func (c *Character) GetLastMoveDir() pixel.Vec {
	return c.lastMoveDir
}

func (c *Character) Heal(gain float64) {
	c.hp += gain
	if c.hp > float64(c.maxHP) {
		c.hp = float64(c.maxHP)
	}
}

func (c *Character) GetHitbox() float64 {
	return c.hitbox
}

func (c *Character) GetAttackSpeed() float64 {
	return c.attackSpeed
}

func (c *Character) SetScore(score int) {
	c.score = score
}

func (c *Character) GetScore() int {
	return c.score
}

func (c *Character) GamestateUpdate() {
	enemyCount := 0
	for i := range c.enemies {
		enemyCount += len(c.enemies[i].enemies)
	}
	if enemyCount < 100 {
		rand.Seed(time.Now().UnixNano())
		spawnRate := determineSpawnRate(c.score)
		if rand.Float64()*1000 <= spawnRate {
			i := rand.Intn(len(c.enemies))
			c.enemies[i].Spawn()
		}
	}
	for i := range c.projectiles {
		c.projectiles[i].TryShoot(c)
	}
	if time.Since(c.lastRegenTime).Seconds() > 1 {
		c.Heal(c.regen)
		c.lastRegenTime = time.Now()
	}
}

func (c *Character) ProlongLifeSpan() {
	for i := range c.projectiles {
		c.projectiles[i].ProlongLifeSpan()
	}
	c.lastRegenTime = time.Now()
}

func determineSpawnRate(score int) float64 {
	//at 0 score return 6
	//at 100 score return 600
	//over 100 score return 600
	//all in between is a linear interpolation
	if score <= 100 {
		return float64(score)*6 + 6
	}
	return 600
}

func (c *Character) RemoveDead() int {
	count := 0
	for i := len(c.enemies) - 1; i >= 0; i-- {
		count += c.enemies[i].RemoveDead()
	}
	for i := len(c.projectiles) - 1; i >= 0; i-- {
		c.projectiles[i].RemoveExpired()
	}

	return count
}

func (c *Character) GetMaxHP() int {
	return c.maxHP
}

func (c *Character) GetSpeed() float64 {
	return c.speed
}

func (c *Character) GetDamage() float64 {
	return c.damage
}

func (c *Character) SetDamage(damage float64) {
	c.damage = damage
}

func (c *Character) SetMaxHP(maxHP int) {
	// fmt.Println("Setting max hp to: ", maxHP)
	c.maxHP = maxHP
}

func (c *Character) SetHP(hp float64) {
	c.hp = hp
}

func (c *Character) SetSpeed(speed float64) {
	c.speed = speed
}

func (c *Character) SetAttackSpeed(attackSpeed float64) {
	c.attackSpeed = attackSpeed
}

func (c *Character) SetScale(scale float64) {
	c.scale = scale
}

func (c *Character) GetEnemiesPointer() *[]EnemyGroup {
	return &c.enemies
}

func (c *Character) GetProjectilesPointer() *[]ProjectileGroup {
	return &c.projectiles
}

func (c *Character) GetUpgradesPointer() *[]Upgrade {
	return &c.upgrades
}

func (c *Character) GetAvailableUpgradesPointer() []*Upgrade {
	return c.AvailableUpgrades
}

func (c *Character) AddUpgrade(upgrade *Upgrade) {
	c.upgrades = append(c.upgrades, *upgrade)
}

func (c *Character) HasUpgrade(upgrade string) bool {
	for i := range c.upgrades {
		if c.upgrades[i].GetName() == upgrade {
			return true
		}
	}
	return false
}

func (c *Character) GetAvailableUpgrades(n int) []*Upgrade {
	//We need to list through all the upgrades and find any that are not invalid
	//validation criteria:
	// get the Count from the upgrade, we must not have equal or more of that upgrade in our list
	// get the required upgrade, we must have that upgrade in our list
	upgrades := []*Upgrade{}
	selectedUpgrades := []*Upgrade{}
	for i := range c.AvailableUpgrades {
		name := c.AvailableUpgrades[i].GetName()
		count := c.AvailableUpgrades[i].GetCount()
		//count how many of this upgrade we have
		upgradeCount := 0
		for j := range c.upgrades {
			if c.upgrades[j].GetName() == name {
				upgradeCount++
			}
		}
		hasPrerequisit := c.AvailableUpgrades[i].HasPrerequisits(c)
		if upgradeCount >= count && hasPrerequisit {
			// fmt.Println("Adding upgrade: ", c.AvailableUpgrades[i].GetName())
			upgrades = append(upgrades, c.AvailableUpgrades[i])
		} else {
			// fmt.Println("Not adding upgrade: %v upgrade count: %v hasPrerequisit: %v", c.AvailableUpgrades[i].GetName(), upgradeCount, hasPrerequisit)
		}
	}
	// fmt.Println("Available Upgrades: ", upgrades)
	selectedIndexes := Util.IntList{}
	for i := 0; i < n; i++ {
		if len(upgrades) == 0 {
			break
		}
		index := rand.Intn(len(upgrades))
		for selectedIndexes.Contains(index) {
			index = rand.Intn(len(upgrades))
		}
		selectedIndexes = append(selectedIndexes, index)
		selectedUpgrades = append(selectedUpgrades, upgrades[index])
	}

	return selectedUpgrades
}
