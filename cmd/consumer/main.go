package main

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopCh := make(chan bool)

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	c, err := initConsumer(sugar)
	if err != nil {
		sugar.Errorf("config error %v", err)
	}

	e := c.SubscribeTopics([]string{"transactions"}, nil)
	if e != nil {
		sugar.Errorf("subscribe to topic error %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go consume(c, ctx, sugar)

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("Received interrupt signal, exiting...")
				stopCh <- true
				return
			default:
			}
		}
	}(ctx)

	interruptListener(cancel)

	sugar.Info("Running kafka consumer...")
	<-stopCh
}

func initConsumer(sugar *zap.SugaredLogger)(*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
		"go.events.channel.enable": true,
		"enable.partition.eof": false,
		"session.timeout.ms": 6000,
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	sugar.Infof("Created Consumer %v\n", c)

	return c, nil
}

func consume(consumer *kafka.Consumer, ctx context.Context, sugar *zap.SugaredLogger) {
	defer func() {
		sugar.Infof("Stopped consumer\n")

		err := consumer.Close()
		if err != nil {
			sugar.Errorf("subscribe to topic error %v", err)
		}
	}()

	for {
		select {
		case s := <-ctx.Done():
			sugar.Infof("Stopping kafka consumer...%v", s)
			return
		case event := <-consumer.Events():
			switch e := event.(type) {
			case *kafka.Message:
				sugar.Infof("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			case kafka.AssignedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Assign(e.Partitions)
			case kafka.RevokedPartitions:
				fmt.Fprintf(os.Stderr, "%% %v\n", e)
				consumer.Unassign()
			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
			}
		}
	}
}

func interruptListener(cancel context.CancelFunc) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()
}

