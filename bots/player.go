package bots

import (
	"fmt"
	santorini "santorini/pkg"
	"santorini/pkg/color"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// PlayerBot allows a player to make moves on the game board
type PlayerBot struct {
	name   string
	Team   int
	logger *logrus.Logger
	Board  *santorini.Board
}

func NewPlayerBot(team int, board *santorini.Board, logger *logrus.Logger) santorini.TurnSelector {
	p := &PlayerBot{
		Team:   team,
		Board:  board,
		logger: logger,
		name:   "PlayerBot",
	}
	fmt.Printf("[PlayerBot] Enter name for player (team %d): ", team)
	var ret string
	fmt.Scanln(&ret)
	if ret != "" {
		p.name = ret
	}
	return p
}

func (p PlayerBot) IsDeterministic() bool {
	return false
}

func (p PlayerBot) GetInput(prompt string, args ...interface{}) string {
	prompt = strings.TrimRight(fmt.Sprintf(prompt, args...), ": \t\n")
	fmt.Printf("[%s] %s: ", p.name, prompt)
	var ret string
	if _, err := fmt.Scanln(&ret); err != nil {
		panic(err)
	}
	return ret
}

func (p PlayerBot) GetChoice(prompt string, options map[string]interface{}) interface{} {
	if len(options) == 1 {
		return 0
	}
	prompt = strings.TrimRight(prompt, ": \t\n")

	answers := make([]interface{}, 0, len(options))
	i := 0
	for name, choice := range options {
		fmt.Printf("\t%d. %s\n", i, name)
		i++
		answers = append(answers, choice)
	}
	for {
		fmt.Printf("[%s] %s: ", p.name, prompt)
		var answer string
		fmt.Scanln(&answer)
		if i, err := strconv.ParseInt(answer, 0, 64); err == nil {
			if int(i) < len(answers) {
				return answers[i]
			}
		}
	}
}

func (p PlayerBot) Name() string {
	return p.name
}

func (p *PlayerBot) SelectTurn() *santorini.Turn {
	// Get the workers
	workers := p.Board.GetWorkerTiles(p.Team)

	movableWorkerTiles := make(map[santorini.Tile][]santorini.Tile)

	options := make(map[string]interface{})
	for _, worker := range workers {
		movableWorkerTiles[worker] = p.Board.GetMoveableTiles(worker)
		if len(movableWorkerTiles[worker]) > 0 {
			name := fmt.Sprintf("%sWorker %d (%d, %d)%s",
				color.GetWorkerColor(p.Team, worker.GetWorker()),
				worker.GetWorker(),
				worker.GetX(),
				worker.GetY(),
				color.Reset)
			options[name] = worker
		}
	}

	worker := p.GetChoice("Choose a worker", options).(santorini.Tile)

	turn := &santorini.Turn{
		Team:   p.Team,
		Worker: worker.GetWorker(),
	}

	options = make(map[string]interface{})
	for _, tile := range movableWorkerTiles[worker] {
		heightDiff := tile.GetHeight() - worker.GetHeight()
		label := ""
		if heightDiff > 0 {
			label = "up 1 tile"
		} else if heightDiff == -1 {
			label = "down 1 tile"
		} else if heightDiff == -1 {
			label = "down 2 tiles"
		}

		if tile.GetHeight() == 3 {
			label = "Winning Move!"
		}
		name := fmt.Sprintf("%s (%d,%d)", GetTileDir(worker, tile), tile.GetX(), tile.GetY())
		if label != "" {
			name += " - " + label
		}
		options[name] = tile
	}

	turn.MoveTo = p.GetChoice("Choose a tile to move to", options).(santorini.Tile)

	fmt.Println()
	p.simulateMove(worker, turn.MoveTo)
	fmt.Println()

	options = make(map[string]interface{})
	for _, tile := range p.Board.GetBuildableTiles(p.Team, worker.GetWorker(), turn.MoveTo) {
		name := fmt.Sprintf("%s (%d,%d)", GetTileDir(turn.MoveTo, tile), tile.GetX(), tile.GetY())
		label := ""
		if tile.GetHeight() == 3 {
			label = "Cap tile"
		}
		if label != "" {
			name += " - " + label
		}
		options[name] = tile
	}

	turn.Build = p.GetChoice("Choose a tile to build on", options).(santorini.Tile)
	return turn
}

// Print out the board after a move has been taken
func (p *PlayerBot) simulateMove(src, dst santorini.Tile) {
	copyBoard := &santorini.Board{
		Tiles: p.Board.Tiles,
		Size:  p.Board.Size,
		Teams: make(map[int]bool),
	}

	copyBoard.PlaceWorker(0, 0, src.GetX(), src.GetY())
	copyBoard.PlaceWorker(src.GetTeam(), src.GetWorker(), dst.GetX(), dst.GetY())
	fmt.Println(copyBoard)
}

// ⇐⇑⇒⇓⇖⇗⇘⇙
// ←↑→↓↖↗↘↙

func GetTileDir(src, dst santorini.Tile) string {
	dx := dst.GetX() - src.GetX()
	dy := dst.GetY() - src.GetY()
	if dx < 0 && dy == 0 {
		return "←" // left
	}
	if dx == 0 && dy < 0 {
		return "↑" // up
	}
	if dx > 0 && dy == 0 {
		return "→" // right
	}
	if dx == 0 && dy > 0 {
		return "↓" // down
	}
	if dx < 0 && dy < 0 {
		return "↖" // up, left
	}
	if dx > 0 && dy < 0 {
		return "↗" // up, right
	}
	if dx > 0 && dy > 0 {
		return "↘" // down, right
	}
	if dx < 0 && dy > 0 {
		return "↙" // down, left
	}
	return "unkowno"
}
