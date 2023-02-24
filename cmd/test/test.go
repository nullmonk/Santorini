package main

import (
	"os"
	"santorini/bots"
	"santorini/santorini"
)

func playGame(lgr *santorini.GameLog, board santorini.Board, bts []santorini.TurnSelector) int {
	turnCount := 0
	for {
		turnCount++
		for i, b := range bts {
			turn := b.SelectTurn(board)
			lgr.LogTurn(turn, b.Name())
			if v, _ := board.PlayTurn(turn); v {
				lgr.Comment("Engine", "%s wins in %d moves", b.Name(), turnCount)
				return i
			}
		}
	}
	return 0
}

func main() {
	board := santorini.NewBoard()
	lgr := santorini.NewGameLog(board, "AuBot", "KyleBot")
	bts := []santorini.TurnSelector{
		bots.NewAuBot("model2.json", false)(1, board, lgr),
		bots.NewKyleBot(2, board, lgr),
		//bots.NewAuBot("model2.json", false)(2, board, lgr),
	}
	playGame(lgr, board, bts)
	lgr.Dump(os.Stdout)
	board = santorini.NewBoard()
	lgr = santorini.NewGameLog(board, "AuBot", "KyleBot")
	bts = []santorini.TurnSelector{
		bots.NewKyleBot(1, board, lgr),
		bots.NewAuBot("model2.json", false)(2, board, lgr),
	}
	playGame(lgr, board, bts)
	lgr.Dump(os.Stdout)
}
