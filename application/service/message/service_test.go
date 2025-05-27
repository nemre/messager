package message_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"messager/application/service/message"
	entity "messager/domain/message"
	"messager/infrastructure/database/postgresql"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Create(ctx context.Context, msg *entity.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *mockRepository) FindAllByStatus(ctx context.Context, status entity.Status) ([]entity.Message, error) {
	args := m.Called(ctx, status)
	return args.Get(0).([]entity.Message), args.Error(1)
}

func (m *mockRepository) FindByID(ctx context.Context, id string) (*entity.Message, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.Message), args.Error(1)
}

func (m *mockRepository) UpdateAllStatusesByStatus(ctx context.Context, from, to entity.Status) error {
	args := m.Called(ctx, from, to)
	return args.Error(0)
}

func (m *mockRepository) CreateSentInfo(ctx context.Context, messageID, t string) error {
	args := m.Called(ctx, messageID, t)
	return args.Error(0)
}

type mockClient struct {
	mock.Mock
}

func (m *mockClient) SendMessage(ctx context.Context, msg entity.Message) (string, error) {
	args := m.Called(ctx, msg)
	return args.String(0), args.Error(1)
}

func validMessage() entity.Message {
	return entity.Message{
		ID:      uuid.New().String(),
		Content: "This is a valid message content",
		Phone:   "+905551234567",
		Status:  entity.StatusPending,
	}
}

func TestService_Create(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	msg := validMessage()

	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("Create", ctx, mock.AnythingOfType("*message.Message")).Return(nil)
		svc := message.New(repo, cli)
		got, err := svc.Create(ctx, msg)
		assert.NoError(t, err)
		assert.Equal(t, msg.Content, got.Content)
		repo.AssertExpectations(t)
	})

	t.Run("invalid message", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		invalid := msg
		invalid.Content = "short"
		svc := message.New(repo, cli)
		got, err := svc.Create(ctx, invalid)
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("Create", ctx, mock.AnythingOfType("*message.Message")).Return(errors.New("db error"))
		svc := message.New(repo, cli)
		_, err := svc.Create(ctx, msg)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_ListByStatus(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	status := entity.StatusPending
	msg := validMessage()

	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindAllByStatus", ctx, status).Return([]entity.Message{msg}, nil)
		svc := message.New(repo, cli)
		got, err := svc.ListByStatus(ctx, status)
		assert.NoError(t, err)
		assert.Len(t, got, 1)
		repo.AssertExpectations(t)
	})

	t.Run("invalid status", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		svc := message.New(repo, cli)
		_, err := svc.ListByStatus(ctx, "INVALID")
		assert.Error(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindAllByStatus", ctx, status).Return([]entity.Message{}, errors.New("db error"))
		svc := message.New(repo, cli)
		_, err := svc.ListByStatus(ctx, status)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_Process(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("UpdateAllStatusesByStatus", ctx, entity.StatusPending, entity.StatusSent).Return(nil)
		svc := message.New(repo, cli)
		err := svc.Process(ctx)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("UpdateAllStatusesByStatus", ctx, entity.StatusPending, entity.StatusSent).Return(errors.New("db error"))
		svc := message.New(repo, cli)
		err := svc.Process(ctx)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}

func TestService_Sent(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	msg := validMessage()
	msg.Status = entity.StatusSent

	t.Run("success", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindByID", ctx, msg.ID).Return(&msg, nil)
		cli.On("SendMessage", ctx, msg).Return("sent-id", nil)
		repo.On("CreateSentInfo", ctx, "sent-id", mock.AnythingOfType("string")).Return(nil)
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, msg)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
		cli.AssertExpectations(t)
	})

	t.Run("invalid message", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		invalid := msg
		invalid.ID = ""
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, invalid)
		assert.Error(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindByID", ctx, msg.ID).Return((*entity.Message)(nil), errors.New("not found"))
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, msg)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("status not eligible", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		m := msg
		m.Status = entity.StatusPending
		repo.On("FindByID", ctx, m.ID).Return(&m, nil)
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, m)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("client error", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindByID", ctx, msg.ID).Return(&msg, nil)
		cli.On("SendMessage", ctx, msg).Return("", errors.New("client error"))
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, msg)
		assert.Error(t, err)
		repo.AssertExpectations(t)
		cli.AssertExpectations(t)
	})

	t.Run("repo CreateSentInfo error", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindByID", ctx, msg.ID).Return(&msg, nil)
		cli.On("SendMessage", ctx, msg).Return("sent-id", nil)
		repo.On("CreateSentInfo", ctx, "sent-id", mock.AnythingOfType("string")).Return(errors.New("db error"))
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, msg)
		assert.Error(t, err)
		repo.AssertExpectations(t)
		cli.AssertExpectations(t)
	})

	t.Run("not found with ErrNoRows", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindByID", ctx, msg.ID).Return((*entity.Message)(nil), postgresql.ErrNoRows)
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, msg)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("FindByID returns error with non-nil message", func(t *testing.T) {
		repo := new(mockRepository)
		cli := new(mockClient)
		repo.On("FindByID", ctx, msg.ID).Return(&msg, errors.New("unexpected error"))
		svc := message.New(repo, cli)
		err := svc.Sent(ctx, msg)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
