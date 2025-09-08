package kafka

import (
	"context"
	"encoding/json"
	"log"
	"errors"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writers map[string]*kafka.Writer
}

func NewKafkaPublisher(brokers []string) ports.EventPublisher {
	return &KafkaPublisher{
		writers: map[string]*kafka.Writer{
			"accounts": {
				Addr:     kafka.TCP(brokers...),
				Topic:    "accounts",
				Balancer: &kafka.LeastBytes{},
			},
			"journals": {
				Addr:     kafka.TCP(brokers...),
				Topic:    "journals",
				Balancer: &kafka.LeastBytes{},
			},
		},
	}
}

func (p *KafkaPublisher) PublishAccountCreated(ctx context.Context, a db.Account) error {
	return p.publish(ctx, "accounts", "account.created", a)
}

func (p *KafkaPublisher) PublishJournalCreated(ctx context.Context, j db.JournalEntry) error {
	return p.publish(ctx, "journals", "journal.created", j)
}

func (p *KafkaPublisher) PublishAccountUpdated(ctx context.Context, a db.Account) error {
	return p.publish(ctx, "accounts", "account.updated", a)
}

func (p *KafkaPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, "accounts", "account.deleted", id)
}

func (p *KafkaPublisher) PublishJournalUpdated(ctx context.Context, j db.JournalEntry) error {
	return p.publish(ctx, "journals", "journal.updated", j)
}

func (p *KafkaPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, "journals", "journal.deleted", id)
}

func (p *KafkaPublisher) publish(ctx context.Context, topic, key string, v interface{}) error {
	writer, ok := p.writers[topic]
	if !ok {
		return errors.New("kafka writer not found for topic: " + topic)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	msg := kafka.Message{
		Key:   []byte(key),
		Value: data,
	}
	if err := writer.WriteMessages(ctx, msg); err != nil {
		log.Printf("kafka publish error topic=%s: %v", topic, err)
		return err
	}
	return nil
}
