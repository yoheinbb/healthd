package infrastructure

type MockRepository struct {
	status int
}

func NewMockRepository(status int) *MockRepository {
	return &MockRepository{status: status}
}

func (mr *MockRepository) GetStatus() int {
	return mr.status
}

func (mr *MockRepository) UpdateStatus() error {
	return nil
}
