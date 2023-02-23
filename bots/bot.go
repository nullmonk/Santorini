package bots

import (
	"fmt"
	"santorini/santorini"

	"github.com/sirupsen/logrus"
)

// RankFunction is a function that takes a turn, and ranks how good it is. The move with the highest rank gets chosen
type RankFunction func(*Bot, santorini.Board, *santorini.Turn) int

// RankAllFunction gets a list of all the moves to allow the bot to do whatever it wants with them
// NOTE: If RankAllFunction is defined for a bot, RankFunction will not be called
type RankAllFunction func(*Bot, santorini.Board, []*santorini.Turn) *santorini.Turn

// Bot implements a basic bot that allows for easy implementation of new bots
type Bot struct {
	name          string
	team          uint8
	EnemyTeams    []uint8
	Logger        *logrus.Logger
	deterministic bool

	rankF    RankFunction
	rankAllF RankAllFunction
}

// NewBasicBot builds a basic bot that sorts turn using rf
func NewBasicBot(name string, rf RankFunction, deterministic bool) santorini.BotInitializer {
	// If f is nil, panic
	if rf == nil {
		panic(fmt.Errorf("No RankFunction specified for %s", name))
	}
	return func(team uint8, board santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
		return &Bot{
			name:          name,
			rankF:         rf,
			team:          team,
			EnemyTeams:    GetEnemyTeams(board, team),
			deterministic: deterministic,
		}
	}
}

// NewAdvancedBot builds a bot that sorts all the functions using a RankAllFunction. This bot requires the user to do a bit more work
func NewAdvancedBot(name string, rf RankAllFunction, deterministic bool) santorini.BotInitializer {
	// If f is nil, panic
	if rf == nil {
		panic(fmt.Errorf("No RankFunction specified for %s", name))
	}
	return func(team uint8, board santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
		return &Bot{
			name:          name,
			rankAllF:      rf,
			team:          team,
			EnemyTeams:    GetEnemyTeams(board, team),
			deterministic: deterministic,
		}
	}
}

func (r Bot) Name() string {
	return r.name
}

func (r Bot) Team() uint8 {
	return r.team
}

func (r Bot) IsDeterministic() bool {
	return r.deterministic
}

func (r Bot) Log(fmtstr string, args ...interface{}) {
	if r.Logger == nil {
		return
	}
	r.Logger.Debugf(r.name+": ", fmt.Sprintf(fmtstr, args...))
}

func (r Bot) GameOver(win bool) {
	if win {
		r.Log("won the game!")
	} else {
		r.Log("lost the game")
	}
}

// StopWin returns a list of moves that can block the enemy, the returned moves are then passed to the rank function.
// Note: This overrides ALL other moves in the candidate pool, we must try to stop the enemy
func (r Bot) StopWin(b santorini.Board, candidates []*santorini.Turn) []*santorini.Turn {
	newCandidates := make([]*santorini.Turn, 0, 4)
	OnEnemyWinningMoves(b, r.EnemyTeams, func(t *santorini.Turn) error {
		for _, defense := range candidates {
			if t.MoveTo.SameLocation(defense.Build) {
				// Add it to our new candidate pool
				newCandidates = append(newCandidates, defense)
			}
		}
		return nil
	})
	return newCandidates
}

// SelectTurn at random, returns nil if no move can be made
func (r *Bot) SelectTurn(b santorini.Board) *santorini.Turn {
	// Get the valid turns from here
	candidates := b.ValidTurns(r.team)
	if candidates == nil || len(candidates) == 0 {
		return nil
	}

	// If we can win this turn, do it
	if wins := GetWinningMoves(candidates); len(wins) > 0 {
		return wins[0]
	}

	// If we need to defend this turn, do it
	if defences := r.StopWin(b, candidates); len(defences) > 0 {
		candidates = defences
	}

	if r.rankAllF != nil {
		return r.rankAllF(r, b, candidates)
	}

	bestWeight := -11111111111111
	bestIndex := 0
	for i, t := range candidates {
		rank := r.rankF(r, b, t)
		if rank > bestWeight {
			bestWeight = rank
			bestIndex = i
		}
	}
	return candidates[bestIndex]
}
