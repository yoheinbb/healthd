package repository

type Status interface {
	GetStatus() string
	UpdateStatus() error
}
