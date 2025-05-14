package service

type Job interface {
	Name() string
	Run() error
}
