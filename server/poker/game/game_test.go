package game

import (
	"spoker/poker/card"
	"spoker/poker/player"
	"testing"
)

func TestPair_PocketPair(t *testing.T) {
	game := Game{
		deck:    card.NewDeck(),
		players: make(map[PlayerId]*player.Player),
		cards: []card.Card{
			card.NewCard(card.CLUB, 2),
			card.NewCard(card.SPADE, 3),
			card.NewCard(card.HEART, 5),
			card.NewCard(card.DIAMOND, 8),
			card.NewCard(card.CLUB, 9),
		}}
	if !game.AddPlayer(1) {
		t.Fatal("did not add player")
	}
	p, ok := game.players[1]
	if !ok {
		t.Fatal("could not get player")
	}

	club := card.NewCard(card.CLUB, card.ACE)
	spade := card.NewCard(card.SPADE, card.ACE)
	p.GiveHand([]card.Card{
		club,
		spade,
	})

	allCards := game.getSortedCards(1)
	pair := game.FindPair(allCards)

	if pair == nil {
		t.Fatal("pair was unexpectedly nil")
	}

	if pair[0] == club {
		if pair[1] != spade {
			t.Fatalf("second card was bad. expected %v, got %v", spade, pair[1])
		}
	} else if pair[0] == spade {
		if pair[1] != club {
			t.Fatalf("second card was bad. expected %v, got %v", club, pair[1])
		}
	} else {
		t.Fatalf("first card is %v", pair[0])
	}
}

func TestPair_TablePair(t *testing.T) {
	club := card.NewCard(card.CLUB, 2)
	spade := card.NewCard(card.SPADE, 2)
	game := Game{
		deck:    card.NewDeck(),
		players: make(map[PlayerId]*player.Player),
		cards: []card.Card{
			club,
			spade,
			card.NewCard(card.HEART, 5),
			card.NewCard(card.DIAMOND, 8),
			card.NewCard(card.CLUB, 9),
		}}
	if !game.AddPlayer(1) {
		t.Fatal("did not add player")
	}
	p, ok := game.players[1]
	if !ok {
		t.Fatal("could not get player")
	}

	p.GiveHand([]card.Card{
		card.NewCard(card.CLUB, 3),
		card.NewCard(card.SPADE, 6),
	})

	allCards := game.getSortedCards(1)
	pair := game.FindPair(allCards)

	if pair == nil {
		t.Fatal("pair was unexpectedly nil")
	}

	if pair[0] == club {
		if pair[1] != spade {
			t.Fatalf("second card was bad. expected %v, got %v", spade, pair[1])
		}
	} else if pair[0] == spade {
		if pair[1] != club {
			t.Fatalf("second card was bad. expected %v, got %v", club, pair[1])
		}
	} else {
		t.Fatalf("first card is %v", pair[0])
	}
}

func TestPair_PlayerPlusTable(t *testing.T) {
	club := card.NewCard(card.CLUB, 2)
	spade := card.NewCard(card.SPADE, 2)
	game := Game{
		deck:    card.NewDeck(),
		players: make(map[PlayerId]*player.Player),
		cards: []card.Card{
			club,
			card.NewCard(card.SPADE, 6),
			card.NewCard(card.HEART, 5),
			card.NewCard(card.DIAMOND, 8),
			card.NewCard(card.CLUB, 9),
		}}
	if !game.AddPlayer(1) {
		t.Fatal("did not add player")
	}
	p, ok := game.players[1]
	if !ok {
		t.Fatal("could not get player")
	}

	p.GiveHand([]card.Card{
		spade,
		card.NewCard(card.CLUB, 3),
	})

	allCards := game.getSortedCards(1)
	pair := game.FindPair(allCards)

	if pair == nil {
		t.Fatal("pair was unexpectedly nil")
	}

	if pair[0] == club {
		if pair[1] != spade {
			t.Fatalf("second card was bad. expected %v, got %v", spade, pair[1])
		}
	} else if pair[0] == spade {
		if pair[1] != club {
			t.Fatalf("second card was bad. expected %v, got %v", club, pair[1])
		}
	} else {
		t.Fatalf("first card is %v", pair[0])
	}
}

func TestPair_NoPair(t *testing.T) {
	game := Game{
		deck:    card.NewDeck(),
		players: make(map[PlayerId]*player.Player),
		cards: []card.Card{
			card.NewCard(card.CLUB, 2),
			card.NewCard(card.SPADE, 6),
			card.NewCard(card.HEART, 5),
			card.NewCard(card.DIAMOND, 8),
			card.NewCard(card.CLUB, 9),
		}}
	if !game.AddPlayer(1) {
		t.Fatal("did not add player")
	}
	p, ok := game.players[1]
	if !ok {
		t.Fatal("could not get player")
	}

	p.GiveHand([]card.Card{
        card.NewCard(card.DIAMOND, card.KING),
		card.NewCard(card.CLUB, 3),
	})

	allCards := game.getSortedCards(1)
	pair := game.FindPair(allCards)

	if pair != nil {
		t.Fatal("pair was unexpectedly not nil")
	}
}
