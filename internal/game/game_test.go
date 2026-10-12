package game_test

import (
	"testing"

	"github.com/wasituf/cheddar/internal/game"
)

func TestGame_GoToStart(t *testing.T) {
	tests := []struct {
		name string
		pgn  string
	}{
		{
			name: "moves game back to starting position",
			pgn:  `[Event "Test Event"] 1. Nf3 c6 2. c4 d5 3. Nd4 dxc4`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGameFromPGN(tt.pgn)
			g.GoToStart()

			const startingFEN = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
			if got := g.Position().String(); got != startingFEN {
				t.Errorf("GoToStart() FEN = %q, want %q", got, startingFEN)
			}
		})
	}
}

func TestGame_GoToEnd(t *testing.T) {
	tests := []struct {
		name string
		pgn  string
	}{
		{
			name: "moves game forward to end position",
			pgn:  `[Event "Test Event"] 1. Nf3 c6 2. c4 d5 3. Nd4 dxc4`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGameFromPGN(tt.pgn)
			g.State.GoBack()
			g.State.GoBack()
			g.GoToEnd()

			const endFEN = "rnbqkbnr/pp2pppp/2p5/8/2pN4/8/PP1PPPPP/RNBQKB1R w KQkq - 0 4"
			if got := g.Position().String(); got != endFEN {
				t.Errorf("GoToEnd() FEN = %q, want %q", got, endFEN)
			}
		})
	}
}

func TestGame_GoToMove(t *testing.T) {
	tests := []struct {
		name      string // description of this test case
		pgn       string
		startBack int
		targetPly int
		targetFEN string
	}{
		{
			name:      "moves game to a given position higher than the current ply",
			pgn:       `[Event "Test Event"] 1. Nf3 c6 2. c4 d5 3. Nd4 dxc4`,
			startBack: 2,
			targetPly: 7,
			targetFEN: "rnbqkbnr/pp2pppp/2p5/8/2pN4/8/PP1PPPPP/RNBQKB1R w KQkq - 0 4",
		},
		{
			name:      "moves game to a given position lower than the current ply",
			pgn:       `[Event "Test Event"] 1. Nf3 c6 2. c4 d5 3. Nd4 dxc4`,
			startBack: 0,
			targetPly: 2,
			targetFEN: "rnbqkbnr/pppppppp/8/8/8/5N2/PPPPPPPP/RNBQKB1R b KQkq - 1 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := game.NewGameFromPGN(tt.pgn)
			for range tt.startBack {
				g.State.GoBack()
			}
			g.GoToMove(tt.targetPly)

			if got := g.Position().String(); got != tt.targetFEN {
				t.Errorf("GoToMove() FEN = %q, want %q", got, tt.targetFEN)
			}
		})
	}
}
