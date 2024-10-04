package domain

import "github.com/yoheinbb/healthd/internal/util/constant"

type Status struct {
	Status string
}

func NewStatus() *Status {
	return &Status{Status: constant.Maintenance}
}
