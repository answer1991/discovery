package main

import (
	"context"
	"github.com/answer1991/discovery/engine"
	"github.com/urfave/cli"
	"os"
	"os/signal"
	"syscall"
)

const (
	flagConsul                  = "consul"
	flagDockerServiceNameLabels = "docker-service-name-labels"
	flagDockerServiceTagLabels  = "docker-service-tag-labels"
	flagTtl                     = "ttl"

	defaultServiceNameLabelKey = "io.answer1991.service.name"
	defaultServiceTagLabelKey  = "io.answer1991.service.tags"
)

var (
	defaultServiceNameLabel = new(cli.StringSlice)
	defaultServiceTagLabel  = new(cli.StringSlice)
)

func init() {
	defaultServiceNameLabel.Set(defaultServiceNameLabelKey)
	defaultServiceTagLabel.Set(defaultServiceTagLabelKey)
}

var serveCommand = cli.Command{
	Name:  "serve",
	Usage: "serve registrator as a daemon",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  flagConsul,
			Value: "127.0.0.1:8500",
			Usage: "consul address",
		},
		cli.StringSliceFlag{
			Name:  flagDockerServiceNameLabels,
			Value: defaultServiceNameLabel,
			Usage: "container label key which point to service name",
		},
		cli.StringSliceFlag{
			Name:  flagDockerServiceTagLabels,
			Value: defaultServiceTagLabel,
			Usage: "container label key which point to service tag",
		},
		cli.IntFlag{
			Name:  flagTtl,
			Value: 60,
			Usage: "comparator ttl",
		},
	},
	Action: func(cliContext *cli.Context) error {
		consul := cliContext.String(flagConsul)
		ttl := cliContext.Int(flagTtl)
		appLabels := cliContext.StringSlice(flagDockerServiceNameLabels)
		appTagLabels := cliContext.StringSlice(flagDockerServiceTagLabels)

		theEngine, err := engine.NewEngine(ttl, consul, appLabels, appTagLabels)

		if nil != err {
			return err
		}

		signalC := make(chan os.Signal, 1)
		signal.Notify(signalC, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

		ctx, cancel := context.WithCancel(context.Background())
		theEngine.Start(ctx)

		select {
		case <-signalC:
			cancel()
		}

		return nil
	},
}
