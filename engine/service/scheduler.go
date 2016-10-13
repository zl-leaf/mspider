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
    State IState
    Validator msg.ISchedulerValidator
}

func (this *SchedulerService) Start() error {
    nextState, err := this.State.Next(WorkingStateCode)
    if err != nil {
        return err
    }
    this.State = nextState
    go this.listen(this.EventListener)
    go this.push()
    return nil
}

func (this *SchedulerService) Stop() error {
    nextState, err := this.State.Next(FreeStateCode)
    if err != nil {
        return err
    }
    this.State = nextState
    return nil
}

func (this *SchedulerService) SetScheduler(s *scheduler.Scheduler) {
    this.Scheduler = s
}

func (this *SchedulerService) listen(listenerChan chan string) {
    for {
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
        if this.State.Code() != WorkingStateCode {
            break
        }
        this.EventPublisher <- u
        time.Sleep(time.Duration(this.Scheduler.Interval))
    }
}

func (this *SchedulerService) do(request string) {
    if this.State.Code() != WorkingStateCode {
        return
    }
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
    schedulerService.State = createState(FreeStateCode)
    return
}