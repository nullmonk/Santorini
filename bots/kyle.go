package bots

import (
	santorini "santorini/pkg"
)

type KyleBot struct {
	Team      int
	EnemyTeam int
	Board     *santorini.Board
}

func NewKyleBot(team int, board *santorini.Board) KyleBot {
	enemy := 2
	if team == 2 {
		enemy = 1
	}
	return KyleBot{
		Team:      team,
		EnemyTeam: enemy,
		Board:     board,
	}
}

func (bot KyleBot) SelectTurn() *santorini.Turn {
	// worker1 := bot.Board.GetWorkerTile(bot.Team, 1)
	// worker2 := bot.Board.GetWorkerTile(bot.Team, 2)
	// enemyWorker1 := bot.Board.GetWorkerTile(bot.EnemyTeam, 1)
	// enemyWorker2 := bot.Board.GetWorkerTile(bot.EnemyTeam, 2)

	candidates := bot.Board.GetValidTurns(bot.Team)
	if candidates == nil {
		return nil
	}

	var (
		maxWeight int = -1000
		bestIndex     = 0
	)

	// Always block a win if possible
	enemyCandidates := bot.Board.GetValidTurns(bot.EnemyTeam)

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

		// Initialize Weight
		weight := 0

		// Weight based on move height
		// weight += 10 * candidate.MoveTo.GetHeight()

		// Weight based on build height
		// if !bot.hasNearbyEnemyWorker(bot.Board.GetSurroundingTiles(candidate.Build.GetX(), candidate.Build.GetY())...) {
		// 	weight += 10 * (candidate.Build.GetHeight() + 1)
		// } else {
		// 	weight -= 50 * (candidate.Build.GetHeight() + 1)
		// }

		// Avoid building near enemy workers
		// var enemyReach []santorini.Tile
		// enemyMovement := append(
		// 	bot.Board.GetSurroundingTiles(enemyWorker1.GetX(), enemyWorker1.GetY()),
		// 	bot.Board.GetSurroundingTiles(enemyWorker2.GetX(), enemyWorker2.GetY())...,
		// )
		// for _, tile := range enemyMovement {
		// 	if candidate.Build.GetX() == tile.GetX() || candidate.Build.GetY() == tile.GetY() {
		// 		weight -= 10
		// 	}
		// 	enemyReach = append(enemyReach, bot.Board.GetSurroundingTiles(tile.GetX(), tile.GetY())...)
		// }
		// for _, tile := range enemyReach {
		// 	if candidate.Build.GetX() == tile.GetX() || candidate.Build.GetY() == tile.GetY() {
		// 		weight -= 5
		// 	}
		// }

		// Prefer building on edge tiles
		// if candidate.Build.GetX() == 0 ||
		// 	candidate.Build.GetX() == bot.Board.Size-1 ||
		// 	candidate.Build.GetY() == 0 ||
		// 	candidate.Build.GetY() == bot.Board.Size-1 {
		// 	weight += 10
		// }

		// Avoid moves that enable an enemy win next turn
		thoughtBoard := bot.copyBoard()
		thoughtBoard.PlayTurn(candidate)
		futureEnemyCandidates := thoughtBoard.GetValidTurns(bot.EnemyTeam)
		for _, enemyCandidate := range futureEnemyCandidates {
			if enemyCandidate.IsVictory() {
				weight -= 10000
				break
			}
		}

		if weight > maxWeight {
			bestIndex = index
		}
	}

	return &candidates[bestIndex]
}

func (bot KyleBot) Name() string {
	return "Kyle Bot"
}

func (bot KyleBot) copyBoard() santorini.Board {
	tiles := make([]santorini.Tile, len(bot.Board.Tiles))
	copy(tiles, bot.Board.Tiles)
	return santorini.Board{
		Size:  bot.Board.Size,
		Tiles: tiles,
	}
}

func (bot KyleBot) hasNearbyEnemyWorker(tiles ...santorini.Tile) bool {
	for _, tile := range tiles {
		if tile.GetTeam() == bot.EnemyTeam {
			return true
		}
	}

	return false
}
