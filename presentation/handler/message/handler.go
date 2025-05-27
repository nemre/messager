package message

import (
	"messager/domain/message"
	service "messager/domain/message"
	"messager/infrastructure/server"
	job "messager/presentation/job/message"
)

type Handler interface {
	create(ctx server.RequestContext) (any, error)
}

type handler struct {
	service service.Service
	job     job.Job
}

func New(router server.Router, service message.Service, job job.Job) Handler {
	h := handler{
		service: service,
		job:     job,
	}

	router.AddRoute("POST /messages", h.create)
	router.AddRoute("GET /messages", h.listByStatus)
	router.AddRoute("POST /messages/jobs", h.startJob)
	router.AddRoute("DELETE /messages/jobs", h.stopJob)

	return &h
}
