package main

import "net/http"

type Server struct {
	addr             string
	liveEventManager *LiveEventManager[Metrics]
	metricsDao       IMetricsDao
}

func newServer(addr string, liveEventManager *LiveEventManager[Metrics], metricsDao IMetricsDao) *Server {
	return &Server{addr: addr, liveEventManager: liveEventManager, metricsDao: metricsDao}
}

func (s *Server) run(doneChan chan<- bool) {
	go func() {
		defer func() {
			doneChan <- true
		}()
		mux := http.NewServeMux()
		handler := newHandlers(s)

		mux.HandleFunc("/api/ws", handler.HandleWebSockets)
		mux.HandleFunc("GET /api/metrics", handler.GetRouteMetricsForEndpointHandler)

		http.ListenAndServe(s.addr, mux)
	}()
}
