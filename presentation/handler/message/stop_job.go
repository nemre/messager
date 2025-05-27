package message

import (
	"messager/infrastructure/server"
)

type stopJobResponse struct {
	Stopped bool `json:"stopped"`
}

func (h *handler) stopJob(ctx server.RequestContext) (any, error) {
	h.job.Stop()

	return stopJobResponse{
		Stopped: true,
	}, nil
}
