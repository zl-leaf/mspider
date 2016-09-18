package service
import(
    "fmt"
    "time"
    "encoding/json"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
)

type SpiderService struct {
    Spiders map[string]*spider.Spider
    EventPublisher chan string
    Listener *DownloaderService
    State int
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

func (this *SpiderService) AddSpider(s *spider.Spider) {
    this.Spiders[s.ID()] = s
}

func (this *SpiderService) listen(listenerChan chan string) {
    for {
        value := <- listenerChan
        go this.do(value)
    }
}

func (this *SpiderService) do(content string) {
    if this.State == StopState {
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
        this.EventPublisher <- redirect
    }
}

func (this *SpiderService) getSpider(u string) (targetSpider *spider.Spider, err error) {
    matchResult := false
    for _,s := range this.Spiders {
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

func CreateSpiderService() (spiderService *SpiderService) {
    spiderService = &SpiderService{}
    spiderService.Spiders = make(map[string]*spider.Spider, 0)
    spiderService.EventPublisher = make(chan string)
    return
}