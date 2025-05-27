package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
)

type Router interface {
	AddRoute(pattern string, handler func(ctx RequestContext) (any, error))
}

type router struct {
	mux            *http.ServeMux
	idHeader       string
	onRequestStart func(ctx RequestContext)
	onRequestEnd   func(ctx RequestContext, status uint16)
	onRequestError func(ctx RequestContext, status uint16, message string, err error)
	onRequestPanic func(ctx RequestContext, status uint16, message string, err error, stackTrace string)
}

func (r *router) AddRoute(pattern string, handler func(ctx RequestContext) (any, error)) {
	r.mux.HandleFunc(pattern, r.wrapHandler(handler))
}

func (r *router) wrapHandler(handler func(ctx RequestContext) (any, error)) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, request *http.Request) {
		ctx := newRequestContext(responseWriter, request, r.idHeader)
		setDefaultHeaders(responseWriter, r.idHeader, ctx)
		defer r.recoverFromPanic(responseWriter, ctx)
		r.handleRequest(responseWriter, ctx, handler)
	}
}

func setDefaultHeaders(responseWriter http.ResponseWriter, idHeader string, ctx RequestContext) {
	responseWriter.Header().Set("Content-Type", "application/json")
	if idHeader != "" {
		responseWriter.Header().Set(idHeader, ctx.GetID())
	}
}

func (r *router) recoverFromPanic(responseWriter http.ResponseWriter, ctx RequestContext) {
	if rec := recover(); rec != nil {
		status := http.StatusInternalServerError
		message := "An unexpected error occurred."
		err := fmt.Errorf("%v", rec)

		_ = json.NewEncoder(responseWriter).Encode(map[string]string{
			"message": message,
			"error":   err.Error(),
		})

		if r.onRequestPanic != nil {
			r.onRequestPanic(ctx, uint16(status), message, err, string(debug.Stack()))
		}
	}
}

func (r *router) handleRequest(responseWriter http.ResponseWriter, ctx RequestContext, handler func(ctx RequestContext) (any, error)) {
	response, err := handler(ctx)
	if err != nil {
		r.handleError(responseWriter, ctx, err)

		return
	}

	_ = json.NewEncoder(responseWriter).Encode(response)

	if r.onRequestEnd != nil {
		r.onRequestEnd(ctx, http.StatusOK)
	}
}

func (r *router) handleError(responseWriter http.ResponseWriter, ctx RequestContext, err error) {
	status := http.StatusInternalServerError
	message := "An unexpected error occurred."

	var re *requestError
	if errors.As(err, &re) {
		status = int(re.status)
		message = re.message
	}

	responseWriter.WriteHeader(status)

	_ = json.NewEncoder(responseWriter).Encode(map[string]string{
		"message": message,
		"error":   err.Error(),
	})

	if r.onRequestError != nil {
		r.onRequestError(ctx, uint16(status), message, err)
	}
}
