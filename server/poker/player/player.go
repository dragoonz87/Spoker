package player

import (
	"spoker/assert"
	"spoker/poker/card"
)

type Player struct {
	Cards []card.Card
	money uint
}

func New() Player {
	return Player{Cards: make([]card.Card, 0), money: 0}
}

func (p *Player) GiveHand(cards []card.Card) {
    assert.Equals(len(cards), 2)

    p.Cards = cards
}

func (p *Player) SetMoney(amount uint) {
    p.money = amount
}
