package card

import "testing"

func TestShuffle(t *testing.T) {
	d1 := NewDeck()
	d2 := NewDeck()

	d1.Shuffle()

	identical := true
	for i := range 52 {
		if d1.cards[i] != d2.cards[i] {
			identical = false
			break
		}
	}

	if identical {
		t.Fatal("shuffle did not shuffle the cards")
	}
}
