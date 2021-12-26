package main

import (
	"crypto/rand"
	"math/big"
	santorini "santorini/pkg"
)

// RandomSelector will play the game randomly
type RandomSelector struct {
	Team  int
	Board *santorini.Board
}

// SelectTurn at random, returns nil if no move can be made
func (r RandomSelector) SelectTurn() *santorini.Turn {
	candidates := r.Board.GetValidTurns(r.Team)
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
