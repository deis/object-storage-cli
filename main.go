package main

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/deis/object-storage-cli/actions"
	"github.com/deis/object-storage-cli/config"
)

var (
	debugStr = "DEIS_DEBUG"
)

func main() {
	var debug bool
	debugValStr := os.Getenv(debugStr)
	if debugValStr == "true" || debugValStr == "1" {
		debug = true
	}

	app := cli.NewApp()
	app.Name = "objstorage"
	app.Usage = "Use a variety of different object storage systems with a single tool"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  config.StorageTypeFlag,
			Value: config.S3StorageType.String(),
			Usage: "Specify the type of the object storage system",
		},
	}
	app.Commands = []cli.Command{
		actions.DownloadCommand(debug),
		actions.UploadCommand(debug),
	}
	app.Run(os.Args)
}
