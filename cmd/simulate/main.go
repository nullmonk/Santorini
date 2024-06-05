package main

import (
	"fmt"
	"os"
	"santorini/santorini"
	"strconv"
	"sync"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
)

type options struct {
	threadCount int `help:"Number of threads to use"`
	simCount    int
}

type overallstats struct {
	bot1Wins int
	bot2Wins int
	// Calculate average round count
	sumRounds int
	pb        *progressbar.ProgressBar
	dumpLoss  bool
}

func (stats *overallstats) update(sim *santorini.Game) {
	if sim.Victor == nil {
		os.Stdout.WriteString(sim.GetTextLog())
		return
	}
	// The first bot can be team 1 or team 2 depending on the round number
	if int(sim.Victor.Team) == sim.Id%2+1 {
		stats.bot1Wins++
	} else {
		stats.bot2Wins++
		// Keep track of the losses
		if stats.dumpLoss {
			os.Stdout.WriteString(sim.GetTextLog())
		}
	}
	stats.sumRounds += sim.TurnCount / 2
	if stats.pb != nil {
		stats.pb.Describe(fmt.Sprintf("%03d / %03d", stats.bot1Wins, stats.bot2Wins))
		stats.pb.Add(1)
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Chose two bots to simulate. Bots will alternate going first. Deterministic bots will only run 1 game each.")
		fmt.Printf("USAGE: %s bot1 bot2 [numRounds]\n", os.Args[0])
		os.Exit(1)
	}

	bot1, err := santorini.PlayerFromFile(os.Args[1])
	if err != nil {
		fmt.Printf("%s is not a known bot. %s\n", os.Args[1], err)
		os.Exit(1)
	}
	bot2, err := santorini.PlayerFromFile(os.Args[2])
	if err != nil {
		fmt.Printf("%s is not a known bot. %s\n", os.Args[1], err)
		os.Exit(1)
	}
	opts := &options{
		threadCount: 10,
		simCount:    1000,
	}

	/* TODO
	if b1.IsDeterministic() && b2.IsDeterministic() {
		opts.simCount = 2
	}
	*/
	if len(os.Args) > 3 {
		if i, err := strconv.ParseInt(os.Args[3], 10, 64); err == nil {
			opts.simCount = int(i)
		} else {
			fmt.Println("Cannot parse integer", os.Args[3])
			os.Exit(1)
		}
	}

	logrus.Infof("Running %d simulations between %s and %s", opts.simCount, bot1.Name, bot2.Name)
	stats := &overallstats{
		pb:       progressbar.Default(int64(opts.simCount), "0 / 0"),
		dumpLoss: true,
	}

	wg := new(sync.WaitGroup)
	wg2 := new(sync.WaitGroup)
	sims := make(chan *santorini.Game)
	completedSims := make(chan *santorini.Game)

	logrus.Debugf("Starting %d workers", opts.threadCount)
	for i := 0; i < opts.threadCount; i++ {
		wg.Add(1)
		go runner(wg, sims, completedSims)
	}
	wg2.Add(1)
	go statistician(wg2, completedSims, stats)

	// run all the sim
	for i := 0; i < opts.simCount; i++ {
		var sim *santorini.Game
		if i%2 == 0 {
			sim = santorini.NewGame(i, nil, bot1, bot2)
		} else {
			sim = santorini.NewGame(i, nil, bot2, bot1)
		}
		sims <- sim
	}

	// Wait for all the sims to finish
	logrus.Debug("Waiting for runners to finish")
	close(sims)
	wg.Wait()
	// Wait for the stats to finish
	logrus.Debug("Waiting for stats to finish")
	close(completedSims)
	wg2.Wait()

	logrus.WithFields(map[string]interface{}{
		"bot1":             bot1.Name,
		"bot1_wins":        stats.bot1Wins,
		"bot2":             bot2.Name,
		"bot2_wins":        stats.bot2Wins,
		"avg_round_length": stats.sumRounds / opts.simCount,
		"num_rounds":       opts.simCount,
	}).Info("Simulation Complete")
}

func runner(wg *sync.WaitGroup, sims chan *santorini.Game, results chan *santorini.Game) {
	defer wg.Done()
	defer logrus.Debug("Runner finished")
	for sim := range sims {
		sim.Finish()
		results <- sim
	}
}

func statistician(wg *sync.WaitGroup, results chan *santorini.Game, stats *overallstats) {
	defer wg.Done()
	for sim := range results {
		stats.update(sim)
	}
}
