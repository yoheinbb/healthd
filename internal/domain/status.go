package domain

import "github.com/yoheinbb/healthd/internal/util/constant"

type Status struct {
	s string
}

func NewStatus() *Status {
	return &Status{s: constant.MAINTENANCE}
}
func (s *Status) GetStatus() string {
	return s.s
}
func (s *Status) UpdateStatus(status string) {
	s.s = status

}
