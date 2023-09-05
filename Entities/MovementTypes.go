package Entities

type ProjectileMovementType struct {
	Name   string
	speeds []float64
	Func   func(*Projectile)
}

var projectileMovementTypeMap = map[string]struct {
	Func              func(*Projectile)
	NumExpectedSpeeds int
}{
	"linear":      {Func: (*Projectile).linearMove, NumExpectedSpeeds: 1},      //speeed
	"osscilating": {Func: (*Projectile).osscilatingMove, NumExpectedSpeeds: 3}, //speed, amplitude, frequency
}

func NewProjectileMovementType(name string, speeds []float64) *ProjectileMovementType {
	if _, ok := projectileMovementTypeMap[name]; !ok {
		panic("Invalid movement type " + name)
	}
	if len(speeds) != projectileMovementTypeMap[name].NumExpectedSpeeds {
		panic("Invalid number of speeds for movement type " + name)
	}
	return &ProjectileMovementType{Name: name, speeds: speeds, Func: projectileMovementTypeMap[name].Func}
}

func NewLiniarMovement(speed float64) *ProjectileMovementType {
	return NewMovementTypeBuilder().
		WithName("linear").
		WithSpeeds([]float64{speed}).
		Build()
}

func NewOsscilatingMovement(speed, amplitude, frequency float64) *ProjectileMovementType {
	return NewMovementTypeBuilder().
		WithName("osscilating").
		WithSpeeds([]float64{speed, amplitude, frequency}).
		Build()
}

type EnemyMovementType struct {
	Name   string
	speeds []float64
	Func   func(*EnemyGroup)
}

var enemyMovementTypeMap = map[string]struct {
	Func              func(*EnemyGroup)
	NumExpectedSpeeds int
}{
	"Direct Follow":     {Func: (*EnemyGroup).FollowDirect, NumExpectedSpeeds: 0},
	"Orthogonal Follow": {Func: (*EnemyGroup).FollowOrth, NumExpectedSpeeds: 0},
	"Diagonal Follow":   {Func: (*EnemyGroup).FollowDiag, NumExpectedSpeeds: 0},
	"Spiral Follow":     {Func: (*EnemyGroup).FollowSpiral, NumExpectedSpeeds: 0},
}

func NewEnemyMovementType(name string, speeds []float64) *EnemyMovementType {
	if _, ok := enemyMovementTypeMap[name]; !ok {
		panic("Invalid movement type " + name)
	}
	if len(speeds) != enemyMovementTypeMap[name].NumExpectedSpeeds {
		panic("Invalid number of speeds for movement type " + name)
	}
	return &EnemyMovementType{Name: name, speeds: speeds, Func: enemyMovementTypeMap[name].Func}
}

func NewDirectFollow() *EnemyMovementType {
	return NewEnemyMovementType("Direct Follow", []float64{})
}

func NewOrthogonalFollow() *EnemyMovementType {
	return NewEnemyMovementType("Orthogonal Follow", []float64{})
}

func NewDiagonalFollow() *EnemyMovementType {
	return NewEnemyMovementType("Diagonal Follow", []float64{})
}

func NewSpiralFollow() *EnemyMovementType {
	return NewEnemyMovementType("Spiral Follow", []float64{})
}
