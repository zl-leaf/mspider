package service
import(
    "fmt"
    "time"
    "encoding/json"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/logger"
)

type DownloadResponse struct {
    URL string `json:"url"`
    Html string `json:"html"`
}

type DownloaderService struct {
    mDownloaders map[string]*downloader.Downloader
    engineListener chan string
    eventPublisher chan string
    listener IService
    state int
}

func (this *DownloaderService) EventPublisher() chan string {
    return this.eventPublisher
}

func (this *DownloaderService) Start() error {
    this.state = WorkingState
    go this.listen(this.listener.EventPublisher())
    return nil
}

func (this *DownloaderService) Stop() error {
    this.state = StopState
    stopChan := make(chan string)
    go func(stopChan chan string) {
        for _,d := range this.mDownloaders {
            logger.Info("downloader id: %s wait for stop", d.ID())
            if d.State() != downloader.FreeState {
                for {
                    time.Sleep(time.Duration(1) * time.Second)
                    if d.State() == downloader.FreeState {
                        break
                    }
                }
            }
            logger.Info("downloader id: %s has stop", d.ID())
        }
        stopChan <- "stop"
    }(stopChan)
    <- stopChan
    return nil
}

func (this *DownloaderService) AddListener(s IService) error {
    this.listener = s
    return nil
}

func (this *DownloaderService) AddDownloader(d *downloader.Downloader) {
    this.mDownloaders[d.ID()] = d
}

func (this *DownloaderService) listen(listenerChan chan string) {
    for {
        value := <- listenerChan
        go this.do(value)
    }
}

func (this *DownloaderService) do(u string) {
    if this.state == StopState {
        return
    }
    d,err := this.getDownloader()
    if err != nil {
        return
    }
    html,err := d.Request(u)
    defer d.Relase()
    logger.Info("downloader id: %s download url: %s.", d.ID(), u)
    if err != nil {
        return
    }
    resp := DownloadResponse{URL:u, Html:html}
    respJson,err := json.Marshal(resp)
    if err == nil {
        this.eventPublisher <- string(respJson)
    }
    return
}

func (this *DownloaderService) getDownloader() (dr *downloader.Downloader, err error) {
    findResult := false
    for _,d := range this.mDownloaders {
        if d.State() == downloader.FreeState {
            findResult = true
            dr = d
            break
        }
    }
    if !findResult {
        err = fmt.Errorf("can not find suitable downloader")
    }
    return
}

func CreateDownloaderService(engineListener chan string) (downloaderService *DownloaderService) {
    downloaderService = &DownloaderService{}
    downloaderService.mDownloaders = make(map[string]*downloader.Downloader, 0)
    downloaderService.engineListener = engineListener
    downloaderService.eventPublisher = make(chan string)
    return
}