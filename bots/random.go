package bots

import (
	"crypto/rand"
	"math/big"
	santorini "santorini/pkg"

	"github.com/sirupsen/logrus"
)

// RandomSelector will play the game randomly
type RandomSelector struct {
	Team       int
	EnemyTeams []int
	Board      *santorini.Board
}

func NewRandomBot(team int, board *santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
	eTeams := make(map[int]bool, 1)
	for _, tile := range board.Tiles {
		if tile.IsOccupied() && tile.GetTeam() != team {
			eTeams[tile.GetTeam()] = true
		}
	}
	enemyTeams := make([]int, 0, len(eTeams))
	for team, _ := range eTeams {
		enemyTeams = append(enemyTeams, team)
	}
	return &RandomSelector{
		Team:       team,
		Board:      board,
		EnemyTeams: enemyTeams,
	}
}

func (r RandomSelector) Name() string {
	return "RandomBot"
}

func (r RandomSelector) IsDeterministic() bool {
	return false
}

func (r RandomSelector) testReturn(t *santorini.Turn) *santorini.Turn {
	if t.Team != r.Team {
		panic("bad team")
	}
	return t
}

// SelectTurn at random, returns nil if no move can be made
func (r RandomSelector) SelectTurn() *santorini.Turn {
	candidates := r.Board.GetValidTurns(r.Team)

	if candidates == nil {
		return nil
	}
	if len(candidates) == 1 {
		return r.testReturn(&candidates[0])
	}

	// Always block a win if possible
	enemyCandidates := r.Board.GetValidTurns(r.EnemyTeam)

	for index, candidate := range candidates {
		// Always take a victory turn
		if candidate.IsVictory() {
			return r.testReturn(&candidate)
		}
	}

	// Always block a win if possible
	for _, team := range r.EnemyTeams {
		for _, turn := range r.Board.GetValidTurns(team) {
			if turn.IsVictory() {
				// Find a move that can block it
				for _, defense := range candidates {
					if turn.MoveTo.GetX() == defense.Build.GetX() && turn.MoveTo.GetY() == defense.Build.GetY() {
						// This is the move
						return r.testReturn(&defense)
					}
				}
			}
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

	return r.testReturn(&candidates[n.Int64()])
}
