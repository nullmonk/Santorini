package bots

import (
	"fmt"
	"santorini/santorini"

	"github.com/sirupsen/logrus"
)

const maxDepth = 1

type KyleBot struct {
	Team      uint8
	EnemyTeam uint8
	Board     santorini.Board
}

func NewKyleBot(team uint8, board santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
	enemy := uint8(2)
	if team == 2 {
		enemy = 1
	}
	return &KyleBot{
		Team:      team,
		EnemyTeam: enemy,
		Board:     board,
	}
}

func (bot KyleBot) SelectTurn(b santorini.Board) *santorini.Turn {
	bot.Board = b
	candidates := b.ValidTurns(bot.Team)
	if len(candidates) == 0 {
		return nil
	}

	var (
		maxWeight int = -1000
		bestIndex     = 0
	)

	// Always block a win if possible
	enemyCandidates := bot.Board.ValidTurns(bot.EnemyTeam)

	for index, candidate := range candidates {
		// Always take a victory turn
		if candidate.IsWinningMove() {
			return candidate
		}

		// Always block a win
		for _, enemyCandidate := range enemyCandidates {
			if enemyCandidate.IsWinningMove() {
				if enemyCandidate.MoveTo.GetX() == candidate.Build.GetX() && enemyCandidate.MoveTo.GetY() == candidate.Build.GetY() {
					return candidates[index]
				}
			}
		}

		weight := bot.getWeight(candidate)

		if weight > maxWeight {
			maxWeight = weight
			bestIndex = index
		}
	}

	if bestIndex >= len(candidates) {
		panic(fmt.Errorf("what do %d %d", bestIndex, len(candidates)))
	}
	return candidates[bestIndex]
}

func (bot KyleBot) Name() string {
	return "KyleBot"
}

func (bot KyleBot) IsDeterministic() bool {
	return true
}

func (bot KyleBot) GameOver(win bool) {}

func (bot KyleBot) getWeight(candidate *santorini.Turn) int {
	// Initialize Weight
	weight := 0

	// Prefer to move up
	weight += int(candidate.MoveTo.GetHeight()) * 20

	// Prefer to cover the most tiles
	weight += len(bot.Board.GetSurroundingTiles(candidate.MoveTo.GetX(), candidate.MoveTo.GetY()))

	// Prefer to build high if no enemies are near
	if !bot.hasNearbyEnemyWorker(candidate.Worker.GetTeam(), candidate.Build) {
		weight += int(candidate.Build.GetHeight()+1) * 10
	}

	// Don't build what you cannot reach
	if candidate.MoveTo.GetHeight() < candidate.Build.GetHeight() {
		weight -= 50
	}

	// Ponder the moves to come
	thoughtBoard := bot.Board.Clone()
	thoughtBoard.PlayTurn(candidate)

	// Prefer moves that enable us to win next turn
	futureCandidates := thoughtBoard.ValidTurns(bot.Team)
	for _, futureCandidate := range futureCandidates {
		if futureCandidate.IsWinningMove() {
			weight += 1000
		}
	}

	// Avoid moves that enable an enemy win next turn
	futureEnemyCandidates := thoughtBoard.ValidTurns(bot.EnemyTeam)
	for _, futureEnemyCandidate := range futureEnemyCandidates {
		if futureEnemyCandidate.IsWinningMove() {
			weight -= 100000
		}
	}

	return weight
}

func (bot KyleBot) hasNearbyEnemyWorker(friendly uint8, tile santorini.Tile) bool {
	surroundingTiles := bot.Board.GetSurroundingTiles(tile.GetX(), tile.GetY())
	for _, surroundingTile := range surroundingTiles {
		if team := surroundingTile.GetTeam(); team != 0 && team != friendly {
			// Enemy worker cannot navigate to the new tile if built
			if surroundingTile.GetHeight() < tile.GetHeight() {
				continue
			}

			return true
		}
	}

	return false
}
