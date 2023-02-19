package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/internal/generated/db"
)

type Server struct {
	GinEngine  *gin.Engine
	DB         *db.Queries
	Config     *Configuration
	EncKeySize uint32
}

func NewServer(eng *gin.Engine, conn *db.Queries, apic *Configuration, encs uint32) *Server {
	return &Server{
		GinEngine:  eng,
		DB:         conn,
		Config:     apic,
		EncKeySize: encs,
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.GinEngine
}
