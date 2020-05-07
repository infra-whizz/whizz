package main

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"

	"github.com/infra-whizz/whizz"
	whizz_cli "github.com/infra-whizz/whizz/cli"
	wzlib_utils "github.com/infra-whizz/wzlib/utils"
	"github.com/isbm/go-nanoconf"
	"github.com/urfave/cli/v2"
)

func prepareLogger(client *whizz.WzClient, ctx *cli.Context) {
	switch ctx.String("log") {
	case "quiet":
		client.MuteLogger()
	case "trace":
	default:
		fmt.Printf("Unknown logger option: %s\n", ctx.String("log"))
		os.Exit(wzlib_utils.EX_GENERIC)
	}
}

func runner(ctx *cli.Context) error {
	client := whizz.NewWhizzClient()
	prepareLogger(client, ctx)

	client.Call()

	return nil
}

// managePKI manages the PKI
func managePki(ctx *cli.Context) error {
	conf := nanoconf.NewConfig(ctx.String("config"))
	client := whizz.NewWhizzClient()
	prepareLogger(client, ctx)

	if ctx.Bool("generate") {
		pkiDir := conf.Root().String("pki-path", "")
		client.GetLogger().Debugf("Generating PKI keys into %s", pkiDir)
		if err := client.GetCryptoBundle().GetRSA().GenerateKeyPair(pkiDir); err != nil {
			msg := "Error generating PKI: %s"
			if client.GetLogger().GetLevel() != logrus.TraceLevel {
				fmt.Printf(msg, err.Error())
			} else {
				client.GetLogger().Errorf(msg, err.Error())
			}
		}
	}

	return nil
}

// loadPKI loads default keypair location
func loadPKI(client *whizz.WzClient, pkiDir string) {
	if err := client.GetCryptoBundle().GetRSA().LoadPEMKeyPair(pkiDir); err != nil {
		client.GetLogger().Errorf("Error loading PKI keys: %s", err.Error())
		os.Exit(wzlib_utils.EX_GENERIC)
	}
}

// run the client
func client(ctx *cli.Context) error {
	conf := nanoconf.NewConfig(ctx.String("config"))
	client := whizz.NewWhizzClient()
	prepareLogger(client, ctx)

	loadPKI(client, conf.Root().String("pki-path", ""))

	if ctx.String("client") == "accept" && (ctx.Bool("all") || len(ctx.StringSlice("finger")) > 0) {
		client.Boot()
		defer client.Stop()
		missing := client.Accept(ctx.StringSlice("finger")...)
		if len(missing) > 0 {
			fmt.Println("Following fingerprints as new systems was not found:")
			for idx, fp := range missing {
				fmt.Printf("%d. %s\n", idx+1, fp)
			}
		}
	} else if ctx.String("client") == "reject" && (ctx.Bool("all") || len(ctx.StringSlice("finger")) > 0) {
		client.Boot()
		defer client.Stop()
		client.Reject(ctx.StringSlice("finger")...)
	} else if ctx.String("client") == "delete" && len(ctx.StringSlice("finger")) > 0 {
		client.Boot()
		defer client.Stop()
		missing := client.Delete(ctx.StringSlice("finger")...)
		if len(missing) > 0 {
			fmt.Println("Following fingerprints as deleted systems was not found:")
			for idx, fp := range missing {
				fmt.Printf("%d. %s\n", idx+1, fp)
			}
		}
	} else if ctx.String("list") == "new" || ctx.String("list") == "rejected" {
		client.Boot()
		defer client.Stop()
		fmtr := whizz_cli.NewWhizzCliFormatter()
		switch ctx.String("list") {
		case "new":
			clients := client.ListNew()
			client.GetLogger().Debugf("Found %d new client(s)", len(clients))
			for idx, clientData := range clients {
				fmtr.HostnameWithFp(idx+1, clientData["Fqdn"].(string), clientData["RsaFp"].(string))
			}
		case "rejected":
			clients := client.ListRejected()
			client.GetLogger().Debugf("Found %d rejected client(s)", len(clients))
			for idx, clientData := range clients {
				fmtr.HostnameWithFp(idx+1, clientData["Fqdn"].(string), clientData["RsaFp"].(string))
			}
		}
	} else if ctx.String("search") != "" {
		client.Boot()
		defer client.Stop()
		client.GetLogger().Debugln("Searching for", ctx.String("search"), "host[s]")
		whizz_cli.NewWhizzCliFormatter().ListSystems(client.Search(ctx.String("search")))
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
			Name:   "pki",
			Usage:  "Manage public key infrastructure",
			Action: managePki,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:    "generate",
					Usage:   "Generate keypair. NOTE: this is possible only if no keys are present.",
					Aliases: []string{"g"},
				},
			},
		},
		{
			Name:   "client",
			Usage:  "Operations with the clients",
			Action: client,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "client",
					Usage:   "Do (accept|reject|delete) clients. Requires --all or --finger. NOTE: --all does not apply to 'delete'",
					Aliases: []string{"c"},
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
				&cli.StringFlag{
					Name:    "search",
					Usage:   "Search clients by hostname. Wildcards supported.",
					Aliases: []string{"s"},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
