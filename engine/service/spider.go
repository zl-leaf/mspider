package service
import(
    "fmt"
    "time"
    "encoding/json"
    "github.com/zl-leaf/mspider/engine/msg"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
)

const(
    stopSpiderWait = 1
    getSpiderRetryNum = 3
    getSpiderRetryWait = 1
)

type SpiderService struct {
    Spiders map[string]*spider.Spider
    EventPublisher chan string
    Listener *DownloaderService
    State int
    MessageHandler msg.ISpiderMessageHandler
}

func (this *SpiderService) Start() error {
    this.State = WorkingState
    go this.listen(this.Listener.EventPublisher)

    for _,s := range this.Spiders {
        for _,u := range s.StartURLs() {
            this.EventPublisher <- u
        }
    }
    return nil
}

func (this *SpiderService) Stop() error {
    this.State = StopState
    stopChan := make(chan string)
    go func(stopChan chan string) {
        for _,s := range this.Spiders {
            logger.Info(logger.SYSTEM, "spider id: %s wait for stop", s.ID)
            if s.State != spider.FreeState {
                for {
                    time.Sleep(time.Duration(stopSpiderWait) * time.Second)
                    if s.State == spider.FreeState {
                        break
                    }
                }
            }
            logger.Info(logger.SYSTEM, "spider id: %s has stop", s.ID)
        }
        stopChan <- "stop"
    }(stopChan)
    <- stopChan
    return nil
}

func (this *SpiderService) AddSpider(s *spider.Spider) {
    this.Spiders[s.ID] = s
}

func (this *SpiderService) listen(listenerChan chan string) {
    for {
        value := <- listenerChan
        var requestData msg.DownloadResponse
        err := json.Unmarshal([]byte(value), &requestData)
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }
        request, err := this.MessageHandler.HandleRequest(requestData)
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }
        go this.do(request)
    }
}

func (this *SpiderService) do(request msg.DownloadResponse) {
    if this.State == StopState {
        return
    }

    s, err := this.getSpider(request.URL)
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
        return
    }
    err = s.Do(request.URL, request.Html)
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
        return
    }
    defer s.Relase()
    logger.Info(logger.SYSTEM, "spider id: %s crawl url: %s.", s.ID, request.URL)
    redirects, err := this.MessageHandler.HandleResponse(s.Redirects())
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
        return
    }
    for _,redirect := range redirects {
        if this.State == StopState {
            break
        }
        this.EventPublisher <- redirect
    }
}

func (this *SpiderService) getSpider(u string) (targetSpider *spider.Spider, err error) {
    findResult := false
    for i := 0; i < getSpiderRetryNum; i++ {
        matchResult := false
        for _,s := range this.Spiders {
            if matchResult,_ = s.MatchRules(u); matchResult {
                if s.State == spider.FreeState {
                    targetSpider = s
                    findResult = true
                    break
                }
            }
        }

        if !matchResult {
            err = fmt.Errorf("can not find suitable spider for url %s", u)
            break
        }
        if findResult {
            break
        } else {
            time.Sleep(time.Duration(getSpiderRetryWait) * time.Second)
        }
    }

    if !findResult {
        err = fmt.Errorf("can not find free spider")
    }
    return
}

func CreateSpiderService() (spiderService *SpiderService) {
    spiderService = &SpiderService{}
    spiderService.Spiders = make(map[string]*spider.Spider, 0)
    spiderService.EventPublisher = make(chan string)
    spiderService.MessageHandler = &msg.SpiderMessageHandler{}
    return
}