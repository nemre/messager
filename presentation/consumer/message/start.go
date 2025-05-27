package message

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"messager/domain/message"
)

func (c *consumer) Start() {
	for {
		event, err := c.reader.ReadMessage(context.Background())
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			c.onError(fmt.Errorf("consumer.reader.ReadMessage: %w", err))

			continue
		}

		c.wg.Add(1)

		var key struct {
			ID string `json:"id"`
		}

		if err = json.Unmarshal(event.Key, &key); err != nil {
			c.onError(fmt.Errorf("json.Unmarshal: %w", err))
			c.wg.Done()

			continue
		}

		if key.ID == "" {
			c.onError(errors.New("message id is empty"))
			c.wg.Done()

			continue
		}

		var value struct {
			Before struct {
				Status string `json:"status"`
			} `json:"before"`
			After struct {
				Status string `json:"status"`
			} `json:"after"`
		}

		if err = json.Unmarshal(event.Value, &value); err != nil {
			c.onError(fmt.Errorf("json.Unmarshal: %w", err))
			c.wg.Done()

			continue
		}

		if value.Before.Status != string(message.StatusPending) {
			c.onError(fmt.Errorf("message %s before status is not PENDING", key.ID))
			c.wg.Done()

			continue
		}

		if value.After.Status != string(message.StatusSent) {
			c.onError(fmt.Errorf("message %s after status is not SENT", key.ID))
			c.wg.Done()

			continue
		}

		if err = c.service.Sent(context.Background(), message.Message{
			ID: key.ID,
		}); err != nil {
			c.onError(fmt.Errorf("consumer.service.Sent: %w", err))
		}

		c.wg.Done()
	}
}
