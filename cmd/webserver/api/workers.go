package api

import (
	"context"

	"github.com/romeq/usva/internal/workers"
)

func (s *Server) IncludeServerContextWorker(w workers.Worker) {
	w.Run(context.Background())
}
