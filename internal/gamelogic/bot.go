package gamelogic

type Bot struct {
	Name     string       `json:"name"`
	Value    int          `json:"value"`
	Base_hp  int          `json:"base_hp"`
	Hp       int          `json:"hp"`
	Base_atk int          `json:"base_atk"`
	Atk      int          `json:"atk"`
	Base_def int          `json:"base_def"`
	Def      int          `json:"def"`
	Base_spd int          `json:"base_spd"`
	Spd      int          `json:"spd"`
	Huge     bool         `json:"huge"`
	IsAlive  bool         `json:"is_alive"`
	Skills   []SkillIndex `json:"skills"`
}

type BotNum int

const (
	None      BotNum = -1
	BasicBot  BotNum = 0
	SpiderBot BotNum = 1
	BigBot    BotNum = 2
)

func newBot(name string, value int, hp int, atk int, def int, spd int, huge bool) *Bot {
	return &Bot{
		Name:     name,
		Value:    value,
		Base_hp:  hp,
		Hp:       hp,
		Base_atk: atk,
		Atk:      atk,
		Base_def: def,
		Def:      def,
		Base_spd: spd,
		Spd:      spd,
		Huge:     huge,
		IsAlive:  true,
		Skills:   make([]SkillIndex, 0, 4),
	}
}

func MakeBot(num BotNum) *Bot {
	switch num {
	case BasicBot:
		return newBot("Basic Bot", 10, 50, 10, 5, 10, false)
	case SpiderBot:
		return newBot("Spider Bot", 20, 100, 20, 5, 10, false)
	case BigBot:
		return newBot("Big Bot", 200, 15, 10, 5, 10, true)
	default:
		return nil
	}
}

func (b *Bot) AddSkill(skill SkillIndex) bool {
	if len(b.Skills) < 4 {
		b.Skills = append(b.Skills, skill)
		return true
	}
	return false
}

func (b *Bot) RemoveSkill(skill SkillIndex) {
	for i, s := range b.Skills {
		if s == skill {
			b.Skills = append(b.Skills[:i], b.Skills[i+1:]...)
			break
		}
	}
}

// Get skill from the global SkillList by index. Returns a pointer to the skill.
func (b *Bot) GetSkill(skill SkillIndex) *Skill {
	return &SkillList[skill]
}

func (b *Bot) TakeDamage(damage int) {
	actualDamage := damage - b.Def
	actualDamage = max(actualDamage, 1) // Ensure that damage cannot be less than 1
	b.Hp -= actualDamage
	if b.Hp < 0 {
		b.Hp = 0
	}
	if b.Hp == 0 {
		b.IsAlive = false
	}
}

func (b *Bot) Attack(target *Bot, power int) {
	if !b.IsAlive {
		return
	}
	target.TakeDamage(b.Atk * power)
}
