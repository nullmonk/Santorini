package bots

import (
	santorini "santorini/pkg"
)

type KyleBot struct {
	Team      int
	EnemyTeam int
	Board     *santorini.Board
}

func (bot KyleBot) SelectTurn() *santorini.Turn {
	worker1 := bot.Board.GetWorkerTile(bot.Team, 1)
	worker2 := bot.Board.GetWorkerTile(bot.Team, 2)
	enemyWorker1 := bot.Board.GetWorkerTile(bot.EnemyTeam, 1)
	enemyWorker2 := bot.Board.GetWorkerTile(bot.EnemyTeam, 2)

	candidates := bot.Board.GetValidTurns(bot.Team)

	var (
		maxWeight int = -1000
		bestIndex     = 0
	)

	// Always take a victory turn
	for index, candidate := range candidates {
		if candidate.IsVictory() {
			return &candidate
		}

		// Initialize Weight
		weight := 0

		// Weight based on height
		if candidate.Worker == 1 {
			if worker1.GetHeight() < candidate.MoveTo.GetHeight() {
				weight += 10
			} else if worker1.GetHeight() < candidate.MoveTo.GetHeight() {
				weight -= 10
			}
		} else if candidate.Worker == 2 {
			if worker2.GetHeight() < candidate.MoveTo.GetHeight() {
				weight += 20
			} else if worker2.GetHeight() < candidate.MoveTo.GetHeight() {
				weight -= 20
			}
		}

		// Avoid building near enemy workers
		var enemyReach []santorini.Tile
		enemyMovement := append(
			bot.Board.GetSurroundingTiles(enemyWorker1.GetX(), enemyWorker1.GetY()),
			bot.Board.GetSurroundingTiles(enemyWorker2.GetX(), enemyWorker2.GetY())...,
		)
		for _, tile := range enemyMovement {
			if candidate.Build.GetX() == tile.GetX() || candidate.Build.GetY() == tile.GetY() {
				weight -= 10
			}
			enemyReach = append(enemyReach, bot.Board.GetSurroundingTiles(tile.GetX(), tile.GetY())...)
		}
		for _, tile := range enemyReach {
			if candidate.Build.GetX() == tile.GetX() || candidate.Build.GetY() == tile.GetY() {
				weight -= 5
			}
		}

		// Prefer building on edge tiles
		if candidate.Build.GetX() == 0 ||
			candidate.Build.GetX() == bot.Board.Size-1 ||
			candidate.Build.GetY() == 0 ||
			candidate.Build.GetY() == bot.Board.Size-1 {
			weight += 1
		}

		// Avoid moves that enable the enemy to win
		// thoughtBoard := bot.getBoardCopy()
		// thoughtBoard.PlayTurn(candidate)

		// enemyCandidates := bot.Board.GetValidTurns(bot.EnemyTeam)
		// for _, enemyCandidate := range enemyCandidates {
		// 	if enemyCandidate.IsVictory() {
		// 		weight = -1000
		// 		break
		// 	}
		// }

		if weight > maxWeight {
			bestIndex = index
		}
	}

	return &candidates[bestIndex]
}

// func (bot KyleBot) getBoardCopy() (thoughtBoard santorini.Board) {
// 	copier.Copy(&thoughtBoard, *bot.Board)
// 	return
// }
