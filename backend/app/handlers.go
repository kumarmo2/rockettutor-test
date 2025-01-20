package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
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

func (h *Handlers) GetRouteMetricsForEndpointHandler(w http.ResponseWriter, r *http.Request) {
	route := r.URL.Query().Get("route")
	latestWindowInMinutes := r.URL.Query().Get("latestWindowInMinutes")

	if route == "" {
		msg := "Empty route not valid"
		newErrResult[any](&msg).WriteHttpResponse(w, 400)
		return
	}
	if latestWindowInMinutes == "" {
		msg := "Empty latestWindowInMinutes not valid"
		newErrResult[any](&msg).WriteHttpResponse(w, 400)
		return
	}
	latestWindowInMinutesInt, err := strconv.Atoi(latestWindowInMinutes)
	if err != nil {
		msg := "Invalid latestWindowInMinutes not valid"
		newErrResult[any](&msg).WriteHttpResponse(w, 400)
		return
	}
	fromTime := time.Now().UTC().Add(-time.Minute * time.Duration(latestWindowInMinutesInt))
	metrics, err := h.server.metricsDao.GetMetricsDataForEndpoint(route, fromTime)
	if err != nil {
		fmt.Println("error while getting metrics: ", err)
		msg := "Internal server error. Please try again!!"
		newErrResult[any](&msg).WriteHttpResponse(w, 500)
		return
	}

	fmt.Println("route:", route, "latestWindowInMinutes:", latestWindowInMinutes)
	// res := map[string]string{"k1": "v1"}
	newOkResult[[]Metrics, any](&metrics).WriteHttpResponse(w, 200)
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

	endClientChan := make(chan bool, 1)

	go func() {
		for {
			if _, _, err := conn.NextReader(); err != nil {
				endClientChan <- true
				break
			}
		}
	}()

	go func() {
		for {
			msg := <-msgChan
			r := newOkResult[Metrics, any](&msg)
			bytes, _ := json.Marshal(&r)
			// fmt.Println("writing data to client")
			conn.WriteMessage(websocket.TextMessage, bytes)
		}
	}()

	<-endClientChan

	fmt.Println("HandleWebSockets ending")
}
