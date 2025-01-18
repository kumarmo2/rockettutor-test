package main

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Metrics struct {
	Id           string
	Route        string
	Method       string
	StatusCode   uint16
	ResponseTime float64
	CreatedOn    time.Time
}

func NewMetricsFromRawPayload(payload *HttpRouteMetricsPayload) *Metrics {
	return &Metrics{
		Id:           uuid.New().String(),
		Route:        payload.Route,
		StatusCode:   payload.StatusCode,
		ResponseTime: payload.ResponseTime,
		Method:       getHttpMethodString(payload.Method),
		CreatedOn:    time.Now().UTC(),
	}
}

func getHttpMethodString(b byte) string {
	var method string
	switch b {
	case 1:
		method = "GET"
	case 2:
		method = "POST"
	case 3:
		method = "PUT"
	case 4:
		method = "DELETE"
	case 5:
		method = "PATCH"
	case 6:
		method = "OPTIONS"
	case 7:
		method = "HEAD"
	case 8:
		method = "TRACE"
	case 9:
		method = "CONNECT"
	default:
		method = strconv.Itoa(int(b))
	}
	return method

}
