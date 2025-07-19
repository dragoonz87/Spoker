package card

import (
	"math/rand/v2"
	"spoker/assert"
)

type Suit uint8
type Value uint8

const (
	SPADE Suit = iota
	CLUB
	HEART
	DIAMOND
)

const (
	JACK Value = iota + 11
	QUEEN
	KING
	ACE
)

type Card struct {
	suit  Suit
	value Value
}

func (c *Card) Eq(other Card) bool {
    return c.value == other.value
}

func (c *Card) Gt(other Card) bool {
    return c.value > other.value
}

func NewCard(suit Suit, value Value) Card {
    return Card{suit, value}
}

type Deck struct {
	cards []Card
}

func NewDeck() Deck {
	var cards []Card

	for s := range 4 {
		for v := range 13 {
			var c Card

			if v == 1 {
				c.value = ACE
			} else {
				c.value = Value(v)
			}

			c.suit = Suit(s)

			cards = append(cards, c)
		}
	}

	return Deck{cards: cards}
}

func (d *Deck) Shuffle() {
	rand.Shuffle(len(d.cards), func(i, j int) {
		d.cards[i], d.cards[j] = d.cards[j], d.cards[i]
	})
}

func (d *Deck) Draw(num uint) []Card {
	assert.NotEquals(0, num)

	var cards []Card

	for range num {
		cards = append(cards, d.pop())
	}

	return cards
}

func (d *Deck) pop() Card {
	c := d.cards[0]
	d.cards = d.cards[1:]

	return c
}
