package scheduler
import(
    "github.com/zl-leaf/mspider/scheduler/queue"
)

type Scheduler struct {
    Queue *queue.Queue
    Interval int
}

func New() (scheduler *Scheduler, err error) {
    scheduler = &Scheduler{}
    scheduler.Queue = queue.New()
    return
}

func (this *Scheduler) Add(url string) {
    this.Queue.Add(url)
}

func (this *Scheduler) Head() (value string, err error) {
    e, err := this.Queue.Head()
    if err == nil {
        value = e.Value.(string)
    }
    return
}

func (this *Scheduler) Empty() bool {
    return this.Queue.Empty()
}