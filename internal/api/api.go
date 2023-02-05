package api

import (
	"github.com/gin-gonic/gin"
	"github.com/romeq/usva/internal/db"
)

type Server struct {
	engine            *gin.Engine
	db                *db.Queries
	api               *Configuration
	encryptionKeySize uint32
}

func NewServer(eng *gin.Engine, conn *db.Queries, apic *Configuration, encs uint32) *Server {
	return &Server{
		engine:            eng,
		db:                conn,
		api:               apic,
		encryptionKeySize: encs,
	}
}

func (s *Server) GetRouter() *gin.Engine {
	return s.engine
}
