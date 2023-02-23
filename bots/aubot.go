package bots

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"santorini/santorini"

	"github.com/sirupsen/logrus"
)

// Bot implements a basic bot that allows for easy implementation of new bots
type AuBot struct {
	team       uint8
	EnemyTeams []uint8
	Logger     *logrus.Logger

	model    *santorini.AuNet
	modelFil string

	trainingModel *santorini.AuNet
	learn         bool
}

func NewAuBot(filename string, learn bool) santorini.BotInitializer {
	model, err := santorini.LoadAuNet(filename)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Loaded the model %v\n", model)
	return func(team uint8, board santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
		return &AuBot{
			team:          team,
			EnemyTeams:    GetEnemyTeams(board, team),
			model:         model,
			modelFil:      filename,
			Logger:        logger,
			trainingModel: santorini.NewAuNet(),
			learn:         learn,
		}
	}
}

func (r *AuBot) Name() string {
	return "AuBot"
}

func (r *AuBot) Team() uint8 {
	return r.team
}

func (r *AuBot) IsDeterministic() bool {
	return false
}

func (r *AuBot) Log(fmtstr string, args ...interface{}) {
	r.Logger.Debugf("AuBot: ", fmt.Sprintf(fmtstr, args...))
}

func (r *AuBot) GameOver(win bool) {
	if !r.learn {
		return
	}

	// update our model with the training model
	if win {
		r.model.Add(r.trainingModel)
	} else {
		r.model.Sub(r.trainingModel)
	}
	if len(r.modelFil) > 0 {
		r.Log("saving updated model to %s, win=%v", r.modelFil, win)
		r.model.Save(r.modelFil)
	}

	// reset the traing model
	r.trainingModel = santorini.NewAuNet()
}

// StopWin returns a list of moves that can block the enemy, the returned moves are then passed to the rank function.
// Note: This overrides ALL other moves in the candidate pool, we must try to stop the enemy
func (r *AuBot) StopWin(b santorini.Board, candidates []*santorini.Turn) []*santorini.Turn {
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
func (r *AuBot) SelectTurn(b santorini.Board) *santorini.Turn {
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

	bestWeight := -11111111111111
	bestIndexes := make([]int, 0, 1)
	// Multiple Best Ranks?
	for i, t := range candidates {
		rank := r.model.RankTurn(b, t)
		if rank > bestWeight {
			bestWeight = rank
			bestIndexes = []int{i}
		} else if rank == bestWeight {
			bestIndexes = append(bestIndexes, i)
		}
	}

	turn := candidates[bestIndexes[0]]

	if len(bestIndexes) > 1 {
		// pick ONE of the best options
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(bestIndexes))))
		if err != nil {
			panic(err)
		}
		turn = candidates[bestIndexes[n.Int64()]]
	}

	// update our training model with the turn
	r.trainingModel.AddTurn(b, turn)
	return turn
}

func (r AuBot) TrainingModel() *santorini.AuNet {
	return r.trainingModel
}
