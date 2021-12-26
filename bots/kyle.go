package bots

import (
	santorini "santorini/pkg"
)

const maxDepth = 1

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

		weight := bot.getWeight(candidate)

		if weight > maxWeight {
			maxWeight = weight
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

func (bot KyleBot) getWeight(candidate santorini.Turn) int {
	// Initialize Weight
	weight := 0

	// Prefer to move up
	weight += candidate.MoveTo.GetHeight() * 20

	// Prefer to cover the most tiles
	weight += len(bot.Board.GetSurroundingTiles(candidate.MoveTo.GetX(), candidate.MoveTo.GetY()))

	// Prefer to build high if no enemies are near
	if !bot.hasNearbyEnemyWorker(candidate.Team, candidate.Build) {
		weight += (candidate.Build.GetHeight() + 1) * 10
	}

	// Don't build what you cannot reach
	if candidate.MoveTo.GetHeight() < candidate.Build.GetHeight() {
		weight -= 50
	}

	// Ponder the moves to come
	thoughtBoard := bot.copyBoard()
	thoughtBoard.PlayTurn(candidate)

	// Prefer moves that enable us to win next turn
	futureCandidates := thoughtBoard.GetValidTurns(bot.Team)
	for _, futureCandidate := range futureCandidates {
		if futureCandidate.IsVictory() {
			weight += 1000
		}
	}

	// Avoid moves that enable an enemy win next turn
	futureEnemyCandidates := thoughtBoard.GetValidTurns(bot.EnemyTeam)
	for _, futureEnemyCandidate := range futureEnemyCandidates {
		if futureEnemyCandidate.IsVictory() {
			weight -= 100000
		}
	}

	return weight
}

func (bot KyleBot) hasNearbyEnemyWorker(friendly int, tile santorini.Tile) bool {
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
