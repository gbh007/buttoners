package rabbitmq

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func (c *Client[T]) Connect(_ context.Context) (err error) {
	c.conn, err = amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s/",
		c.user, c.pass, c.addr,
	))
	if err != nil {
		return fmt.Errorf("%w: Connect: %w", ErrRabbitMQClient, err)
	}

	c.ch, err = c.conn.Channel()
	if err != nil {
		return fmt.Errorf("%w: Connect: %w", ErrRabbitMQClient, err)
	}

	c.queue, err = c.ch.QueueDeclare(
		c.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("%w: Connect: %w", ErrRabbitMQClient, err)
	}

	return nil
}

func (c *Client[T]) Close() error {
	var errs []error

	if c.ch != nil {
		err := c.ch.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if c.conn != nil {
		err := c.conn.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	err := fmt.Errorf("%w: Close", ErrRabbitMQClient)

	for _, e := range errs {
		err = fmt.Errorf("%w: %w", err, e)
	}

	return err
}
