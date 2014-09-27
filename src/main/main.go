package main

import (
	"log"
	"oo"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "oo"
	app.Version = oo.Version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "mode,m",
			Value: "nwrdc",
			Usage: "n:New, w:Write, r:Rename, d:Delete, c:Chmod",
		},
		cli.BoolFlag{
			Name:  "show,s",
			Usage: "show executed commands",
		},
	}
	app.Action = func(c *cli.Context) {
		args := c.Args()
		if len(args) < 2 {
			cli.ShowAppHelp(c)
			return
		}
		o, err := oo.New(c.String("mode"), args[0], args[1], c.Bool("show"))
		if err != nil {
			log.Fatal(err)
		}
		log.Fatal(o.Watch())
	}
	app.Usage = `watch dir and execute command

EXAMPLE:
	# show modified files
	$ oo . "echo {{.}}" 

	# rsync modified files
	$ oo -m="nwrc" -s a "rsync -azvc a/{{.Rel}} b/{{.Rel.Dir}}/" 

	# show deleted files
	$ oo -m="d" . "echo deleted {{.}}"`
	app.Run(os.Args)
}
