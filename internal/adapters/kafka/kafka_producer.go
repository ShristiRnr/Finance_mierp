package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"
	"fmt"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/segmentio/kafka-go"
)

type KafkaPublisher struct {
	writers map[string]*kafka.Writer
}

// NewKafkaPublisher initializes Kafka writers for accounts and journals topics
func NewKafkaPublisher(brokers []string) ports.EventPublisher {
	return &KafkaPublisher{
		writers: map[string]*kafka.Writer{
			"accounts": {
				Addr:         kafka.TCP(brokers...),
				Topic:        "accounts",
				Balancer:     &kafka.LeastBytes{},
				RequiredAcks: kafka.RequireAll,
				Async:        false,
			},
			"journals": {
				Addr:         kafka.TCP(brokers...),
				Topic:        "journals",
				Balancer:     &kafka.LeastBytes{},
				RequiredAcks: kafka.RequireAll,
				Async:        false,
			},
			"accruals": {
				Addr:         kafka.TCP(brokers...),
				Topic:        "accruals",
				Balancer:     &kafka.LeastBytes{},
				RequiredAcks: kafka.RequireAll,
				Async:        false,
			},
			"allocation_rules": {
				Addr:         kafka.TCP(brokers...),
				Topic:        "allocation_rules",
				Balancer:     &kafka.LeastBytes{},
				RequiredAcks: kafka.RequireAll,
				Async:        false,
			},
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	writer, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("writer not found for topic %s", topic)
	}
	msg := kafka.Message{Key: []byte(key), Value: payload}
	return writer.WriteMessages(ctx, msg)
}


func (p *KafkaPublisher) PublishAccountCreated(ctx context.Context, a db.Account) error {
	return p.publish(ctx, "accounts", "account.created", a)
}

func (p *KafkaPublisher) PublishAccountUpdated(ctx context.Context, a db.Account) error {
	return p.publish(ctx, "accounts", "account.updated", a)
}

func (p *KafkaPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, "accounts", "account.deleted", id)
}

func (p *KafkaPublisher) PublishJournalCreated(ctx context.Context, j db.JournalEntry) error {
	return p.publish(ctx, "journals", "journal.created", j)
}

func (p *KafkaPublisher) PublishJournalUpdated(ctx context.Context, j db.JournalEntry) error {
	return p.publish(ctx, "journals", "journal.updated", j)
}

func (p *KafkaPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, "journals", "journal.deleted", id)
}

// --- Accrual events ---
func (p *KafkaPublisher) PublishAccrualCreated(ctx context.Context, a db.Accrual) error {
	return p.publish(ctx, "accruals", "accrual.created", a)
}

func (p *KafkaPublisher) PublishAccrualUpdated(ctx context.Context, a db.Accrual) error {
	return p.publish(ctx, "accruals", "accrual.updated", a)
}

func (p *KafkaPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, "accruals", "accrual.deleted", id)
}

// --- AllocationRule events ---
func (p *KafkaPublisher) PublishAllocationRuleCreated(ctx context.Context, rule db.AllocationRule) error {
	return p.publish(ctx, "allocation_rules", "allocation_rule.created", rule)
}
func (p *KafkaPublisher) PublishAllocationRuleUpdated(ctx context.Context, rule db.AllocationRule) error {
	return p.publish(ctx, "allocation_rules", "allocation_rule.updated", rule)
}

func (p *KafkaPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, "allocation_rules", "allocation_rule.deleted", id)
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
		Time:  time.Now(),
	}

	// simple retry mechanism
	for i := 0; i < 3; i++ {
		if err := writer.WriteMessages(ctx, msg); err != nil {
			log.Printf("kafka publish error attempt=%d topic=%s: %v", i+1, topic, err)
			time.Sleep(time.Second * time.Duration(i+1))
			continue
		}
		return nil
	}
	return err
}

// Close gracefully closes all writers
func (p *KafkaPublisher) Close() error {
	for _, w := range p.writers {
		if err := w.Close(); err != nil {
			return err
		}
	}
	return nil
}

//
// ---------------- Consumer ----------------
//

type KafkaConsumer struct {
	readers map[string]*kafka.Reader
}

// NewKafkaConsumer initializes Kafka readers for accounts and journals topics
func NewKafkaConsumer(brokers []string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		readers: map[string]*kafka.Reader{
			"accounts": kafka.NewReader(kafka.ReaderConfig{
				Brokers:  brokers,
				Topic:    "accounts",
				GroupID:  groupID,
				MinBytes: 10e3, // 10KB
				MaxBytes: 10e6, // 10MB
			}),
			"journals": kafka.NewReader(kafka.ReaderConfig{
				Brokers:  brokers,
				Topic:    "journals",
				GroupID:  groupID,
				MinBytes: 10e3,
				MaxBytes: 10e6,
			}),
		},
	}
}

// Consume starts consuming messages from the given topic
func (c *KafkaConsumer) Consume(ctx context.Context, topic string, handler func(key, value []byte) error) error {
	reader, ok := c.readers[topic]
	if !ok {
		return errors.New("kafka reader not found for topic: " + topic)
	}

	for {
		select {
		case <-ctx.Done():
			log.Printf("stopping consumer for topic=%s", topic)
			return ctx.Err()

		default:
			m, err := reader.ReadMessage(ctx)
			if err != nil {
				// exit if context cancelled
				if errors.Is(err, context.Canceled) {
					return nil
				}
				log.Printf("kafka consume error topic=%s: %v", topic, err)
				continue
			}

			log.Printf("Consumed message topic=%s key=%s value=%s", topic, string(m.Key), string(m.Value))
			if handler != nil {
				if err := handler(m.Key, m.Value); err != nil {
					log.Printf("handler error: %v", err)
				}
			}
		}
	}
}

// Close gracefully closes all readers
func (c *KafkaConsumer) Close() error {
	for _, r := range c.readers {
		if err := r.Close(); err != nil {
			return err
		}
	}
	return nil
}
