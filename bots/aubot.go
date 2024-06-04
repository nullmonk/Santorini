package bots

import (
	"crypto/rand"
	"math/big"
	"santorini/santorini"
)

// Bot implements a basic bot that allows for easy implementation of new bots
type AuBot struct {
	team       uint8
	EnemyTeams []uint8
	Logger     *santorini.GameLog

	model    *santorini.AuNet
	modelFil string

	turnCount     int
	trainingModel *santorini.AuNet
	learn         bool
	save          bool
}

func NewAuBot(filename string, learn ...bool) santorini.BotInitializer {
	model, err := santorini.LoadAuNet(filename)
	if err != nil {
		panic(err)
	}
	lrn := false // Learn as we go
	if len(learn) > 0 {
		lrn = learn[0]
	}
	save := false // save progress as we go
	if len(learn) > 1 {
		save = learn[1]
	}
	return func(team uint8, board santorini.Board, logger *santorini.GameLog) santorini.TurnSelector {
		return &AuBot{
			team:          team,
			EnemyTeams:    GetEnemyTeams(board, team),
			model:         model,
			modelFil:      filename,
			Logger:        logger,
			trainingModel: santorini.NewAuNet(),
			learn:         lrn,
			save:          save,
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
	if r.Logger == nil {
		return
	}
	r.Logger.Comment("AuBot", fmtstr, args...)
}

func (r *AuBot) GameOver(win bool) {
	// update our model with the training model
	if !r.learn {
		r.Log("this model does not learn")
		return // this model does not learn
	}
	if win {
		r.model.Add(r.trainingModel)
	} else {
		r.model.Sub(r.trainingModel)
	}
	if len(r.modelFil) > 0 && r.save {
		r.Log("saving updated model to %s, win=%v", r.modelFil, win)
		r.model.Save(r.modelFil)
	}
	r.Log("this model does not learn")

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
	bestIndexes := make(map[string]*santorini.Turn, 1)
	lastI := 0
	// Multiple Best Ranks?
	r.turnCount++
	for i, t := range candidates {
		rank := r.model.RankTurn(b, t)
		if rank > bestWeight {
			bestWeight = rank
			bestIndexes = make(map[string]*santorini.Turn, 1)
			// reset indexes and add ours
			bestIndexes[t.String()] = t
			lastI = i
		} else if rank == bestWeight {
			bestIndexes[t.String()] = t
			lastI = i
		}
	}

	turn := candidates[lastI]

	if len(bestIndexes) > 1 {
		// pick ONE of the best options

		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(bestIndexes))))
		if err != nil {
			panic(err)
		}
		i := 0
		for _, t := range bestIndexes {
			if i == int(n.Int64()) {
				turn = t
				break
			}
			i++
		}
		// dump these turn states bc we want to find the differences and calculate them
		if r.turnCount > 5 {
			turns := make([]string, 0, len(bestIndexes))
			for k := range bestIndexes {
				if turn.String() == k {
					k = "*" + k
				}
				turns = append(turns, k)
			}
			r.Log("Could not decide between the following turns [weight=%d]: %s", bestWeight, turns)
		}
	}

	// update our training model with the turn
	r.trainingModel.AddTurn(b, turn)
	return turn
}

func (r AuBot) TrainingModel() *santorini.AuNet {
	return r.trainingModel
}
