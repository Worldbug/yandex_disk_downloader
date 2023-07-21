package main

import (
	"net/url"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"downloader/internal/clients/yandex_disk"
	_ "downloader/internal/logger"
	"downloader/internal/services/disk_storage"
	"downloader/internal/services/downloader"
	"downloader/internal/services/monitoring/cli_monitor"
	"downloader/internal/services/processor"
)

func main() {
	app := cobra.Command{
		Use: "downloader",
	}

	app.AddCommand(newDownloadCMD())
	app.AddCommand(newDaemonCMD())

	if err := app.Execute(); err != nil {
		os.Exit(1)
	}
}

func newDownloadCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "download",
		Short: "downloader download [Yandex disk URL] [Threads count]",
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

			proc := processor.NewProcessor(
				downloader,
				disk_storage,
				yandexCli,
				cliMonitor,
			)

			if err := proc.DownloadDirectory(ctx, uint(threads), url.String()); err != nil {
				zerolog.Ctx(ctx).
					Err(err).
					Str("url", url.String()).
					Msg("failed to extract tasks")
				os.Exit(1)
			}
		},
	}
}

func newDaemonCMD() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "downloader run",
		Run: func(cmd *cobra.Command, args []string) {
			// ctx := cmd.Context()

			// downloader := downloader.NewHTTPDownloader()
			// disk_storage := disk_storage.NewDiskStorage(".")
			// yandexCli := yandex_disk.NewYandexDiskClient()

			// monitor := api_monitor.NewStreamMonitor()

			// processor := processor.NewProcessor(
			// 	0,
			// 	downloader,
			// 	disk_storage,
			// 	yandexCli,
			// 	monitor,
			// )

			// host := ""

			// apiHandler := api.NewAPI(
			// 	processor,
			// 	monitor,
			// 	gin.Default(),
			// 	host,
			// )

			// if err := apiHandler.Run(ctx); err != nil {
			// 	zerolog.Ctx(ctx).
			// 		Err(err).
			// 		Msg("Failed to run API")
			// 	os.Exit(1)
			// }
		},
	}
}
