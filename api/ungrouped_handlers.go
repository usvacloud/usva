package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

var errEmptyResponse = errors.New("content was not found")

func (s *Server) RestrictionsHandler(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"maxSize": s.api.MaxSingleUploadSize,
	})
}

func (s *Server) NotFoundHandler(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusNotFound)
	ctx.File("public/404.html")
}
