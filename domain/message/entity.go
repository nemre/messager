package message

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nyaruka/phonenumbers"
)

const (
	StatusPending Status = "PENDING"
	StatusSent    Status = "SENT"

	minContentLength   = 10
	maxContentLength   = 255
	defaultPhoneRegion = "TR"
)

var (
	ErrMessageDoesNotValidForCreate        = errors.New("message does not valid for create")
	ErrMessageDoesNotValidForListByStatus  = errors.New("message does not valid for list by status")
	ErrMessageDoesNotValidForSent          = errors.New("message does not valid for sent")
	ErrMessageNotFound                     = errors.New("message not found")
	ErrMessageStatusDoesNotEligibleForSent = errors.New("message status does not eligible for sent")
)

type Message struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Content   string
	Phone     string
	Status    Status
}

type Status string

func (m *Message) NewErrMessageDoesNotValidForCreate() error {
	return ErrMessageDoesNotValidForCreate
}

func (m *Message) NewErrMessageDoesNotValidForListByStatus() error {
	return ErrMessageDoesNotValidForListByStatus
}

func (m *Message) NewErrMessageDoesNotValidForSent() error {
	return ErrMessageDoesNotValidForSent
}

func (m *Message) NewErrMessageNotFound() error {
	return ErrMessageNotFound
}

func (m *Message) NewErrMessageStatusDoesNotEligibleForSent() error {
	return ErrMessageStatusDoesNotEligibleForSent
}

func (m *Message) ValidateForCreate() error {
	if m.Content == "" {
		return errors.New("message content must be provided")
	}

	if len(m.Content) < minContentLength {
		return fmt.Errorf("message content must be at least %d characters long", minContentLength)
	}

	if len(m.Content) > maxContentLength {
		return fmt.Errorf("message content must not exceed %d characters", maxContentLength)
	}

	if strings.TrimSpace(m.Content) != m.Content {
		return errors.New("message content must not contain leading or trailing whitespace")
	}

	if m.Phone == "" {
		return errors.New("message phone must be provided")
	}

	if strings.TrimSpace(m.Phone) != m.Phone {
		return errors.New("message phone must not contain leading or trailing whitespace")
	}

	if _, err := phonenumbers.Parse(m.Phone, defaultPhoneRegion); err != nil {
		return errors.New("message phone must be a valid phone number")
	}

	if m.Status != StatusPending {
		return errors.New("message status must be pending")
	}

	return nil
}

func (m *Message) ValidateForListByStatus() error {
	if m.Status == "" {
		return errors.New("message status must be provided")
	}

	switch m.Status {
	case StatusPending, StatusSent:
		return nil
	default:
		return errors.New("message status must be one of PENDING or SENT")
	}
}

func (m *Message) ValidateForSent() error {
	if m.ID == "" {
		return errors.New("message id must be provided")
	}

	if err := uuid.Validate(m.ID); err != nil {
		return fmt.Errorf("message id must be a valid uuid: %w", err)
	}

	return nil
}
