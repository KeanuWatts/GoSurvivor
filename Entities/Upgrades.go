package Entities

type UpgradeBuilder struct {
	name         string
	description  string
	prerequisits []string
	function     func(*Character)
	count        int
}

func NewUpgradeBuilder() *UpgradeBuilder {
	return &UpgradeBuilder{}
}

func (b *UpgradeBuilder) WithName(name string) *UpgradeBuilder {
	b.name = name
	return b
}

func (b *UpgradeBuilder) WithCount(count int) *UpgradeBuilder {
	b.count = count
	return b
}

func (b *UpgradeBuilder) WithDescription(description string) *UpgradeBuilder {
	b.description = description
	return b
}

func (b *UpgradeBuilder) WithPrerequisit(prerequisit []string) *UpgradeBuilder {
	b.prerequisits = prerequisit
	return b
}

func (b *UpgradeBuilder) WithFunction(function func(*Character)) *UpgradeBuilder {
	b.function = function
	return b
}

func (b *UpgradeBuilder) Build() *Upgrade {
	return &Upgrade{
		Name:         b.name,
		Description:  b.description,
		Prerequisits: b.prerequisits,
		Func:         b.function,
	}
}

type Upgrade struct {
	Name         string
	Description  string
	Prerequisits []string
	Func         func(*Character)
	Count        int
}

func GetDefaultUpgrades(c *Character) []*Upgrade {

	// OscilatingMovement := NewMovementTypeBuilder().
	// WithName("osscilating").
	// WithSpeeds([]float64{5, 5, 8}).
	// Build()

	// LiniarMovemeent := NewMovementTypeBuilder().
	// WithName("liniar").
	// WithSpeeds([]float64{2}).
	// Build()

	upgrades := []Upgrade{
		*NewUpgradeBuilder().WithName("Max HP").WithDescription("Increases your Max HP by 10").WithCount(0).WithFunction(func(ca *Character) {
			ca.SetMaxHP(int(float64(ca.GetMaxHP()) * 1.1))
		}).Build(),
		*NewUpgradeBuilder().WithName("Heal").WithDescription("Increases your HP by 10").WithCount(0).WithFunction(func(ca *Character) {
			ca.SetHP(ca.GetHP() + 10)
		}).Build(),
		*NewUpgradeBuilder().WithName("attack speed").WithDescription("Increases your attack speed by 10%").WithCount(10).WithFunction(func(ca *Character) {
			ca.SetAttackSpeed(ca.GetAttackSpeed() * 1.1)
		}).Build(),
		*NewUpgradeBuilder().WithName("damage").WithDescription("Increases your damage by 10%").WithCount(0).WithFunction(func(ca *Character) {
			ca.SetDamage(ca.GetDamage() * 1.1)
		}).Build(),
		*NewUpgradeBuilder().WithName("Move speed").WithDescription("Increases your Move speed by 10%").WithCount(10).WithFunction(func(ca *Character) {
			ca.SetSpeed(ca.GetSpeed() * 1.1)
		}).Build(),
		*NewUpgradeBuilder().WithName("Regen").WithDescription("Incraeses regen by 0.1 hp per second").WithCount(10).WithFunction(func(ca *Character) {
			ca.SetRegen(ca.GetRegen() + 0.1)
		}).Build(),
		*NewUpgradeBuilder().WithName("Basic Attack attack speed upgrade").WithDescription("Increases your basic attack attack speed by 10%").WithCount(10).WithFunction(func(ca *Character) {
			pgs := ca.GetProjectilesPointer()
			for i := range *pgs {
				if (*pgs)[i].GetName() == "Basic Attack" {
					(*pgs)[i].SetAttackSpeed((*pgs)[i].GetAttackSpeed() * 1.1)
				}
			}
		}).Build(),
		*NewUpgradeBuilder().WithName("Basic Attack damage upgrade").WithDescription("Increases your basic attack damage by 10%").WithCount(10).WithFunction(func(ca *Character) {
			pgs := ca.GetProjectilesPointer()
			for i := range *pgs {
				if (*pgs)[i].GetName() == "Basic Attack" {
					(*pgs)[i].SetDamage((*pgs)[i].GetDamage() * 0.9)
				}
			}
		}).Build(),
		*NewUpgradeBuilder().WithName("Basic Attack range upgrade").WithDescription("Increases your basic attack life span by 10%").WithCount(10).WithFunction(func(ca *Character) {
			pgs := ca.GetProjectilesPointer()
			for i := range *pgs {
				if (*pgs)[i].GetName() == "Basic Attack" {
					(*pgs)[i].SetLifeSpan(int64(float64((*pgs)[i].GetLifeSpan()) * 1.1))
				}
			}
		}).Build(),
		*NewUpgradeBuilder().WithName("Basic Attack projectile size upgrade").WithDescription("Increases your basic attack projectile size by 10%").WithCount(10).WithFunction(func(ca *Character) {
			pgs := ca.GetProjectilesPointer()
			for i := range *pgs {
				if (*pgs)[i].GetName() == "Basic Attack" {
					(*pgs)[i].SetScale((*pgs)[i].GetScale() * 2)
				}
			}
		}).Build(),
		*NewUpgradeBuilder().WithName("Basic Attack projectile speed upgrade").WithDescription("Increases your basic attack projectile speed by 10%").WithCount(10).WithFunction(func(ca *Character) {
			pgs := ca.GetProjectilesPointer()
			for i := range *pgs {
				if (*pgs)[i].GetName() == "Basic Attack" {
					(*pgs)[i].SetMovementTypeSpeed((*pgs)[i].GetMovementTypeSpeed() * 1.1)
				}
			}
		}).Build(),
		*NewUpgradeBuilder().WithName("Type 2 Attack").WithDescription("Adds a new attack").WithCount(10).WithFunction(func(ca *Character) {
			moveType := *NewOsscilatingMovement(3, 5, 8)
			pg :=
				NewProjectileGroupBuilder().
					WithWin(ca.win).
					WithName("Type 2 Attack").
					WithChar(ca).
					WithSpriteFile("Util/Assets/Bullet.png").
					WithScale(0.2).
					WithLifeSpan(1000).
					WithDamage(10).
					WithMovementType(moveType).
					WithAttackPerSec(0.5).
					Build()

			ca.AddProjectileGroup(*pg)
		}).Build(),
	}
	result := make([]*Upgrade, len(upgrades))
	for i := range upgrades {
		result[i] = &upgrades[i]
	}
	return result
}

func (u *Upgrade) GetName() string {
	return u.Name
}

func (u *Upgrade) GetDescription() string {
	return u.Description
}

func (u *Upgrade) GetPrerequisit() []string {
	return u.Prerequisits
}

func (u *Upgrade) GetFunc() func(*Character) {
	return u.Func
}

func (u *Upgrade) GetCount() int {
	return u.Count
}

func (u *Upgrade) SetCount(count int) {
	u.Count = count
}

func (u *Upgrade) AddCount(count int) {
	u.Count += count
}

func (u *Upgrade) DecrementCount() {
	u.Count--
}

func (u *Upgrade) IncrementCount() {
	u.Count++
}

func (u *Upgrade) HasPrerequisits(c *Character) bool {
	if u.Prerequisits == nil {
		return true
	}
	for i := range u.Prerequisits {
		if !c.HasUpgrade(u.Prerequisits[i]) {
			return false
		}
	}
	return true
}

func (u *Upgrade) ApplyUpgrade(c *Character) {
	u.Func(c)
}
