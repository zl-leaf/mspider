package service
import(
    "fmt"
    "time"
    "encoding/json"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
)

type SpiderService struct {
    mSpiders map[string]*spider.Spider
    engineListener chan string
    eventPublisher chan string
    listener IService
    state int
}

func (this *SpiderService) EventPublisher() chan string {
    return this.eventPublisher
}

func (this *SpiderService) Start() error {
    this.state = WorkingState
    go this.listen(this.listener.EventPublisher())

    for _,s := range this.mSpiders {
        for _,u := range s.StartURLs() {
            this.eventPublisher <- u
        }
    }
    return nil
}

func (this *SpiderService) Stop() error {
    this.state = StopState
    stopChan := make(chan string)
    go func(stopChan chan string) {
        for _,s := range this.mSpiders {
            logger.Info("spider id: %s wait for stop", s.ID())
            if s.State() != spider.FreeState {
                for {
                    time.Sleep(time.Duration(1) * time.Second)
                    if s.State() == spider.FreeState {
                        break
                    }
                }
            }
            logger.Info("spider id: %s has stop", s.ID())
        }
        stopChan <- "stop"
    }(stopChan)
    <- stopChan
    return nil
}

func (this *SpiderService) AddListener(s IService) error {
    this.listener = s
    return nil
}

func (this *SpiderService) AddSpider(s *spider.Spider) {
    this.mSpiders[s.ID()] = s
}

func (this *SpiderService) listen(listenerChan chan string) {
    for {
        value := <- listenerChan
        go this.do(value)
    }
}

func (this *SpiderService) do(content string) {
    if this.state == StopState {
        return
    }
    var dresp DownloadResponse
    err := json.Unmarshal([]byte(content), &dresp)
    if err != nil {
        return
    }
    s, err := this.getSpider(dresp.URL)
    if err != nil {
        return
    }
    s.Do(dresp.URL, dresp.Html)
    defer s.Relase()
    logger.Info("spider id: %s crawl url: %s.", s.ID(), dresp.URL)
    redirects := s.Redirects()
    for _,redirect := range redirects {
        this.eventPublisher <- redirect
    }
}

func (this *SpiderService) getSpider(u string) (targetSpider *spider.Spider, err error) {
    matchResult := false
    for _,s := range this.mSpiders {
        if s.State() != spider.FreeState {
            continue
        }
        if matchResult,_ = s.MatchRules(u); matchResult {
            targetSpider = s
            break
        }
    }

    if !matchResult {
        err = fmt.Errorf("can not find suitable spider for url %s", u)
    }
    return
}

func CreateSpiderService(engineListener chan string) (spiderService *SpiderService) {
    spiderService = &SpiderService{}
    spiderService.mSpiders = make(map[string]*spider.Spider, 0)
    spiderService.engineListener = engineListener
    spiderService.eventPublisher = make(chan string)
    return
}