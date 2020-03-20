package main

import (
	"os"

	"github.com/infra-whizz/whizz"
	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
)

func run(ctx *cli.Context) error {
	client := whizz.NewWhizzClient()
	client.Call()

	return nil
}

func main() {
	appname := "whizz"
	confpath := nanoconf.NewNanoconfFinder(appname).DefaultSetup(nil)
	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    appname,
		Usage:   "Ansible on Steroids",
		Action:  run,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to configuration file",
				Required: false,
				Value:    confpath.SetDefaultConfig(confpath.FindFirst()).FindDefault(),
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
