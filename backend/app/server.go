package main

import "net/http"

type Server struct {
	addr string
}

func newServer(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) run(doneChan chan<- bool) {
	go func() {
		defer func() {
			doneChan <- true
		}()
		mux := http.NewServeMux()
		handler := newHandlers(s)

		mux.HandleFunc("/api/ws", handler.HandleWebSockets)
		mux.HandleFunc("GET /api/metrics", handler.GetMetricsHandler)

		http.ListenAndServe(s.addr, mux)
	}()
}
