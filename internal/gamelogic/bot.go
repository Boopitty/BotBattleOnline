package gamelogic

type Bot struct {
	Name    string `json:"name"`
	Health  int    `json:"health"`
	Power   int    `json:"power"`
	Defense int    `json:"defense"`
	Huge    bool   `json:"huge"`
	IsAlive bool   `json:"is_alive"`
}

func newBot(name string, health int, power int, defense int, huge bool) *Bot {
	return &Bot{
		Name:    name,
		Health:  health,
		Power:   power,
		Defense: defense,
		Huge:    huge,
		IsAlive: true,
	}
}

func MakeBasicBot() *Bot {
	return newBot("Basic Bot", 50, 10, 5, false)
}

func MakeSpiderBot() *Bot {
	return newBot("Spider Bot", 100, 20, 5, false)
}

func (b *Bot) TakeDamage(damage int) {
	actualDamage := damage - b.Defense
	actualDamage = max(actualDamage, 1) // Ensure that damage cannot be less than 1
	b.Health -= actualDamage
	if b.Health < 0 {
		b.Health = 0
	}
	if b.Health == 0 {
		b.IsAlive = false
	}
}
