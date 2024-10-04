package usecase

import (
	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

type Status struct {
	Status *domain.Status
}

func NewStatus(status *domain.Status) *Status {
	return &Status{Status: status}
}

func (s *Status) SetFailed() {
	s.Status.Status = constant.Failed
}

func (s *Status) SetSucess() {
	s.Status.Status = constant.Success
}

func (s *Status) GetStatus() string {
	return s.Status.Status
}
