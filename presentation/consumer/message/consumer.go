package message

import (
	"sync"

	"messager/domain/message"

	"github.com/segmentio/kafka-go"
)

type Consumer interface {
	Start()
	Stop() error
}

type consumer struct {
	service message.Service
	reader  *kafka.Reader
	onError func(err error)
	wg      *sync.WaitGroup
}

func New(service message.Service, brokers []string, groupID, topic string, onError func(err error)) (Consumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		StartOffset: kafka.LastOffset,
	})

	if onError == nil {
		onError = func(err error) {}
	}

	return &consumer{
		service: service,
		reader:  reader,
		onError: onError,
		wg:      new(sync.WaitGroup),
	}, nil
}
