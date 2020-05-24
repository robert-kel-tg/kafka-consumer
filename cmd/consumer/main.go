package main

import (
	"context"
	"fmt"
	"github.com/robertke/kafka-consumer/pkg/infrastructure/config"
	"github.com/robertke/kafka-consumer/pkg/infrastructure/db"
	logger "github.com/robertke/kafka-consumer/pkg/infrastructure/log"
	"go.uber.org/zap"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	stopCh := make(chan bool)

	l, err := logger.NewLogger(logger.DebugLogConfig)
	if err != nil {
		log.Fatalf("could not init logger: %v", err)
	}

	sqlDB, err := db.Connect(conf.DB.Driver, conf.DB)
	if err != nil {
		log.Fatalf("could not init db connection: %v", err)
	}

	if err := db.Migrate(conf.DB.MigrationsPath, sqlDB); err != nil {
		log.Fatalf("could run migrations: %v", err)
	}

	sugar := l.Sugar()

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
		"bootstrap.servers": "broker:29092",
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

