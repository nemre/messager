package message

import (
	"messager/domain/message"
	"messager/infrastructure/client"
)

type service struct {
	repository message.Repository
	client     client.Client
}

func New(repository message.Repository, client client.Client) message.Service {
	return &service{
		repository: repository,
		client:     client,
	}
}
