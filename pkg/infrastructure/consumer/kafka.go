package consumer

import (
	"context"
	"github.com/robertke/kafka-consumer/pkg/infrastructure/config"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
)

type (
	Consumer struct {
		consumer   *kafka.Consumer
		topics     []string
		log        *zap.Logger
	}

	ConsumedMessage struct {
		Topic string
		Body  []byte
	}

	ExecFn func(ctx context.Context, msg ConsumedMessage) error
)

func New(conf *config.Config, log *zap.Logger) (*Consumer, error) {

	cnf, err := conf.AddKafkaConf()
	if err != nil {
		return nil, err
	}

	cons, err := kafka.NewConsumer(&cnf)
	if err != nil {
		return nil, err
	}

	res := &Consumer{
		consumer: cons,
		topics: conf.ConsumerConfig.Topics,
		log: log,
	}

	return res, nil
}

func (c *Consumer) Run(ctx context.Context, exec ExecFn) error {
	err := c.consumer.SubscribeTopics(c.topics, nil)
	if err != nil {
		return err
	}

	go func() {
		c.consumeMessages(ctx, exec)
	}()

	return nil
}

func (c *Consumer) consumeMessages(ctx context.Context, exec ExecFn) {

}

