package actions

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/deis/object-storage-cli/config"
	"github.com/docker/distribution/context"
)

// DownloadCommand returns the cli.Command for use in the top level app commands list
func DownloadCommand(debug bool) cli.Command {
	return cli.Command{
		Name:      "download",
		ShortName: "dl",
		Action:    Download(debug),
	}
}

// Download is the cli handler for "download" command
func Download(debug bool) func(c *cli.Context) {
	return func(c *cli.Context) {
		args := c.Args()
		if len(args) < 2 {
			log.Fatalf("This command should be called as 'download $REMOTE_PATH $LOCAL_PATH'")
		}

		remote := args[0]
		local := args[1]

		conf, err := config.FromStorageTypeString(c.GlobalString(config.StorageTypeFlag))
		if err != nil {
			log.Fatal(err)
		}
		driver, err := conf.CreateDriver()
		if err != nil {
			log.Fatal(err)
		}

		ctx := context.Background()
		b, err := driver.GetContent(ctx, remote)
		if err != nil {
			log.Fatalf("Error downloading %s (%s)", remote, err)
		}
		if err := ioutil.WriteFile(local, b, os.ModePerm); err != nil {
			log.Fatalf("Error writing %s to %s (%s)", remote, local, err)
		}

		if debug {
			log.Printf("Successfully copied %s to %s", remote, local)
		}
	}
}
