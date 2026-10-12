package game

import chess "github.com/corentings/chess/v2"

type Game struct {
	State *chess.Game
}

func NewGame() *Game {
	return &Game{State: chess.NewGame()}
}

func NewGameFromPGN(pgnStr string) *Game {
	pgn := GetPGNFromString(pgnStr)
	return &Game{State: chess.NewGame(pgn)}
}

func (g *Game) Position() *chess.Position {
	return g.State.CurrentPosition()
}

func (g *Game) GoToMove(targetPly int) {
	for {
		currentPly := g.Position().Ply()

		if currentPly == targetPly {
			break
		}

		if currentPly > targetPly {
			g.State.GoBack()
		} else if currentPly < targetPly {
			g.State.GoForward()
		}
	}
}

func (g *Game) GoToStart() {
	for {
		cont := g.State.GoBack()
		if !cont {
			break
		}
	}
}

func (g *Game) GoToEnd() {
	for {
		cont := g.State.GoForward()
		if !cont {
			break
		}
	}
}
