package message

import (
	"errors"
	"fmt"
	"time"

	"messager/domain/message"
	"messager/infrastructure/server"
)

type listByStatusRequest struct {
	status string
}

type listByStatusResponse struct {
	Items []listByStatusResponseItem `json:"items"`
}

type listByStatusResponseItem struct {
	ID        string `json:"id,omitempty" example:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"`
	CreatedAt string `json:"createdAt,omitempty" example:"2023-10-27T10:00:00Z"`
	UpdatedAt string `json:"updatedAt,omitempty" example:"2023-10-27T10:00:00Z"`
	Content   string `json:"content,omitempty" example:"Hello from Swagger!"`
	Phone     string `json:"phone,omitempty" example:"+905551234567"`
	Status    string `json:"status,omitempty" example:"PENDING"`
}

// @Summary List messages by status
// @Description Get a list of messages filtered by their status
// @Tags messages
// @Produce json
// @Param status query string true "Message status (e.g., PENDING, SENT)"
// @Success 200 {object} listByStatusResponse
// @Failure 400 {object} server.ErrorResponse "Invalid status parameter"
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /messages [get]
func (h *handler) listByStatus(ctx server.RequestContext) (any, error) {
	request := listByStatusRequest{
		status: ctx.GetQuery("status"),
	}

	messages, err := h.service.ListByStatus(ctx.Context(), message.Status(request.status))
	if errors.Is(err, message.ErrMessageDoesNotValidForListByStatus) {
		return nil, ctx.NewError(server.StatusBadRequest, "Invalid request.", err)
	}
	if err != nil {
		return nil, fmt.Errorf("handler.service.ListByStatus(): %w", err)
	}

	return messagesToListByStatusResponse(messages), nil
}

func messagesToListByStatusResponse(messages []message.Message) *listByStatusResponse {
	response := listByStatusResponse{
		Items: make([]listByStatusResponseItem, 0),
	}

	for _, message := range messages {
		response.Items = append(response.Items, *messageToListByStatusResponseItem(message))
	}

	return &response
}

func messageToListByStatusResponseItem(message message.Message) *listByStatusResponseItem {
	item := listByStatusResponseItem{
		ID:      message.ID,
		Content: message.Content,
		Phone:   message.Phone,
		Status:  string(message.Status),
	}

	if !message.CreatedAt.IsZero() {
		item.CreatedAt = message.CreatedAt.Format(time.RFC3339)
	}

	if !message.UpdatedAt.IsZero() {
		item.UpdatedAt = message.UpdatedAt.Format(time.RFC3339)
	}

	return &item
}
