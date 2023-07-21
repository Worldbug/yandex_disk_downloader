package main

import (
	"net/url"
	"os"
	"strconv"

	"downloader/internal/clients/yandex_disk"
	"downloader/internal/services/disk_storage"
	"downloader/internal/services/downloader"
	"downloader/internal/services/monitoring/cli_monitor"
	"downloader/internal/services/processor"

	"github.com/rs/zerolog"
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

var cli = &cobra.Command{
	Use: "download",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		if len(os.Args) < 3 {
			zerolog.Ctx(ctx).Error().Msg("Not enoughs args ")
			os.Exit(1)
		}

		url, err := url.Parse(os.Args[2])
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Wrong url format")
			os.Exit(1)
		}

		threads, err := strconv.Atoi(os.Args[3])
		if err != nil {
			zerolog.Ctx(ctx).Err(err).Msg("Can`t parse thread count arg")
			os.Exit(1)
		}

		downloader := downloader.NewHTTPDownloader()
		disk_storage := disk_storage.NewDiskStorage(".")
		yandexCli := yandex_disk.NewYandexDiskClient()
		cliMonitor := cli_monitor.NewCliMonitor()

		err = processor.NewProcessor(
			uint(threads),
			downloader,
			disk_storage,
			yandexCli,
			cliMonitor,
		).DownloadDirectory(ctx, url.String())
		if err != nil {
			zerolog.Ctx(ctx).
				Err(err).
				Str("url", url.String()).
				Msg("failed to extract tasks")
			os.Exit(1)
		}
	},
}

var daemon = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {},
}
