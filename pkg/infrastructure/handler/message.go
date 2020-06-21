package handler

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/robertke/kafka-consumer/pkg/domain"
	"github.com/robertke/kafka-consumer/pkg/infrastructure/consumer"
	"go.uber.org/zap"
)

type (
	FooHandler struct {
		log logger
		rw  fooReadWriter
	}

	logger interface {
		Sugar() *zap.SugaredLogger
	}

	fooReadWriter interface {
		fooWriter
		fooReader
	}

	fooWriter interface {
		Insert(ctx context.Context, tx *sqlx.Tx, msg domain.Foo) (*domain.Foo, error)
	}

	fooReader interface {
		GetOne(ctx context.Context, id string) (*domain.Foo, error)
	}
)

func NewMsg(log logger, rw fooReadWriter) *FooHandler {
	return &FooHandler{log, rw}
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
