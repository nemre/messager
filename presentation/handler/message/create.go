package message

import (
	"errors"
	"fmt"

	"messager/domain/message"
	"messager/infrastructure/server"
)

// @Summary Create a new message
// @Description Create a new message with content and phone number
// @Tags messages
// @Accept json
// @Produce json
// @Param message body createRequest true "Message object to be created"
// @Success 201 {object} createResponse
// @Failure 400 {object} server.ErrorResponse "Invalid request"
// @Failure 500 {object} server.ErrorResponse "Internal server error"
// @Router /messages [post]
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

type createRequest struct {
	Content string `json:"content" example:"Hello, world!"`
	Phone   string `json:"phone" example:"+905551234567"`
}

type createResponse struct {
	ID string `json:"id" example:"a1b2c3d4e5f6g7h8i9j0k1l2m3n4o5p6"`
}
