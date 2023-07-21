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
) *API {
	return &API{
		downloader: downloader,
		monitor:    monitor,
		router:     router,
	}
}

type API struct {
	downloader Downloader
	monitor    Monitor

	router *gin.Engine
}

func (api *API) Run(ctx context.Context) {
	apiHandlers := api.router.Group("/api")
	apiHandlers.POST("/create_task", api.createTask)
	apiHandlers.GET("/task_status", api.taskStatus)
}
