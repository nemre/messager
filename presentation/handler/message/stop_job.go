package message

import (
	"messager/infrastructure/server"
)

type stopJobResponse struct {
	Stopped bool `json:"stopped" example:true`
}

// @Summary Stop the message sending job
// @Description Stops the background job that sends pending messages
// @Tags messages
// @Produce json
// @Success 200 {object} stopJobResponse
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /messages/jobs [delete]
func (h *handler) stopJob(ctx server.RequestContext) (any, error) {
	h.job.Stop()

	return stopJobResponse{
		Stopped: true,
	}, nil
}
