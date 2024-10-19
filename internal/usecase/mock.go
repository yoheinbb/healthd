package usecase

type MockStatus struct {
	status string
}

func NewMockStatus(status string) *MockStatus {
	return &MockStatus{status: status}
}

func (m *MockStatus) GetStatus() string {
	return m.status
}
