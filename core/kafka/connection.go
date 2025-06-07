package kafka

import (
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func (c *Client) Connect(createTopic bool) error {
	conn, err := kafka.Dial("tcp", c.addr)
	if err != nil {
		return fmt.Errorf("%w: Connect: %w", ErrKafkaClient, err)
	}

	c.kafkaConn = conn

	if createTopic {
		if c.numPartitions < 1 {
			return fmt.Errorf(
				"%w: Connect: %w: invalid number of partitions %d",
				ErrKafkaClient, ErrFailToCreateTopic, c.numPartitions,
			)
		}

		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             c.topic,
				NumPartitions:     c.numPartitions,
				ReplicationFactor: 1,
			},
		}

		err = c.kafkaConn.CreateTopics(topicConfigs...)
		if err != nil {
			return fmt.Errorf("%w: Connect: %w: %w", ErrKafkaClient, ErrFailToCreateTopic, err)
		}
	}

	if c.groupID != "" {
		c.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{c.addr},
			GroupID: c.groupID,
			Topic:   c.topic,
		})
	}

	c.writer = &kafka.Writer{
		Addr:            kafka.TCP(c.addr),
		Topic:           c.topic,
		WriteBackoffMin: 10 * time.Millisecond,
		WriteBackoffMax: 50 * time.Millisecond,
		BatchTimeout:    100 * time.Millisecond,
	}

	return nil
}

func (c *Client) Close() error {
	var errs []error

	if c.kafkaConn != nil {
		err := c.kafkaConn.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if c.writer != nil {
		err := c.writer.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if c.reader != nil {
		err := c.reader.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	err := fmt.Errorf("%w: Close", ErrKafkaClient)

	for _, e := range errs {
		err = fmt.Errorf("%w: %w", err, e)
	}

	return err
}
