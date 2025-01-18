package main

import (
	"encoding/json"
	"net/http"
)

type ApiResult[T any, E any] struct {
	Ok  *T `json:"ok,omitempty"`
	Err *E `json:"err,omitempty"`
}

func newOkResult[T any, E any](ok *T) *ApiResult[T, E] {
	return &ApiResult[T, E]{Ok: ok}
}

func newErrResult[T any, E any](err *E) *ApiResult[T, E] {
	return &ApiResult[T, E]{Err: err}
}

func (res *ApiResult[T, E]) WriteHttpResponse(w http.ResponseWriter, status int) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(*res)
}
