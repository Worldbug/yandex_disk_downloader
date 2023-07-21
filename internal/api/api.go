package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Downloader interface {
	AddTask(ctx context.Context, threads uint, url string) error
}

type Monitor interface {
	StatusStream(ctx context.Context) // TODO: stream type ?
}

func NewAPI(
	downloader Downloader,
	monitor Monitor,
	router *gin.Engine,
	host string,
) *API {
	return &API{
		downloader: downloader,
		monitor:    monitor,
		router:     router,
		host:       host,
	}
}

type API struct {
	downloader Downloader
	monitor    Monitor

	router *gin.Engine
	host   string
}

func (api *API) Run(ctx context.Context) error {
	apiHandlers := api.router.Group("/api")

	apiHandlers.POST("/create_task", api.createTask)
	apiHandlers.GET("/task_status", api.taskStatus)

	return api.router.Run(api.host)
}
