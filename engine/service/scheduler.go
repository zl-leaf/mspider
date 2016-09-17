package service
import(
    "time"
    "github.com/zl-leaf/mspider/scheduler"
)

type SchedulerService struct {
    mScheduler *scheduler.Scheduler
    engineListener chan string
    eventPublisher chan string
    listener IService
    state int
}

func (this *SchedulerService) EventPublisher() chan string {
    return this.eventPublisher
}

func (this *SchedulerService) Start() error {
    this.state = WorkingState
    go this.listen(this.listener.EventPublisher())
    go this.push()
    return nil
}

func (this *SchedulerService) Stop() error {
    this.state = StopState
    return nil
}

func (this *SchedulerService) AddListener(s IService) error {
    this.listener = s
    return nil
}

func (this *SchedulerService) SetScheduler(s *scheduler.Scheduler) {
    this.mScheduler = s
}

func (this *SchedulerService) listen(listenerChan chan string) {
    for {
        if this.state == StopState {
            break
        }
        value := <- listenerChan
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
        time.Sleep(time.Duration(this.mScheduler.Interval))
    }
}

func (this *SchedulerService) do(content string) {
    this.mScheduler.Add(content)
}

func CreateSchedulerService(engineListener chan string) (schedulerService *SchedulerService) {
    schedulerService = &SchedulerService{}
    schedulerService.engineListener = engineListener
    schedulerService.eventPublisher = make(chan string)
    return
}