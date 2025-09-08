package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ShristiRnr/Finance_mierp/internal/adapters/database/db"
	"github.com/ShristiRnr/Finance_mierp/internal/core/ports"
	"github.com/google/uuid"
)

var (
	ErrInvalidName = errors.New("invalid allocation rule name")
)

type AccrualService struct {
	repo      ports.AccrualRepository
	publisher ports.EventPublisher
	topic     string
}

func NewAccrualService(r ports.AccrualRepository, producer ports.EventPublisher, topic string) *AccrualService {
	return &AccrualService{repo: r, publisher: producer, topic: topic}
}

func (s *AccrualService) Create(ctx context.Context, a db.Accrual) (db.Accrual, error) {
	acc, err := s.repo.Create(ctx, a)
	if err == nil {
		s.publishEvent(ctx, "AccrualCreated", acc)
	}
	return acc, err
}

func (s *AccrualService) Get(ctx context.Context, id uuid.UUID) (db.Accrual, error) {
	return s.repo.Get(ctx, id)
}

func (s *AccrualService) Update(ctx context.Context, a db.Accrual) (db.Accrual, error) {
	acc, err := s.repo.Update(ctx, a)
	if err == nil {
		s.publishEvent(ctx, "AccrualUpdated", acc)
	}
	return acc, err
}

func (s *AccrualService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.publishEvent(ctx, "AccrualDeleted", map[string]any{"id": id})
	}
	return err
}

func (s *AccrualService) List(ctx context.Context, limit, offset int32) ([]db.Accrual, error) {
	return s.repo.List(ctx, limit, offset)
}

// helper: publish Kafka event
func (s *AccrualService) publishEvent(ctx context.Context, eventType string, payload interface{}) {
	data := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}
	b, _ := json.Marshal(data)
	_ = s.publisher.Publish(ctx, s.topic, eventType, b) // ignore error for now
}

type allocationServiceImpl struct {
	repo      ports.AllocationRuleRepository
	publisher ports.EventPublisher
	topic     string
}

func NewAllocationService(
	r ports.AllocationRuleRepository,
	pub ports.EventPublisher,
	topic string,
) ports.AllocationService {
	return &allocationServiceImpl{
		repo:      r,
		publisher: pub,
		topic:     topic,
	}
}

func (s *allocationServiceImpl) CreateRule(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error) {
	if r.Name == "" {
		return db.AllocationRule{}, ErrInvalidName
	}
	rule, err := s.repo.Create(ctx, r)
	if err == nil {
		s.publishEvent(ctx, "AllocationRuleCreated", rule)
	}
	return rule, err
}

func (s *allocationServiceImpl) UpdateRule(ctx context.Context, r db.AllocationRule) (db.AllocationRule, error) {
	rule, err := s.repo.Update(ctx, r)
	if err == nil {
		s.publishEvent(ctx, "AllocationRuleUpdated", rule)
	}
	return rule, err
}

func (s *allocationServiceImpl) DeleteRule(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err == nil {
		s.publishEvent(ctx, "AllocationRuleDeleted", map[string]any{"id": id})
	}
	return err
}

func (s *allocationServiceImpl) GetRule(ctx context.Context, id uuid.UUID) (db.AllocationRule, error) {
	return s.repo.Get(ctx, id)
}

func (s *allocationServiceImpl) ListRules(ctx context.Context, limit, offset int32) ([]db.AllocationRule, error) {
	return s.repo.List(ctx, limit, offset)
}

func (s *allocationServiceImpl) publishEvent(ctx context.Context, eventType string, payload interface{}) {
	data := map[string]interface{}{
		"type":    eventType,
		"payload": payload,
	}
	b, _ := json.Marshal(data)
	_ = s.publisher.Publish(ctx, s.topic, eventType, b) // ignore error for now
}
