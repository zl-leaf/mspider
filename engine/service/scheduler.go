package service
import(
    "github.com/zl-leaf/mspider/scheduler"
    "time"
)

type SchedulerService struct {
    mScheduler *scheduler.Scheduler
    engineListener chan string
    eventPublisher chan string
    listeners []chan string
}

func (this *SchedulerService) EventPublisher() chan string {
    return this.eventPublisher
}

func (this *SchedulerService) Start() error {
    for _,listener := range this.listeners {
        go this.listen(listener)
    }
    go this.push()
    return nil
}

func (this *SchedulerService) AddListener(s IService) error {
    this.listeners = append(this.listeners, s.EventPublisher())
    return nil
}

func (this *SchedulerService) listen(listener chan string) {
    for {
        value := <- listener
        this.do(value)
    }
}

func (this *SchedulerService) push() {
    for {
        u,err := this.mScheduler.Head()

        if err != nil {
            continue
        }

        this.eventPublisher <- u
        time.Sleep(time.Duration(1) * time.Second)
    }
}

func (this *SchedulerService) do(content string) {
    this.mScheduler.Add(content)
}

func CreateSchedulerService(engineListener chan string, s *scheduler.Scheduler) (schedulerService *SchedulerService) {
    schedulerService = &SchedulerService{mScheduler:s}
    schedulerService.engineListener = engineListener
    schedulerService.eventPublisher = make(chan string)
    schedulerService.listeners = make([]chan string, 0)
    return
}