package bots

import (
	"crypto/rand"
	"math/big"
	santorini "santorini/pkg"
)

// RandomSelector will play the game randomly
type RandomSelector struct {
	Team      int
	EnemyTeam int
	Board     *santorini.Board
}

func NewRandomBot(team int, board *santorini.Board) RandomSelector {
	enemyTeam := 1
	if team == 1 {
		enemyTeam = 2
	}

	return RandomSelector{
		Team:      team,
		EnemyTeam: enemyTeam,
		Board:     board,
	}
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

	// Always block a win if possible
	enemyCandidates := r.Board.GetValidTurns(r.EnemyTeam)

	for index, candidate := range candidates {
		// Always take a victory turn
		if candidate.IsVictory() {
			return &candidate
		}

		// Always block a win
		for _, enemyCandidate := range enemyCandidates {
			if enemyCandidate.IsVictory() {
				if enemyCandidate.MoveTo.GetX() == candidate.Build.GetX() && enemyCandidate.MoveTo.GetY() == candidate.Build.GetY() {
					return &candidates[index]
				}
			}
		}
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(len(candidates)-1)))
	if err != nil {
		panic(err)
	}

	return &candidates[n.Int64()]
}
