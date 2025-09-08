package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/segmentio/kafka-go"
)

//
// ---------------- Constants ----------------
//
const (
	TopicAccounts        = "accounts"
	TopicJournals        = "journals"
	TopicAccruals        = "accruals"
	TopicAllocationRules = "allocation_rules"
	TopicAuditEvents     = "audit_events"

	EventAccountCreated = "account.created"
	EventAccountUpdated = "account.updated"
	EventAccountDeleted = "account.deleted"

	EventJournalCreated = "journal.created"
	EventJournalUpdated = "journal.updated"
	EventJournalDeleted = "journal.deleted"

	EventAccrualCreated = "accrual.created"
	EventAccrualUpdated = "accrual.updated"
	EventAccrualDeleted = "accrual.deleted"

	EventAllocationRuleCreated = "allocation_rule.created"
	EventAllocationRuleUpdated = "allocation_rule.updated"
	EventAllocationRuleDeleted = "allocation_rule.deleted"

	EventAuditRecorded = "audit.recorded"
)

//
// ---------------- Publisher ----------------
//

type KafkaPublisher struct {
	writers map[string]*kafka.Writer
}

func newWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
	}
}

func NewKafkaPublisher(brokers []string) ports.EventPublisher {
	topics := []string{
		TopicAccounts,
		TopicJournals,
		TopicAccruals,
		TopicAllocationRules,
		TopicAuditEvents,
	}

	writers := make(map[string]*kafka.Writer, len(topics))
	for _, t := range topics {
		writers[t] = newWriter(brokers, t)
	}

	return &KafkaPublisher{writers: writers}
}

func (p *KafkaPublisher) Publish(ctx context.Context, topic, key string, payload []byte) error {
	writer, ok := p.writers[topic]
	if !ok {
		return fmt.Errorf("writer not found for topic %s", topic)
	}
	msg := kafka.Message{Key: []byte(key), Value: payload}
	return writer.WriteMessages(ctx, msg)
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

	// retry with exponential backoff and context timeout
	for i := 0; i < 3; i++ {
		ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		err = writer.WriteMessages(ctxTimeout, msg)
		cancel()
		if err == nil {
			return nil
		}
		log.Printf("kafka publish error attempt=%d topic=%s: %v", i+1, topic, err)
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return err
}

// --- Account events ---
func (p *KafkaPublisher) PublishAccountCreated(ctx context.Context, a *db.Account) error {
	return p.publish(ctx, TopicAccounts, EventAccountCreated, a)
}
func (p *KafkaPublisher) PublishAccountUpdated(ctx context.Context, a *db.Account) error {
	return p.publish(ctx, TopicAccounts, EventAccountUpdated, a)
}
func (p *KafkaPublisher) PublishAccountDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicAccounts, EventAccountDeleted, id)
}

// --- Journal events ---
func (p *KafkaPublisher) PublishJournalCreated(ctx context.Context, j *db.JournalEntry) error {
	return p.publish(ctx, TopicJournals, EventJournalCreated, j)
}
func (p *KafkaPublisher) PublishJournalUpdated(ctx context.Context, j *db.JournalEntry) error {
	return p.publish(ctx, TopicJournals, EventJournalUpdated, j)
}
func (p *KafkaPublisher) PublishJournalDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicJournals, EventJournalDeleted, id)
}

// --- Accrual events ---
func (p *KafkaPublisher) PublishAccrualCreated(ctx context.Context, a *db.Accrual) error {
	return p.publish(ctx, TopicAccruals, EventAccrualCreated, a)
}
func (p *KafkaPublisher) PublishAccrualUpdated(ctx context.Context, a *db.Accrual) error {
	return p.publish(ctx, TopicAccruals, EventAccrualUpdated, a)
}
func (p *KafkaPublisher) PublishAccrualDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicAccruals, EventAccrualDeleted, id)
}

// --- AllocationRule events ---
func (p *KafkaPublisher) PublishAllocationRuleCreated(ctx context.Context, rule *db.AllocationRule) error {
	return p.publish(ctx, TopicAllocationRules, EventAllocationRuleCreated, rule)
}
func (p *KafkaPublisher) PublishAllocationRuleUpdated(ctx context.Context, rule *db.AllocationRule) error {
	return p.publish(ctx, TopicAllocationRules, EventAllocationRuleUpdated, rule)
}
func (p *KafkaPublisher) PublishAllocationRuleDeleted(ctx context.Context, id string) error {
	return p.publish(ctx, TopicAllocationRules, EventAllocationRuleDeleted, id)
}

// --- Audit events ---
func (p *KafkaPublisher) PublishAuditRecorded(ctx context.Context, e *db.AuditEvent) error {
	return p.publish(ctx, TopicAuditEvents, EventAuditRecorded, e)
}


// Close gracefully closes all writers
func (p *KafkaPublisher) Close() error {
	for topic, w := range p.writers {
		if err := w.Close(); err != nil {
			log.Printf("failed to close writer for topic=%s: %v", topic, err)
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

func newReader(brokers []string, groupID, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Second,
	})
}

func NewKafkaConsumer(brokers []string, groupID string) *KafkaConsumer {
	return &KafkaConsumer{
		readers: map[string]*kafka.Reader{
			TopicAccounts:        newReader(brokers, groupID, TopicAccounts),
			TopicJournals:        newReader(brokers, groupID, TopicJournals),
			TopicAccruals:        newReader(brokers, groupID, TopicAccruals),
			TopicAllocationRules: newReader(brokers, groupID, TopicAllocationRules),
			TopicAuditEvents:     newReader(brokers, groupID, TopicAuditEvents),
		},
	}
}

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
	for topic, r := range c.readers {
		if err := r.Close(); err != nil {
			log.Printf("failed to close reader for topic=%s: %v", topic, err)
			return err
		}
	}
	return nil
}
