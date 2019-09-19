package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli"
)

var version = time.Now().String()

func main() {
	if err := newApp().Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}

func newApp() *cli.App {
	app := cli.NewApp()
	app.Version = version
	app.EnableBashCompletion = true
	app.Name = "gsuite"
	app.Usage = "Google G Suite command line tool"
	// override -v
	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "verbose logging",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:  "user",
			Usage: "Show list of all users",
			Action: func(c *cli.Context) error {
				return cmdUserList(c)
			},
			ArgsUsage: `user`,
		},
	}
	return app
}
