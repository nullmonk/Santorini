package bots

import (
	"crypto/rand"
	"math/big"
	"santorini/santorini"
)

/*
Implent a bot that takes three basic actions
	1. Win if possible *
	2. Stop win if possbile *
	3. Randomly move

*implmented in the Bot class, no work for us
*/

func NewRandomBot() santorini.BotInitializer {
	return NewBasicBot("RandomBot", randomChoice, false)
}

func randomChoice(g *Bot, b santorini.Board, t *santorini.Turn) int {
	// give this move a random rank
	n, err := rand.Int(rand.Reader, big.NewInt(993939))
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}
