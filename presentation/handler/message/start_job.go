package message

import (
	"messager/infrastructure/server"
)

type startJobResponse struct {
	Started bool `json:"started" example:true`
}

// @Summary Start the message sending job
// @Description Starts the background job that sends pending messages
// @Tags messages
// @Produce json
// @Success 200 {object} startJobResponse
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /messages/jobs [post]
func (h *handler) startJob(ctx server.RequestContext) (any, error) {
	h.job.Start()

	return startJobResponse{
		Started: true,
	}, nil
}
