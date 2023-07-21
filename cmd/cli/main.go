package main

import (
	"github.com/spf13/cobra"

	_ "downloader/internal/logger"
)

func main() {
	app := cobra.Command{
		Use: "downloader",
	}

	app.AddCommand(cli)
	app.AddCommand(daemon)

	app.Execute()
}
