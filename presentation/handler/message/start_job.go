package message

import (
	"messager/infrastructure/server"
)

type startJobResponse struct {
	Started bool `json:"started"`
}

func (h *handler) startJob(ctx server.RequestContext) (any, error) {
	h.job.Start()

	return startJobResponse{
		Started: true,
	}, nil
}
