package santorini

type Tile struct {
	Height int     // 0 no building, 4 capped
	Worker *Worker `json:",omitempty"`

	x uint8
	y uint8
}

func (t Tile) IsOccupied() bool {
	return t.Worker != nil
}
func (t Tile) IsOccupiedBy(worker Worker) bool {
	if t.Worker == nil {
		return false
	}

	return t.Worker.Team == worker.Team && t.Worker.Number == worker.Number
}

func (t Tile) IsCapped() bool {
	return t.Height > 3
}

func (t Tile) GetX() uint8 {
	return t.x
}

func (t Tile) GetY() uint8 {
	return t.y
}
