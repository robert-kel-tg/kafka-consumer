package main

import (
	"context"
	"fmt"
	"gopkg.in/confluentinc/confluent-kafka-go.v1/kafka"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	stopCh := make(chan bool)

	c, err := initConsumer()
	if err != nil {
		log.Printf("config error %v", err)
	}

	e := c.SubscribeTopics([]string{"foo", "^aRegex.*[f]oo"}, nil)
	if e != nil {
		log.Printf("subscribe to topic error %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go consume(c, ctx)

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

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		select {
		case <-signalCh:
			cancel()
			return
		}
	}()

	log.Print("Running kafka consumer...")
	<-stopCh
}

func initConsumer()(*kafka.Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
		"go.events.channel.enable": true,
		"enable.partition.eof": true,
		"session.timeout.ms": 6000,
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Created Consumer %v\n", c)

	return c, nil
}

func consume(consumer *kafka.Consumer, ctx context.Context) {
	defer func() {
		fmt.Printf("Stopped consumer\n")

		err := consumer.Close()
		if err != nil {
			log.Printf("subscribe to topic error %v", err)
		}
	}()

	for {
		select {
		case s := <-ctx.Done():
			log.Printf("Stopping kafka consumer...%v", s)
			return
		case event := <-consumer.Events():
			switch e := event.(type) {
			case *kafka.Message:
				fmt.Printf("%% Message on %s:\n%s\n", e.TopicPartition, string(e.Value))
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

