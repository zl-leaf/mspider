package scheduler
import(
    "github.com/zl-leaf/mspider/scheduler/queue"
)

type Scheduler struct {
    queue *queue.Queue
    Interval int
}

func New() (scheduler *Scheduler, err error) {
    scheduler = &Scheduler{}
    scheduler.queue = queue.New()
    return
}

func (this *Scheduler)Add(url string) {
    this.queue.Add(url)
}

func (this *Scheduler)Head() (value string, err error) {
    e, err := this.queue.Head()
    if err == nil {
        value = e.Value.(string)
    }
    return
}