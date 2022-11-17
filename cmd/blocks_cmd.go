package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/cortze/eth2-state-analyzer/pkg/blocks"
	"github.com/cortze/eth2-state-analyzer/pkg/clientapi"
	"github.com/cortze/eth2-state-analyzer/pkg/utils"
)

var BlocksCommand = &cli.Command{
	Name:   "blocks",
	Usage:  "analyze the Beacon Block of a given slot range",
	Action: LaunchBlockMetrics,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "bn-endpoint",
			Usage: "beacon node endpoint (to request the BeaconStates)",
		},
		&cli.IntFlag{
			Name:  "init-slot",
			Usage: "init slot from where to start",
		},
		&cli.IntFlag{
			Name:  "final-slot",
			Usage: "init slot from where to finish",
		},
		&cli.StringFlag{
			Name:  "log-level",
			Usage: "log level: debug, warn, info, error",
		},
		&cli.StringFlag{
			Name:  "db-url",
			Usage: "example: postgresql://beaconchain:beaconchain@localhost:5432/beacon_states_kiln",
		},
		&cli.IntFlag{
			Name:  "workers-num",
			Usage: "example: 50",
		},
		&cli.IntFlag{
			Name:  "db-workers-num",
			Usage: "example: 50",
		}},
}

var logBlocks = logrus.WithField(
	"module", "BlocksCommand",
)

// CrawlAction is the function that is called when running `eth2`.
func LaunchBlockMetrics(c *cli.Context) error {
	coworkers := 1
	dbWorkers := 1
	logRewardsRewards.Info("parsing flags")
	// check if a config file is set
	if !c.IsSet("bn-endpoint") {
		return errors.New("bn endpoint not provided")
	}
	if !c.IsSet("init-slot") {
		return errors.New("final slot not provided")
	}
	if !c.IsSet("final-slot") {
		return errors.New("final slot not provided")
	}
	if c.IsSet("log-level") {
		logrus.SetLevel(utils.ParseLogLevel(c.String("log-level")))
	}
	if !c.IsSet("db-url") {
		return errors.New("db-url not provided")
	}
	if !c.IsSet("workers-num") {
		logRewardsRewards.Infof("workers-num flag not provided, default: 1")
	} else {
		coworkers = c.Int("workers-num")
	}
	if !c.IsSet("db-workers-num") {
		logRewardsRewards.Infof("db-workers-num flag not provided, default: 1")
	} else {
		dbWorkers = c.Int("db-workers-num")
	}
	bnEndpoint := c.String("bn-endpoint")
	initSlot := uint64(c.Int("init-slot"))
	finalSlot := uint64(c.Int("final-slot"))
	dbUrl := c.String("db-url")

	// generate the httpAPI client
	cli, err := clientapi.NewAPIClient(c.Context, bnEndpoint, QueryTimeout)
	if err != nil {
		return err
	}

	// generate the block analyzer
	blockAnalyzer, err := blocks.NewBlockAnalyzer(c.Context, cli, initSlot, finalSlot, dbUrl, coworkers, dbWorkers)
	if err != nil {
		return err
	}

	procDoneC := make(chan struct{})
	sigtermC := make(chan os.Signal)

	signal.Notify(sigtermC, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, syscall.SIGTERM)

	go func() {
		blockAnalyzer.Run()
		procDoneC <- struct{}{}
	}()

	select {
	case <-sigtermC:
		logBlocks.Info("Sudden shutdown detected, controlled shutdown of the cli triggered")
		blockAnalyzer.Close()

	case <-procDoneC:
		logBlocks.Info("Process successfully finish!")
	}
	close(sigtermC)
	close(procDoneC)

	return nil
}
