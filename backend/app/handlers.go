package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type Handlers struct {
	server     *Server
	wsUpgrader websocket.Upgrader
}

func newHandlers(server *Server) *Handlers {
	return &Handlers{server: server, wsUpgrader: websocket.Upgrader{}}
}

func (h *Handlers) HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func (h *Handlers) GetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	res := map[string]string{"k1": "v1"}
	newOkResult[map[string]string, any](&res).WriteHttpResponse(w, 200)
}

func (h *Handlers) HandleWebSockets(w http.ResponseWriter, r *http.Request) {
	conn, err := h.wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	defer conn.Close()
	for {
		time.Sleep(1 * time.Second)
		conn.WriteMessage(websocket.TextMessage, []byte("hello"))
	}
}
