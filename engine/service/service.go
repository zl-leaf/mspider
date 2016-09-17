package service

type IService interface {
    EventPublisher() chan string
    Start() error
    Stop() error
    AddListener(s IService) error
}

const (
    WorkingState = 1
    StopState = 2
)