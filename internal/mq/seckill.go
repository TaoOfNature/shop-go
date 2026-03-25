package mq

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/dawnstack/shop-go/internal/config"
	"github.com/segmentio/kafka-go"
)

type SeckillMessage struct {
	MessageID      string `json:"message_id"`
	UserID         int64  `json:"user_id"`
	ProductID      int64  `json:"product_id"`
	Quantity       int    `json:"quantity"`
	AddressID      int64  `json:"address_id"`
	IdempotencyKey string `json:"idempotency_key"`
}

type SeckillPublisher interface {
	Publish(ctx context.Context, msg SeckillMessage) error
}

type InlinePublisher struct {
	handler func(context.Context, SeckillMessage) error
}

func NewInlinePublisher(handler func(context.Context, SeckillMessage) error) *InlinePublisher {
	return &InlinePublisher{handler: handler}
}

func (p *InlinePublisher) Publish(ctx context.Context, msg SeckillMessage) error {
	return p.handler(ctx, msg)
}

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(cfg config.KafkaConfig) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(cfg.Brokers...),
			Topic:        cfg.Topic,
			RequiredAcks: kafka.RequireOne,
			Balancer:     &kafka.Hash{},
			Async:        false,
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, msg SeckillMessage) error {
	bytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.MessageID),
		Value: bytes,
		Time:  time.Now(),
	})
}

func StartSeckillConsumer(ctx context.Context, cfg config.KafkaConfig, handler func(context.Context, SeckillMessage) error) {
	if !cfg.Enabled {
		return
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})

	go func() {
		defer func() { _ = reader.Close() }()
		for {
			message, err := reader.ReadMessage(ctx)
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				slog.Error("kafka read failed", "error", err)
				continue
			}

			var payload SeckillMessage
			if err := json.Unmarshal(message.Value, &payload); err != nil {
				slog.Error("kafka payload decode failed", "error", err)
				continue
			}

			if err := handler(ctx, payload); err != nil {
				slog.Error("seckill consumer handle failed", "error", err, "message_id", payload.MessageID)
			}
		}
	}()
}
