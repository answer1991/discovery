package main

import (
	"github.com/answer1991/daily-roll-logrus"
	"github.com/urfave/cli"
	"os"
)

var (
	usage  = ``
	logger = drl.GetLogger("registrator")
)

func main() {
	app := cli.NewApp()
	app.Name = "registrator"
	app.Usage = usage
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug output for logging",
		},
	}
	app.Before = func(context *cli.Context) error {
		debug := context.GlobalBool("debug")
		drl.SetEnableStdout(debug)

		//if debug {
		//	drl.SetLevel(logrus.Level(logrus.DebugLevel))
		//} else {
		//	drl.SetLevel(logrus.InfoLevel)
		//}

		return nil
	}
	app.Commands = []cli.Command{
		serveCommand,
	}

	if err := app.Run(os.Args); err != nil {
		logger.WithField("error", err).Info("App Exist with Error")
	}
}
