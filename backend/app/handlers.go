package main

import (
	"encoding/json"
	"net/http"

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

	msgChan := make(chan Metrics)
	sub, err := h.server.liveEventManager.register(msgChan)

	if err != nil {
		msg := "Internal Server Error. Please try again"
		newErrResult[any](&msg).WriteHttpResponse(w, 500)
		return
	}
	defer h.server.liveEventManager.deregister(sub)
	for {
		msg := <-msgChan
		// fmt.Println("time: ", msg.ResponseTime)

		r := newOkResult[Metrics, any](&msg)
		bytes, _ := json.Marshal(&r)
		conn.WriteMessage(websocket.TextMessage, bytes)
	}
}
