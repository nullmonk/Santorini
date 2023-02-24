package santorini

type BotInitializer func(team uint8, board Board, logger *GameLog) TurnSelector

type TurnSelector interface {
	// The name of the bot
	Name() string

	// Perform the next turn for the bot
	SelectTurn(b Board) *Turn

	// True if the bot will perform the same given the same inputs
	IsDeterministic() bool

	// Call this function whenever the game is over
	GameOver(win bool) // Tell the bot it won or lost
}
