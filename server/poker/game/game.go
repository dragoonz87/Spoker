package game

import (
	"slices"
	"spoker/assert"
	"spoker/poker/card"
	"spoker/poker/player"
)

type PlayerId int

type Game struct {
	deck    card.Deck
	players map[PlayerId]*player.Player
	cards   []card.Card
}

func New() Game {
	deck := card.NewDeck()
	deck.Shuffle()
	return Game{deck: deck, players: make(map[PlayerId]*player.Player), cards: make([]card.Card, 5)}
}

func (g *Game) AddPlayer(id PlayerId) bool {
	if len(g.players) == 8 {
		return false
	}

	if _, ok := g.players[id]; ok {
		return false
	}

	player := player.New()
	g.players[id] = &player
	return true
}

func (g *Game) DealHands() {
	for _, p := range g.players {
		p.GiveHand(g.deck.Draw(2))
	}
}

func (g *Game) DealFlop() {
	assert.Equals(len(g.cards), 0)

	g.cards = append(g.cards, g.deck.Draw(3)...)
}

func (g *Game) DealTurn() {
	assert.Equals(len(g.cards), 3)

	g.cards = append(g.cards, g.deck.Draw(1)...)
}

func (g *Game) DealRiver() {
	assert.Equals(len(g.cards), 4)

	g.cards = append(g.cards, g.deck.Draw(1)...)
}

func (g *Game) Start() {
	g.deck = card.NewDeck()
	g.deck.Shuffle()
}

func (g *Game) getSortedCards(id PlayerId) []card.Card {
	player, ok := g.players[id]
	assert.True(ok)
	assert.Equals(len(player.Cards), 2)

	cards := append(g.cards, player.Cards...)

	slices.SortFunc(cards, func(c1, c2 card.Card) int {
		if c1.Eq(c2) {
			return 0
		} else if c1.Gt(c2) {
			return 1
		} else {
			return -1
		}
	})
	slices.Reverse(cards)

	return cards
}

func (g *Game) FindPair(sortedCards []card.Card) []card.Card {
	cards := slices.Clone(sortedCards)
	for len(cards) > 1 {
		c := cards[0]
		cards = cards[1:]
		numMatches := 0
		var match card.Card

		for _, next := range cards {
			if c.Eq(next) {
				numMatches++
				assert.True(numMatches <= 1)
				match = next
			}
		}

		if numMatches == 1 {
			return []card.Card{c, match}
		}
	}

	return nil
}

func (g *Game) FindTriple(sortedCards []card.Card) []card.Card {
	cards := slices.Clone(sortedCards)

	for len(cards) > 2 {
		c := cards[0]
		cards = cards[1:]
		numMatches := 0
		matches := make([]card.Card, 3)
        matches[0] = c

		for _, next := range cards {
			if c.Eq(next) {
				numMatches++
                assert.True(numMatches <= 2)
				matches[numMatches] = next
			}

            if numMatches == 3 {
                return matches
            }
		}
	}

	return nil
}

func (g *Game) FindQuad(sortedCards []card.Card) []card.Card {
	cards := slices.Clone(sortedCards)

	for len(cards) > 3 {
		c := cards[0]
		cards = cards[1:]
		numMatches := 0
		matches := make([]card.Card, 4)
        matches[0] = c

		for _, next := range cards {
			if c.Eq(next) {
				numMatches++
                assert.True(numMatches <= 3)
				matches[numMatches] = next
			}
		}

        if numMatches == 3 {
            return matches
        }
	}

	return nil
}
