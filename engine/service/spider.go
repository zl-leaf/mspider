package service
import(
    "time"
    "github.com/zl-leaf/mspider/engine/msg"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/spider/pool"
    "github.com/zl-leaf/mspider/logger"
)

type SpiderService struct {
    SpiderPool *pool.Pool
    EventListener chan msg.SpiderRequest
    EventPublisher chan string
    State int
    Validator msg.ISpiderValidator
}

func (this *SpiderService) Start() error {
    this.State = WorkingState
    go this.listen(this.EventListener)

    for _,s := range this.SpiderPool.All() {
        for _, u := range s.StartURLs() {
            this.EventPublisher <- u
        }
    }
    return nil
}

func (this *SpiderService) Stop() error {
    this.State = StopState
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
    return nil
}

func (this *SpiderService) listen(listenerChan chan msg.SpiderRequest) {
    for {
        request := <- listenerChan
        if this.Validator != nil {
            if err := this.Validator.Validate(request); err != nil {
                logger.Error(logger.SYSTEM, err.Error())
                continue
            }
        }
        param := spider.Param{URL:request.URL, Data:request.Data, ContentType:request.ContentType}
        s, err := this.SpiderPool.Get(param)
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }
        go this.do(param, s)
    }
}

func (this *SpiderService) do(param spider.Param, s *spider.Spider) {
    if this.State == StopState {
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
        if this.State == StopState {
            break
        }
        this.EventPublisher <- u
    }
    return
}

func CreateSpiderService() (spiderService *SpiderService) {
    spiderService = &SpiderService{}
    spiderService.SpiderPool = pool.New()
    spiderService.EventListener = make(chan msg.SpiderRequest)
    return
}