package main

import (
	"os"

	"downloader/internal/api"
	"downloader/internal/clients/yandex_disk"
	"downloader/internal/services/disk_storage"
	"downloader/internal/services/downloader"
	"downloader/internal/services/monitoring/api_monitor"
	"downloader/internal/services/processor"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var daemon = &cobra.Command{
	Use: "run",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := cmd.Context()

		downloader := downloader.NewHTTPDownloader()
		disk_storage := disk_storage.NewDiskStorage(".")
		yandexCli := yandex_disk.NewYandexDiskClient()

		monitor := api_monitor.NewStreamMonitor()

		processor := processor.NewProcessor(
			0,
			downloader,
			disk_storage,
			yandexCli,
			monitor,
		)

		host := ""

		apiHandler := api.NewAPI(
			processor,
			monitor,
			gin.Default(),
			host,
		)

		if err := apiHandler.Run(ctx); err != nil {
			zerolog.Ctx(ctx).
				Err(err).
				Msg("Failed to run API")
			os.Exit(1)
		}
	},
}
