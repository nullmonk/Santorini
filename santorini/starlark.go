package santorini

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"santorini/santorini/modules"
	"strings"
	"time"

	starjson "go.starlark.net/lib/json"
	starmath "go.starlark.net/lib/math"
	"go.starlark.net/starlark"
	"go.starlark.net/syntax"
)

const ExpireTime = time.Hour * 24 // Games get removed after 24 hours

type Player struct {
	Name   string `json:"name"`   // The name of the bot
	Source any    `json:"source"` // Bot source code
	Team   int    `json:"team"`   // The team number this bot is

	logName string              // NAME (team number)
	globals starlark.StringDict // Bot specific globals

	selectTurn starlark.Value // Bots select turn function
	gameOver   starlark.Value // Bots function that is called on game over
}

func PlayerFromFile(fn string) (*Player, error) {
	b, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	_, name := filepath.Split(fn)
	return &Player{
		Name:   name,
		Source: string(b),
	}, nil
}

func (p *Player) Copy() *Player {
	return &Player{
		Name:   p.Name,
		Source: p.Source,
	}
}

// Have the bot select one of the possible turns
func (g *Game) getTurn(player *Player, options []*Turn) (*Turn, error) {
	turns_starlark := starlark.NewList(nil)
	for _, t := range options {
		turns_starlark.Append(*t)
	}
	args := starlark.Tuple{
		g.board,
		turns_starlark,
	}
	// Pass the turns to the bot
	result, err := starlark.Call(g.thread, player.selectTurn, args, nil)
	if err != nil {
		return nil, fmt.Errorf("invalid result from bot exec: %s", err)
	}

	var turn *Turn
	switch t := result.(type) {
	case Turn:
		turn = &t
	case starlark.String:
		// Got a turn hash
		turn, _ = TurnFromString(0, t.String())
	default:
		return nil, fmt.Errorf("invalid result from bot. Expected Turn or string: type=%s value=%s", t.Type(), t.String())
	}
	if turn == nil {
		return nil, fmt.Errorf("invalid result from bot. Expected Turn or string: type=%s value=%s", result.Type(), result.String())
	}
	return turn, nil
}

// setup the bot
func (p *Player) init(g *Game) error {
	opts := &syntax.FileOptions{
		Set:             true,
		While:           true,
		TopLevelControl: false,
		GlobalReassign:  false,
		Recursion:       true,
	}
	var err error
	p.globals, err = starlark.ExecFileOptions(opts, g.thread, p.Name, p.Source, g.globals)
	if err != nil {
		return err
	}

	// todo validate the bot here
	p.selectTurn = p.globals["take_turn"]
	if p.selectTurn == nil {
		fmt.Printf("%s: %+v\n", p.Name, p.globals)
	}
	p.gameOver = p.globals["game_over"] // TODO: Implement hooks for start and end
	if p.Team != 0 {
		p.logName = fmt.Sprintf("%s (%d)", p.Name, p.Team)
	} else {
		p.logName = p.Name
	}
	return nil
}

type Game struct {
	Id        int       // Unique id for the game or sim number
	Victor    *Player   // Who won
	Players   []*Player // bots
	TurnCount int
	globals   starlark.StringDict // Any global functions or data we want to add to the environment
	thread    *starlark.Thread    // The starlark thread that this game is executing in
	board     *FastBoard
	log       io.Writer // file for the game log
	buf       []string  // buffer to write game logs too
}

func NewGame(id int, log io.Writer, players ...*Player) *Game {
	b, _ := NewBoardFromHash("CAAAAAAAIAAAQAQAAAIAAAAAAA") // Standard 2 player board
	g := &Game{
		Id:      id,
		Players: make([]*Player, 0, len(players)),
		globals: starlark.StringDict{
			"random": modules.Random,
			"json":   starjson.Module,
			"math":   starmath.Module,
			// TODO Add other modules here
		},
		board:  b,
		thread: &starlark.Thread{Name: "santorini"},
		log:    log,
	}

	if log == nil {
		g.log = new(strings.Builder)
	}
	playernames := make([]string, 0, len(players)) // only used for logging
	for i, p := range players {
		p2 := p.Copy() // Always copy the player for a new game
		p2.Team = i + 1
		p2.init(g)
		playernames = append(playernames, p2.Name)
		g.Players = append(g.Players, p2)
	}

	fmt.Fprintf(g.log, "GAME: %d\n", id)
	fmt.Fprintf(g.log, "BOARD: %s\n", b.GameHash())
	fmt.Fprintf(g.log, "TEAMS: %s\n", strings.Join(playernames, " vs. "))
	return g
}

func (g *Game) Exit(err error) {
	if err != nil {
		fmt.Fprintf(g.log, "CRASH: %s\n", err)
	}
}

func (g *Game) NextTurn() (over bool, err error) {
	// whose turn is it
	team := int(g.TurnCount)%len(g.Players) + 1
	bot := g.Players[team-1]

	// get the valid turns for this player
	turns := g.board.ValidTurns(uint8(team))

	if len(turns) == 0 {
		// This bot loses
		g.Comment(bot.logName, "loses game. No more moves")
		g.Victor = g.Players[(bot.Team+1)%len(g.Players)]
		return true, nil
	}

	turn, err := g.getTurn(bot, turns)
	if err != nil {
		return false, err
	}
	v, err := g.board.PlayTurn(turn)
	if err != nil {
		return false, err
	}

	g.TurnCount++
	g.writeComments() // Write all the comment before the turn is dumped
	fmt.Fprintf(g.log, "TURN: %s # %s\n", turn.String(), bot.logName)
	if v {
		// We won?
		g.Victor = bot
	}
	return v, nil
}

func (g *Game) Finish() error {
	// Maybe it could be more?
	for i := 0; i < 1000; i++ {
		v, err := g.NextTurn()
		g.writeComments()
		if err != nil {
			fmt.Fprintf(g.log, "CRASH: %s\n", err)
			return err
		}
		if v {
			fmt.Fprintf(g.log, "# Game %d Complete. %s wins after %d moves\n", g.Id, g.Victor.logName, g.TurnCount)
			return nil
		}
	}
	fmt.Fprintln(g.log, "CRASH: exceeded max number of simulation moves (1000)")
	return fmt.Errorf("exceeded max number of simulation moves (1000)")
}

func (g *Game) Comment(who, format string, a ...interface{}) {
	if who != "" {
		who += ": "
	}
	g.buf = append(g.buf, who+fmt.Sprintf(format, a...))
}

func (g *Game) writeComments() {
	for _, l := range g.buf {
		fmt.Fprintf(g.log, "# %s\n", l)
	}
	g.buf = g.buf[:0]
}

// Get the text log, only possible if the logger was not specified
func (g *Game) GetTextLog() string {
	if log, ok := g.log.(*strings.Builder); ok {
		return log.String()
	}
	return ""
}

func (g *Game) MarshalJSON() ([]byte, error) {
	// If we have the log contents, dump them to the DB
	log := ""
	if b, ok := g.log.(*strings.Builder); ok {
		log = b.String()
	}
	return json.Marshal(struct {
		Board      string    `json:"board"`
		TurnCount  int       `json:"turn_count"`
		Log        string    `json:"log"`
		Players    []*Player `json:"players"`
		Id         int       `json:"id"`
		LastUpdate int64     `json:"last_update"`
	}{
		g.board.GameHash(),
		g.TurnCount,
		log,
		g.Players,
		g.Id,
		time.Now().Unix(),
	})
}

func (g *Game) UnmarshalJSON(data []byte) (err error) {
	dst := struct {
		Board      string    `json:"board"`
		TurnCount  int       `json:"turn_count"`
		Log        string    `json:"log"`
		Players    []*Player `json:"players"`
		Id         int       `json:"id"`
		LastUpdate int64     `json:"last_update"`
	}{}
	if err = json.Unmarshal(data, &dst); err != nil {
		return err
	}

	g.board, err = NewBoardFromHash(dst.Board)
	if err != nil {
		return err
	}
	g.TurnCount = dst.TurnCount

	// Init the players here
	g.Players = dst.Players
	for _, p := range g.Players {
		p.init(g)
	}

	g.Id = dst.Id
	// Get the logs in there
	log := new(strings.Builder)
	if len(dst.Log) > 0 {
		log.WriteString(dst.Log)
	}
	g.log = log

	if time.Now().Unix()-dst.LastUpdate > int64(ExpireTime) {
		return fmt.Errorf("game has expired")
	}
	return nil
}
