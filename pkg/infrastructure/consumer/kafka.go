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

	go c.consumeMessages(ctx, exec)

	return nil
}

func (c *Consumer) consumeMessages(ctx context.Context, execFn ExecFn) {
	sugar := c.log.Sugar()
	defer func() {
		sugar.Infof("Stopped consumer\n")

		err := c.consumer.Close()
		if err != nil {
			sugar.Errorf("subscribe to topic error %v", err)
		}
	}()

	for {
		select {
		case s := <-ctx.Done():
			sugar.Infof("Stopping kafka consumer...%v", s)
			return
		case event := <-c.consumer.Events():
			switch e := event.(type) {
			case *kafka.Message:
				if err := execFn(ctx, ConsumedMessage{
					Topic: *e.TopicPartition.Topic,
					Body:  e.Value,
				}); err != nil {
					c.log.Error("an error occurred", zap.Error(err))
				}
				sugar.Infof("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
			case kafka.Error:
				c.log.Error("error consuming", zap.Error(e))
			case kafka.PartitionEOF:
				c.log.Warn("partition EOF", zap.Any("partition", event))
			case kafka.OffsetsCommitted:
				if e.Error != nil {
					c.log.Error("an error occurred offsets commit", zap.Error(e.Error), zap.Any("offsets", e.Offsets))
				}
			}
		}
	}
}

