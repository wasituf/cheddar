package engine

import (
	"context"
)

type Evaluator interface {
	Evaluate(ctx context.Context, fen string, opts EvalOptions) (Evaluation, error)
	ID() string
	Close() error
}

type Evaluation struct {
	Centipawns int
	IsMate     bool
	BestMove   string // UCI notation
	PV         []string
	Depth      int
}

type EvalOptions struct{}
