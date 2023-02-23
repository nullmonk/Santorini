package santorini

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
)

// Implement some sort of whatever for a custom "neural net", really just a bunch of weights that choose the turn

const numWeights = 20

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
	count   int
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
	tr := t.toVector(b)
	a.mu.Lock()
	defer a.mu.Unlock()
	for i := 0; i < numWeights; i++ {
		a.weights[i] += tr.Pop() * (a.count / 2)
	}
	a.count++
}

func (a *AuNet) RankTurn(b Board, t *Turn) int {
	res := 0
	tr := t.toVector(b)
	for i := 0; i < numWeights; i++ {
		res += a.weights[i] * tr.Pop()
	}
	return res
}

const shrinkFactor = 3

func (a *AuNet) Add(b *AuNet) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range b.weights {
		a.weights[i] += v / shrinkFactor
	}
}

func (a *AuNet) Sub(b *AuNet) {
	a.mu.Lock()
	defer a.mu.Unlock()
	for i, v := range b.weights {
		a.weights[i] -= v / shrinkFactor
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
	d, err := os.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	a = &AuNet{}
	if err = json.Unmarshal(d, &a.weights); err != nil {
		return nil, err
	}
	if len(a.weights) < numWeights {
		a.weights = make([]int, numWeights)
	}
	return a, nil
}

func (t *Turn) toVector(b Board) (tr *trait) {
	tr = NewTrait()
	tr.Push(t.Worker.height < t.MoveTo.height)                         // Jumping Up
	tr.Push(t.Worker.height > t.MoveTo.height)                         // Jumping Down
	tr.Push(t.Worker.height > t.MoveTo.height && t.MoveTo.height == 0) // Jumping Down
	tr.Push(t.MoveTo.height == 2)                                      // Moving up
	tr.Push(t.Build.height == 2)                                       // Building a 3rd row
	tr.Push(t.Build.height == 1)                                       // Building off the ground

	// Both workers off ground?
	//

	var eAtWorker, eAtMove, eAtBuild bool
	var fAtMove, fAtBuild bool
	var tower, standOnTower = true, true // if we are building a high point or starting on a high point
	var buildOnEdge, buildInCorner bool
	var moveOnEdge, moveInCorner bool

	// Track if the enemy can hop up onto our builds
	var eCanWin, eCanJump bool

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
	surrounding = b.GetSurroundingTiles(t.Build.x, t.Build.y)
	buildInCorner = len(surrounding) == 3 // corner
	buildOnEdge = len(surrounding) == 5   // edge
	for _, st := range surrounding {
		if st.team > 0 {
			if st.team != t.Worker.team {
				eAtBuild = true
				if t.Build.height == st.height {
					if st.height == 2 {
						eCanWin = true
					}
					eCanJump = true
				}
			} else {
				fAtBuild = true
			}
		}
		if st.height > t.Build.height {
			tower = false
		}
	}
	tr.Push(eAtWorker, eAtMove, eAtBuild)
	tr.Push(fAtMove, fAtBuild, tower, standOnTower)
	tr.Push(buildInCorner, buildOnEdge, moveInCorner, moveOnEdge)
	tr.Push(eCanJump, eCanWin)
	return
}
