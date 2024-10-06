package usecase

import (
	"context"
	"log"
	"time"

	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/domain/repository"
)

type Status struct {
	status     *domain.Status
	repository repository.Status
}

func NewStatus(status *domain.Status, repository repository.Status) *Status {
	return &Status{status: status, repository: repository}
}

func (s *Status) GetStatus() string {
	return s.status.GetStatus()
}

func (ecs *Status) StartStatusUpdater(ctx context.Context, interval int) error {
	intervalTime := time.Duration(interval)

	ticker := time.NewTicker(time.Duration(intervalTime) * time.Second)
	if err := ecs.repository.UpdateStatus(); err != nil {
		log.Printf("%v", err)
	}
	ecs.updateStatus()
	// intervalTime毎にStatusをupdateする
	for {
		select {
		case <-ticker.C:
			if err := ecs.repository.UpdateStatus(); err != nil {
				log.Printf("%v", err)
			}
			ecs.updateStatus()
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Status) updateStatus() {
	s.status.UpdateStatus(s.repository.GetStatus())
}
