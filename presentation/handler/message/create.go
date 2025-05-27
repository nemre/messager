package message

import (
	"errors"
	"fmt"

	"messager/domain/message"
	"messager/infrastructure/server"
)

type createRequest struct {
	Content string `json:"content"`
	Phone   string `json:"phone"`
}

type createResponse struct {
	ID string `json:"id"`
}

func (h *handler) create(ctx server.RequestContext) (any, error) {
	var request createRequest

	if err := ctx.ParseJSONBody(&request); err != nil {
		return nil, ctx.NewError(server.StatusBadRequest, "Invalid request.", err)
	}

	newMessage, err := h.service.Create(ctx.Context(), request.toMessage())
	if errors.Is(err, message.ErrMessageDoesNotValidForCreate) {
		return nil, ctx.NewError(server.StatusBadRequest, "Invalid request.", err)
	}
	if err != nil {
		return nil, fmt.Errorf("handler.service.Create(): %w", err)
	}

	return messageToCreateResponse(*newMessage), nil
}

func (l *createRequest) toMessage() message.Message {
	return message.Message{
		Content: l.Content,
		Phone:   l.Phone,
		Status:  message.StatusPending,
	}
}

func messageToCreateResponse(message message.Message) *createResponse {
	return &createResponse{
		ID: message.ID,
	}
}
