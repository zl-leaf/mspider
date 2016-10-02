package service
import(
    "time"
    "github.com/zl-leaf/mspider/engine/msg"
    "github.com/zl-leaf/mspider/scheduler"
    "github.com/zl-leaf/mspider/logger"
)

type SchedulerService struct {
    Scheduler *scheduler.Scheduler
    EventListener chan string
    EventPublisher chan string
    State int
    Validator msg.ISchedulerValidator
}

func (this *SchedulerService) Start() error {
    this.State = WorkingState
    go this.listen(this.EventListener)
    go this.push()
    return nil
}

func (this *SchedulerService) Stop() error {
    this.State = StopState
    return nil
}

func (this *SchedulerService) SetScheduler(s *scheduler.Scheduler) {
    this.Scheduler = s
}

func (this *SchedulerService) listen(listenerChan chan string) {
    for {
        if this.State == StopState {
            break
        }
        request := <- listenerChan
        this.do(request)
    }
}

func (this *SchedulerService) push() {
    for {
        if this.Scheduler.Empty() {
            continue
        }
        u, err := this.Scheduler.Head()
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }
        if this.State == StopState {
            break
        }
        this.EventPublisher <- u
        time.Sleep(time.Duration(this.Scheduler.Interval))
    }
}

func (this *SchedulerService) do(request string) {
    if this.Validator != nil {
        if err := this.Validator.Validate(request); err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            return
        }
    }
    this.Scheduler.Add(request)
}

func CreateSchedulerService() (schedulerService *SchedulerService) {
    schedulerService = &SchedulerService{}
    schedulerService.EventListener = make(chan string, 10)
    return
}