package usecase

import (
	"reflect"
	"testing"

	"github.com/yoheinbb/healthd/internal/domain"
	"github.com/yoheinbb/healthd/internal/domain/repository"
	"github.com/yoheinbb/healthd/internal/infrastructure"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

func TestNewStatus(t *testing.T) {
	type args struct {
		domain     *domain.Status
		repository repository.Status
	}
	tests := []struct {
		name string
		args args
		want *Status
	}{
		{
			name: "new status",
			args: args{
				domain:     &domain.Status{},
				repository: &infrastructure.MockRepository{},
			},
			want: &Status{
				domain:     &domain.Status{},
				repository: &infrastructure.MockRepository{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStatus(tt.args.domain, tt.args.repository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStatus_GetStatus(t *testing.T) {
	type fields struct {
		domain     *domain.Status
		repository repository.Status
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "success status",
			fields: fields{
				domain:     domain.NewStatus(),
				repository: infrastructure.NewMockRepository(0),
			},
			want: constant.SUCCESS,
		},
		{
			name: "failed status",
			fields: fields{
				domain:     domain.NewStatus(),
				repository: infrastructure.NewMockRepository(1),
			},
			want: constant.FAILED,
		},
		{
			name: "MAINTENANCE status",
			fields: fields{
				domain:     domain.NewStatus(),
				repository: infrastructure.NewMockRepository(-1),
			},
			want: constant.MAINTENANCE,
		},
		{
			name: "default status",
			fields: fields{
				domain:     domain.NewStatus(),
				repository: infrastructure.NewMockRepository(2),
			},
			want: constant.FAILED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Status{
				domain:     tt.fields.domain,
				repository: tt.fields.repository,
			}
			// get status from repository to domain
			s.updateStatus()
			if got := s.GetStatus(); got != tt.want {
				t.Errorf("Status.GetStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
