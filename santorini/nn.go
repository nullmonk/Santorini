package santorini

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Implement some sort of whatever for a custom "neural net", really just a bunch of weights that choose the turn

const numWeights = 28

type trait struct {
	t uint32
}

func NewTrait() *trait {
	return &trait{
		t: 0,
	}
}

func (t *trait) Push(values ...bool) (ok bool) {
	for _, v := range values {
		if t.t == 0xffffffff {
			return false
		}
		t.t <<= 1
		if v {
			t.t |= 1
		}
	}
	return true
}

func (t *trait) Pop() int {
	v := int(t.t & 1)
	t.t >>= 1
	return v
}

func (t *trait) String() string {
	return fmt.Sprintf("%032s", strconv.FormatUint(uint64(t.t), 2))
}

type AuNet struct {
	mu      sync.Mutex
	weights []int
}

func NewAuNet() *AuNet {
	return &AuNet{
		weights: make([]int, numWeights),
	}
}

func (a *AuNet) AddTurn(b Board, t *Turn) {
	if t == nil {
		return
	}
	tr := t.toVectorBin(b)
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < numWeights; i++ {
		a.weights[i] += tr.Pop()
	}
}

func (a *AuNet) RankTurn(b Board, t *Turn) int {
	res := 0
	tr := t.toVectorBin(b)
	for i := 0; i < numWeights; i++ {
		res += a.weights[i] * tr.Pop()
	}
	return res
}

func (a *AuNet) Add(b *AuNet) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range b.weights {
		a.weights[i] += v / 10
	}
}

func (a *AuNet) Sub(b *AuNet) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range b.weights {
		a.weights[i] -= v / 10
	}
}

func (a *AuNet) Save(fn string) error {
	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	a.mu.Lock()
	defer a.mu.Unlock()
	return enc.Encode(a.weights)
}

func LoadAuNet(fn string) (a *AuNet, err error) {
	a = &AuNet{}
	d, err := os.ReadFile(fn)
	if err == nil {
		if err = json.Unmarshal(d, &a.weights); err != nil {
			return nil, err
		}
	}
	if len(a.weights) < numWeights {
		a.weights = make([]int, numWeights)
	}
	return a, nil
}

func (t *Turn) toVectorBin(b Board) (tr *trait) {
	// Building off the ground

	var eAtWorker, eAtMove, eAtBuild bool
	var fAtMove, fAtBuild bool
	var tower, standOnTower = true, true // if we are building a high point or starting on a high point
	var buildOnEdge, buildInCorner bool
	var moveOnEdge, moveInCorner bool

	var eCanWin, eCanClimb bool
	for _, st := range b.GetSurroundingTiles(t.Worker.x, t.Worker.y) {
		if st.team > 0 && st.team != t.Worker.team {
			eAtWorker = true // enemies
		}
	}
	surrounding := b.GetSurroundingTiles(t.MoveTo.x, t.MoveTo.y)
	moveInCorner = len(surrounding) == 3 // corner
	moveOnEdge = len(surrounding) == 5   // edge
	for _, st := range surrounding {
		if st.team > 0 {
			if st.team != t.Worker.team {
				eAtMove = true
			} else {
				fAtMove = true
			}
		}
		if st.height >= t.Worker.height {
			standOnTower = false
		}
	}
	surroundingAboveGround := 0

	surrounding = b.GetSurroundingTiles(t.Build.x, t.Build.y)
	buildInCorner = len(surrounding) == 3 // corner
	buildOnEdge = len(surrounding) == 5   // edge
	for _, st := range surrounding {
		if st.team > 0 {
			if st.team != t.Worker.team {
				eAtBuild = true
				if st.height == t.Build.height {
					eCanClimb = true
					if st.height == 2 {
						eCanWin = true
					}
				}
			} else {
				fAtBuild = true
			}
		}
		if st.height > t.Build.height {
			tower = false
		}
		if st.height > 0 {
			surroundingAboveGround++
		}
	}
	tr = NewTrait() // 19 traits
	// 10
	tr.Push(t.Worker.height < t.MoveTo.height)                         // Jumping Up
	tr.Push(t.Worker.height > t.MoveTo.height)                         // Jumping Down
	tr.Push(t.Worker.height > t.MoveTo.height && t.MoveTo.height == 0) // Jumping Down
	tr.Push(t.Worker.height == 1)
	tr.Push(t.Worker.height == 2) // Prioritize higher workers?
	tr.Push(t.MoveTo.height == 2) // Moving up
	tr.Push(t.Build.height == 2)  // Building a 3rd row
	tr.Push(t.Build.height == 1)
	tr.Push(surroundingAboveGround == 1)
	tr.Push(t.Build.height > t.MoveTo.height-1)

	// Random traits for the ai to figure out
	tr.Push(t.Build.x == t.Worker.x && t.Build.y == t.Worker.y) // like building in our old spot
	//tr.Push(t.Worker.x != t.MoveTo.x && t.Worker.y != t.MoveTo.y) // like moving diagonally
	// 13
	tr.Push(eAtWorker, eAtMove, eAtBuild)
	tr.Push(fAtMove, fAtBuild, tower, standOnTower)
	tr.Push(buildInCorner, buildOnEdge, moveInCorner, moveOnEdge)
	tr.Push(eCanClimb, eCanWin)
	return
}
