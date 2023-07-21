package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

type createTaskRequest struct {
	URL          string `json:"url"`
	ThreadsCount int    `json:"threads_count"`
}

func (api *API) createTask(ctx *gin.Context) {
	req := &createTaskRequest{}
	if err := ctx.ShouldBindJSON(req); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err := api.downloader.AddTask(ctx, uint(req.ThreadsCount), req.URL)
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		zerolog.Ctx(ctx).
			Err(err).
			Str("url", req.URL).
			Str("client", ctx.ClientIP()).
			Msg("failed to extract tasks")
		return
	}

	ctx.AbortWithStatus(http.StatusCreated)
}
