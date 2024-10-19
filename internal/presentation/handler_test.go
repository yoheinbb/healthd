package presentation

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/yoheinbb/healthd/internal/usecase"
	"github.com/yoheinbb/healthd/internal/util/constant"
)

func TestHandler_HealthdHandler(t *testing.T) {
	type fields struct {
		status     usecase.IStatus
		retSuccess string
		retFailed  string
	}
	tests := []struct {
		name   string
		fields fields
		want   OutputSchema
	}{
		{
			name: "success",
			fields: fields{
				status:     usecase.NewMockStatus(constant.SUCCESS),
				retSuccess: "success",
				retFailed:  "failed",
			},
			want: OutputSchema{Result: "success"},
		},
		{
			name: "failed",
			fields: fields{
				status:     usecase.NewMockStatus(constant.FAILED),
				retSuccess: "success",
				retFailed:  "failed",
			},
			want: OutputSchema{Result: "failed"},
		},
		{
			name: "others",
			fields: fields{
				status:     usecase.NewMockStatus(constant.MAINTENANCE),
				retSuccess: "success",
				retFailed:  "failed",
			},
			want: OutputSchema{Result: "failed"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			h := &Handler{
				status:     tt.fields.status,
				retSuccess: tt.fields.retSuccess,
				retFailed:  tt.fields.retFailed,
			}
			h.HealthdHandler(c)
			assert.Equal(t, 200, w.Code)

			var got OutputSchema
			err := json.Unmarshal(w.Body.Bytes(), &got)
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
