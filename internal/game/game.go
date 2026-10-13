package game

import (
	chess "github.com/corentings/chess/v2"
)

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

func (g *Game) GetPieceFromMove(m *chess.Move) chess.Piece {
	if m == nil || m.Parent() == nil {
		return chess.NoPiece
	}

	p := m.Parent().Position().Board().Piece(m.S1())

	return p
}

func (g *Game) GetLatestMove() *chess.Move {
	if g.State.IsAtStart() {
		return nil
	}

	moves := g.State.Moves()
	ply := g.Position().Ply() - 1

	var move *chess.Move
	for _, m := range moves {
		if m.Ply() == ply {
			move = m
		}
	}

	return move
}

func (g *Game) GetNextMove() *chess.Move {
	if g.State.IsAtEnd() {
		return nil
	}

	moves := g.State.Moves()
	ply := g.Position().Ply()

	var move *chess.Move
	for _, m := range moves {
		if m.Ply() == ply {
			move = m
		}
	}

	return move
}
