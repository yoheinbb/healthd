package repository

type Status interface {
	GetStatus() int
	UpdateStatus() error
}
