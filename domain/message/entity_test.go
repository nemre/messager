package message

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMessage_ValidateForCreate(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid message",
			message: Message{
				Content: "This is a valid message content that meets the minimum length requirement",
				Phone:   "+905551234567",
				Status:  StatusPending,
			},
			wantErr: false,
		},
		{
			name: "empty content",
			message: Message{
				Content: "",
				Phone:   "+905551234567",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message content must be provided",
		},
		{
			name: "content too short",
			message: Message{
				Content: "short",
				Phone:   "+905551234567",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message content must be at least 10 characters long",
		},
		{
			name: "content too long",
			message: Message{
				Content: string(make([]byte, maxContentLength+1)),
				Phone:   "+905551234567",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message content must not exceed 255 characters",
		},
		{
			name: "content with leading whitespace",
			message: Message{
				Content: "  content with leading whitespace",
				Phone:   "+905551234567",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message content must not contain leading or trailing whitespace",
		},
		{
			name: "content with trailing whitespace",
			message: Message{
				Content: "content with trailing whitespace  ",
				Phone:   "+905551234567",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message content must not contain leading or trailing whitespace",
		},
		{
			name: "content with both leading and trailing whitespace",
			message: Message{
				Content: "  content with both leading and trailing whitespace  ",
				Phone:   "+905551234567",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message content must not contain leading or trailing whitespace",
		},
		{
			name: "empty phone",
			message: Message{
				Content: "This is a valid message content",
				Phone:   "",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message phone must be provided",
		},
		{
			name: "phone with leading whitespace",
			message: Message{
				Content: "This is a valid message content",
				Phone:   " +905551234567",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message phone must not contain leading or trailing whitespace",
		},
		{
			name: "phone with trailing whitespace",
			message: Message{
				Content: "This is a valid message content",
				Phone:   "+905551234567 ",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message phone must not contain leading or trailing whitespace",
		},
		{
			name: "phone with both leading and trailing whitespace",
			message: Message{
				Content: "This is a valid message content",
				Phone:   " +905551234567 ",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message phone must not contain leading or trailing whitespace",
		},
		{
			name: "invalid phone number",
			message: Message{
				Content: "This is a valid message content",
				Phone:   "invalid-phone",
				Status:  StatusPending,
			},
			wantErr: true,
			errMsg:  "message phone must be a valid phone number",
		},
		{
			name: "invalid status",
			message: Message{
				Content: "This is a valid message content",
				Phone:   "+905551234567",
				Status:  StatusSent,
			},
			wantErr: true,
			errMsg:  "message status must be pending",
		},
		{
			name: "empty status",
			message: Message{
				Content: "This is a valid message content",
				Phone:   "+905551234567",
				Status:  "",
			},
			wantErr: true,
			errMsg:  "message status must be pending",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.ValidateForCreate()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMessage_ValidateForListByStatus(t *testing.T) {
	tests := []struct {
		name    string
		message Message
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid pending status",
			message: Message{
				Status: StatusPending,
			},
			wantErr: false,
		},
		{
			name: "valid sent status",
			message: Message{
				Status: StatusSent,
			},
			wantErr: false,
		},
		{
			name: "empty status",
			message: Message{
				Status: "",
			},
			wantErr: true,
			errMsg:  "message status must be provided",
		},
		{
			name: "invalid status",
			message: Message{
				Status: "INVALID",
			},
			wantErr: true,
			errMsg:  "message status must be one of PENDING or SENT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.ValidateForListByStatus()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMessage_ValidateForSent(t *testing.T) {
	validUUID := uuid.New().String()
	invalidUUID := "invalid-uuid"

	tests := []struct {
		name    string
		message Message
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid message",
			message: Message{
				ID: validUUID,
			},
			wantErr: false,
		},
		{
			name: "empty id",
			message: Message{
				ID: "",
			},
			wantErr: true,
			errMsg:  "message id must be provided",
		},
		{
			name: "invalid uuid",
			message: Message{
				ID: invalidUUID,
			},
			wantErr: true,
			errMsg:  "message id must be a valid uuid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.ValidateForSent()
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMessage_ErrorMethods(t *testing.T) {
	message := &Message{}

	tests := []struct {
		name     string
		method   func() error
		expected error
	}{
		{
			name:     "NewErrMessageDoesNotValidForCreate",
			method:   message.NewErrMessageDoesNotValidForCreate,
			expected: ErrMessageDoesNotValidForCreate,
		},
		{
			name:     "NewErrMessageDoesNotValidForListByStatus",
			method:   message.NewErrMessageDoesNotValidForListByStatus,
			expected: ErrMessageDoesNotValidForListByStatus,
		},
		{
			name:     "NewErrMessageDoesNotValidForSent",
			method:   message.NewErrMessageDoesNotValidForSent,
			expected: ErrMessageDoesNotValidForSent,
		},
		{
			name:     "NewErrMessageNotFound",
			method:   message.NewErrMessageNotFound,
			expected: ErrMessageNotFound,
		},
		{
			name:     "NewErrMessageStatusDoesNotEligibleForSent",
			method:   message.NewErrMessageStatusDoesNotEligibleForSent,
			expected: ErrMessageStatusDoesNotEligibleForSent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.method()
			assert.Equal(t, tt.expected, err)
		})
	}
}
