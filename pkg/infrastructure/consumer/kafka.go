package consumer

import (
	"context"
	"fmt"
	"github.com/robertke/kafka-consumer/pkg/infrastructure/config"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
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

func (c *Consumer) consumeMessages(ctx context.Context, exec ExecFn) {
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
				sugar.Infof("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				c.consumer.Unassign()
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			}
		}
	}
}

