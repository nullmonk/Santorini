package santorini

import "encoding/json"

type Tile struct {
	team   int // 0 no team, otherwise team number
	worker int // 0 no worker, otherwise a worker is present
	height int // 0 no building, 4 capped
	x      int // x position of the tile
	y      int // y position of the tile
}

func (t Tile) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Team   int `json:"team,omitempty"`
		Worker int `json:"worker,omitempty"`
		Height int `json:"height,omitempty"`
		X      int `json:"x"`
		Y      int `json:"y"`
	}{
		t.team, t.worker, t.height, t.x, t.y,
	})
}

func (t Tile) IsOccupied() bool {
	return t.team != 0 || t.worker != 0
}

func (t Tile) IsOccupiedBy(team, worker int) bool {
	return t.team == team && t.worker == worker
}

func (t Tile) IsCapped() bool {
	return t.height > 3
}

func (t Tile) GetX() int {
	return t.x
}

func (t Tile) GetY() int {
	return t.y
}

func (t Tile) GetTeam() int {
	return t.team
}

func (t Tile) GetWorker() int {
	return t.worker
}

func (t Tile) GetHeight() int {
	return t.height
}
