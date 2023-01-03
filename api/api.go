package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/db"
)

var CookieDomain string

type Server struct {
	engine *gin.Engine
	db     *db.Queries
	api    *Configuration
}

func NewServer(engine *gin.Engine, conn *db.Queries, apiconfiguration *Configuration) *Server {
	return &Server{
		engine: engine,
		api:    apiconfiguration,
		db:     conn,
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.engine
}
