package gamelogic

type targetType string

const (
	SingleEnemy targetType = "single"
	AllEnemies  targetType = "all_enemies"
	SingleAlly  targetType = "single_ally"
	AllAllies   targetType = "all_allies"
	Self        targetType = "self"
	All         targetType = "all"
	AllOthers   targetType = "all_others"
)

type Skill struct {
	name        string
	description string
	power       int
	targetType  targetType
}

type SkillIndex int

const (
	Bash   SkillIndex = 0
	Heal   SkillIndex = 1
	Blast  SkillIndex = 2
	Shield SkillIndex = 3
)

var SkillList = []Skill{
	{name: "Bash", description: "A standard melee attack.", power: 10, targetType: SingleEnemy},
	{name: "Heal", description: "Restores health to an ally.", power: 15, targetType: SingleAlly},
	{name: "Blast", description: "A blast that hits multiple enemies.", power: 20, targetType: AllEnemies},
	{name: "Shield", description: "Creates a protective barrier for self.", power: 0, targetType: Self},
}
