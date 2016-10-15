package engine
import(
    "time"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/spider/pool"
    "github.com/zl-leaf/mspider/logger"
)

const (
    stopSpiderWait = 1
)

type SpiderService struct {
    SpiderPool *pool.Pool
    EventListener chan spider.Param
    EventPublisher chan string
    State IState
    Validate func(request spider.Param) error
}

func (this *SpiderService) Start() error {
    nextState, err := this.State.Next(WorkingStateCode)
    if err != nil {
        return err
    }
    this.State = nextState
    go this.listen(this.EventListener)

    for _,s := range this.SpiderPool.All() {
        for _, u := range s.StartURLs() {
            this.EventPublisher <- u
        }
    }
    return nil
}

func (this *SpiderService) Stop() error {
    nextState, err := this.State.Next(StopStateCode)
    if err != nil {
        return err
    }
    this.State = nextState
    for {
        allFree := true
        for _,free := range this.SpiderPool.States() {
            if !free {
                allFree = false
                break
            }
        }
        if allFree {
            break
        }
        time.Sleep(time.Duration(stopSpiderWait) * time.Second)
    }
    nextState, err = this.State.Next(FreeStateCode)
    if err != nil {
        return err
    }
    this.State = nextState
    return nil
}

func (this *SpiderService) listen(listenerChan chan spider.Param) {
    for {
        param := <- listenerChan
        if this.Validate != nil {
            if err := this.Validate(param); err != nil {
                logger.Error(logger.SYSTEM, err.Error())
                continue
            }
        }
        s, err := this.SpiderPool.Get(param)
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }
        go this.do(param, s)
    }
}

func (this *SpiderService) do(param spider.Param, s *spider.Spider) {
    if this.State.Code() != WorkingStateCode {
        return
    }

    err := s.Do(param)
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
        return
    }
    defer this.SpiderPool.Put(s)
    redirects := s.Redirects()
    logger.Info(logger.SYSTEM, "spider id: %s finish url: %s, got %d redirects", s.ID, param.URL, len(redirects))
    for _, u := range redirects {
        if this.State.Code() != WorkingStateCode {
            break
        }
        this.EventPublisher <- u
    }
    return
}

func CreateSpiderService() (spiderService *SpiderService) {
    spiderService = &SpiderService{}
    spiderService.SpiderPool = pool.New()
    spiderService.EventListener = make(chan spider.Param)
    spiderService.State = createState(FreeStateCode)
    return
}