package message

import "fmt"

func (c *consumer) Stop() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("consumer.reader.Close: %w", err)
	}

	c.wg.Wait()

	return nil
}
