package usecase

import (
	"context"
	"log"
	"time"

	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/domain/repository"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

type IStatus interface {
	GetStatus() string
}

type Status struct {
	domain     *domain.Status
	repository repository.Status
}

func NewStatus(domain_status *domain.Status, repository repository.Status) *Status {
	return &Status{domain: domain_status, repository: repository}
}

func (s *Status) GetStatus() string {
	return s.domain.GetStatus()
}

func (s *Status) StartStatusUpdater(ctx context.Context, interval int) error {
	intervalTime := time.Duration(interval)

	ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
	if err := s.repository.UpdateStatus(); err != nil {
		log.Printf("%v", err)
	}
	s.updateStatus()
	// intervalTime毎にStatusをupdateする
	for {
		select {
		case <-ticker.C:
			if err := s.repository.UpdateStatus(); err != nil {
				log.Printf("%v", err)
			}
			s.updateStatus()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Status) updateStatus() {
	var status string
	statusCode := s.repository.GetStatus()
	switch statusCode {
	case -1:
		status = constant.MAINTENANCE
	case 0:
		status = constant.SUCCESS
	default:
		status = constant.FAILED
	}

	s.domain.UpdateStatus(status)
}
