package server

import "book_keeper/internal/health"

type Handlers struct {
	HealthHandler *health.Handler
}

func (s *Server) InitRoutes(h Handlers) {
	router := s.routerGroups.rootRouter
	router.GET("/sanity", h.HealthHandler.CheckSanity)

}
