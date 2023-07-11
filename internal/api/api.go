package api

import (
	"context"

	"github.com/gin-gonic/gin"
)

type Downloader interface {
	AddTask(ctx context.Context, url string) error
}

type Monitor interface {
	StatusStream(ctx context.Context) // TODO: stream type ?
}

func NewAPI(
	downloader Downloader,
	monitor Monitor,
) *API {
	return &API{
		downloader: downloader,
		monitor:    monitor,
	}
}

type API struct {
	downloader Downloader
	monitor    Monitor
}

func (api *API) CreateJob(ctx *gin.Context)          {}
func (api *API) GetJobStatusStream(ctx *gin.Context) {}
func (api *API) GetJobsList(ctx *gin.Context)        {}
