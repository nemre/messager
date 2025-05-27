package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

const (
	StatusBadRequest uint16 = 400
	StatusNotFound   uint16 = 404
)

type RequestContext interface {
	Context() context.Context
	NewError(status uint16, message string, err error) error
	GetQuery(key string) string
	GetMethod() string
	GetURL() string
	GetID() string
	GetPathValue(key string) string
	ParseJSONBody(object any) error
}

// ErrorResponse represents the structure of an error response.
type ErrorResponse struct {
	Message string `json:"message" example:"Invalid request."`
	Error   string `json:"error" example:"details of the error"`
}

type requestContext struct {
	responseWriter http.ResponseWriter
	request        *http.Request
	id             string
}

type requestError struct {
	status  uint16
	message string
	err     error
}

func newRequestContext(responseWriter http.ResponseWriter, request *http.Request, idHeader string) RequestContext {
	id := request.Header.Get(idHeader)

	if id == "" {
		id = uuid.New().String()
	}

	return &requestContext{
		responseWriter: responseWriter,
		request:        request,
		id:             id,
	}
}

func (r *requestContext) Context() context.Context {
	return r.request.Context()
}

func (r *requestContext) NewError(status uint16, message string, err error) error {
	return &requestError{
		status:  status,
		message: message,
		err:     err,
	}
}

func (r *requestContext) GetQuery(key string) string {
	return r.request.URL.Query().Get(key)
}

func (r *requestContext) GetMethod() string {
	return r.request.Method
}

func (r *requestContext) GetURL() string {
	return r.request.URL.String()
}

func (r *requestContext) GetID() string {
	return r.id
}

func (r *requestContext) GetPathValue(key string) string {
	return r.request.PathValue(key)
}

func (r *requestContext) ParseJSONBody(object any) error {
	err := json.NewDecoder(r.request.Body).Decode(object)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("json.Decoder.Decode(): %w", err)
	}

	return nil
}

func (r *requestError) Error() string {
	if r.err == nil {
		return r.message
	}

	return r.err.Error()
}
