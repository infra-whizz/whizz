package main

import (
	"os"

	whizz_cli "github.com/infra-whizz/whizz/cli"

	"github.com/infra-whizz/whizz"
	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
)

func runner(ctx *cli.Context) error {
	client := whizz.NewWhizzClient()
	client.Call()

	return nil
}

func client(ctx *cli.Context) error {
	client := whizz.NewWhizzClient()
	if ctx.Bool("accept") && (ctx.Bool("all") || len(ctx.StringSlice("finger")) > 0) {
		client.Boot()
		defer client.Stop()

		client.Accept(ctx.StringSlice("finger")...)
	} else if ctx.Bool("reject") && (ctx.Bool("all") || len(ctx.StringSlice("finger")) > 0) {
		client.Boot()
		defer client.Stop()

		client.Reject(ctx.StringSlice("finger")...)
	} else if ctx.String("list") == "new" || ctx.String("list") == "rejected" {
		client.Boot()
		defer client.Stop()
		switch ctx.String("list") {
		case "new":
			fmtr := whizz_cli.NewWhizzCliFormatter()
			for idx, clientData := range client.ListNew() {
				fmtr.HostnameWithFp(idx+1, clientData["Fqdn"].(string), clientData["RsaFp"].(string))
			}
		case "rejected":
			client.ListRejected()
		}
	} else {
		return cli.ShowSubcommandHelp(ctx)
	}

	return nil
}

func main() {
	appname := "whizz"
	confpath := nanoconf.NewNanoconfFinder(appname).DefaultSetup(nil)
	app := &cli.App{
		Version: "0.1 Alpha",
		Name:    appname,
		Usage:   "Ansible on Steroids",
		Action:  runner,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "log",
				Aliases: []string{"l"},
				Usage:   "Set logging level. Choices: 'quiet' or 'trace'.",
				Value:   "quiet",
			},
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Path to configuration file",
				Required: false,
				Value:    confpath.SetDefaultConfig(confpath.FindFirst()).FindDefault(),
			},
		},
	}
	app.Commands = []*cli.Command{
		{
			Name:   "client",
			Usage:  "Operations with the clients",
			Action: client,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "accept",
					Usage:   "Accept new clients. Requires --all or --finger.",
					Aliases: []string{"c"},
				},
				&cli.BoolFlag{
					Name:    "reject",
					Usage:   "Reject new clients. Requires --all or --finger.",
					Aliases: []string{"r"},
				},
				&cli.BoolFlag{
					Name:    "all",
					Usage:   "Mark all",
					Aliases: []string{"a"},
				},
				&cli.StringSliceFlag{
					Name:    "finger",
					Usage:   "Fingerprint (or part of it) that matches the client",
					Aliases: []string{"f"},
				},
				&cli.StringFlag{
					Name:    "list",
					Usage:   "List clients. Choices: (new|rejected)",
					Aliases: []string{"l"},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
