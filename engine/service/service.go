package service

type IService interface {
    EventPublisher() chan string
    Start() error
    AddListener(s IService) error
}