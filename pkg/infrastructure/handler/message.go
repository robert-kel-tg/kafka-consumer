package handler

import (
	"context"
	"github.com/robertke/kafka-consumer/pkg/infrastructure/consumer"
	"go.uber.org/zap"
)

type (
	FooHandler struct {
		log logger
	}

	logger interface {
		Sugar() *zap.SugaredLogger
	}
)

func NewMsg(log logger) *FooHandler {
	return &FooHandler{log}
}

func (f FooHandler) Handle(ctx context.Context, msg consumer.ConsumedMessage) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}
	sugar := f.log.Sugar()

	sugar.Infof("FooHandler: Topic %s Msg %s", msg.Topic, msg.Body)

	return nil
}
