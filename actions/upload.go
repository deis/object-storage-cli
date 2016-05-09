package actions

import (
	"io/ioutil"
	"log"

	"github.com/codegangsta/cli"
	"github.com/deis/object-storage-cli/config"
	"github.com/docker/distribution/context"
)

// UploadCommand returns the cli.Command for use in the top level app commands list
func UploadCommand(debug bool) cli.Command {
	return cli.Command{
		Name:      "upload",
		ShortName: "up",
		Action:    Upload(debug),
	}
}

// Upload is the cli handler for "upload" command
func Upload(debug bool) func(c *cli.Context) {
	return func(c *cli.Context) {
		args := c.Args()
		if len(args) < 2 {
			log.Fatalf("This command should be called as 'upload $LOCAL_PATH $REMOTE_PATH'")
		}
		local := args[0]
		remote := args[1]

		conf, err := config.FromStorageTypeString(c.GlobalString(config.StorageTypeFlag))
		if err != nil {
			log.Fatal(err)
		}
		driver, err := conf.CreateDriver()
		if err != nil {
			log.Fatal(err)
		}

		fBytes, err := ioutil.ReadFile(local)
		if err != nil {
			log.Fatalf("Error reading local file %s (%s)", local, err)
		}
		ctx := context.Background()
		if err := driver.PutContent(ctx, remote, fBytes); err != nil {
			log.Fatalf("Error writing remote object %s (%s)", remote, err)
		}
		if debug {
			log.Printf("Successfully copied %s to %s", remote, local)
		}
	}
}
