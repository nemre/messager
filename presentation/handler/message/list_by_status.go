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
	ID        string `json:"id,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Content   string `json:"content,omitempty"`
	Phone     string `json:"phone,omitempty"`
	Status    string `json:"status,omitempty"`
}

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
