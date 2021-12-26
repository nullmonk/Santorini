package main

import (
	"crypto/rand"
	"math/big"
	santorini "santorini/pkg"
)

// RandomSelector will play the game randomly
type RandomSelector struct {
	Board   *santorini.Board
	Workers []santorini.Worker
}

// SelectTurn at random, returns nil if no move can be made
func (r RandomSelector) SelectTurn() *santorini.Turn {
	var candidates []santorini.Turn
	for _, worker := range r.Workers {
		candidates = append(candidates, worker.GetValidTurns(*r.Board)...)
	}

	if candidates == nil {
		return nil
	}
	if len(candidates) == 1 {
		return &candidates[0]
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates)-1)))
	if err != nil {
		panic(err)
	}

	return &candidates[n.Int64()]
}
