package bots

import "santorini/santorini"

/* Functions that can help bots out */
func GetEnemyTeams(b santorini.Board, myteam uint8) []uint8 {
	teams := b.Teams()
	enemyTeams := make([]uint8, 0, len(teams)-1)
	for _, t := range teams {
		if t != myteam {
			enemyTeams = append(enemyTeams, t)
		}
	}
	return enemyTeams
}

// GetWinningMoves returns a list of all the winning moves in the given pool
func GetWinningMoves(candidates []*santorini.Turn) []*santorini.Turn {
	result := make([]*santorini.Turn, 0, len(candidates))
	for _, t := range candidates {
		if t.MoveTo.GetHeight() == 3 {
			result = append(result, t)
		}
	}
	return result
}

// OnEnemyWinningMoves calls F for every move an enemy team can make that would cause them to win
func OnEnemyWinningMoves(b santorini.Board, enemyTeams []uint8, f func(*santorini.Turn) error) error {
	// Loop through all the other teams
	for _, team := range enemyTeams {
		for _, turn := range b.ValidTurns(team) {
			// Get all the winning moves for that team
			if turn.MoveTo.GetHeight() == 3 {
				if err := f(turn); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
