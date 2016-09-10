package service
import(
    "fmt"
    "encoding/json"
    "github.com/zl-leaf/mspider/spider"
)

type SpiderService struct {
    mSpiders map[string]spider.ISpider
    engineListener chan string
    eventPublisher chan string
    listeners []chan string
}

func (this *SpiderService) EventPublisher() chan string {
    return this.eventPublisher
}

func (this *SpiderService) Start() error {
    for _,listener := range this.listeners {
        go this.listen(listener)
    }

    for _,s := range this.mSpiders {
        for _,u := range s.StartURLs() {
            this.eventPublisher <- u
        }
    }
    return nil
}

func (this *SpiderService) AddListener(s IService) error {
    this.listeners = append(this.listeners, s.EventPublisher())
    return nil
}

func (this *SpiderService) AddSpider(s spider.ISpider) {
    this.mSpiders[s.ID()] = s
}

func (this *SpiderService) listen(listener chan string) {
    for {
        value := <- listener
        this.do(value)
    }
}

func (this *SpiderService) do(content string) {
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
    s.Parse()
    redirects := s.Redirects()
    for _,redirect := range redirects {
        this.eventPublisher <- redirect
    }
    s.Relase()
}

func (this *SpiderService) getSpider(u string) (targetSpider spider.ISpider, err error) {
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
    spiderService.mSpiders = make(map[string]spider.ISpider, 0)
    spiderService.engineListener = engineListener
    spiderService.eventPublisher = make(chan string)
    spiderService.listeners = make([]chan string, 0)
    return
}