package service
import(
    "fmt"
    "time"
    "encoding/json"
    "github.com/zl-leaf/mspider/engine/msg"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/logger"
)

type DownloaderService struct {
    Downloaders map[string]*downloader.Downloader
    EventPublisher chan string
    Listener *SchedulerService
    State int
    MessageHandler msg.IDownloaderMessageHandler
}

func (this *DownloaderService) Start() error {
    this.State = WorkingState
    go this.listen(this.Listener.EventPublisher)
    return nil
}

func (this *DownloaderService) Stop() error {
    this.State = StopState
    stopChan := make(chan string)
    go func(stopChan chan string) {
        for _,d := range this.Downloaders {
            logger.Info("downloader id: %s wait for stop", d.ID)
            if d.State != downloader.FreeState {
                for {
                    time.Sleep(time.Duration(1) * time.Second)
                    if d.State == downloader.FreeState {
                        break
                    }
                }
            }
            logger.Info("downloader id: %s has stop", d.ID)
        }
        stopChan <- "stop"
    }(stopChan)
    <- stopChan
    return nil
}

func (this *DownloaderService) AddDownloader(d *downloader.Downloader) {
    this.Downloaders[d.ID] = d
}

func (this *DownloaderService) listen(listenerChan chan string) {
    for {
        value := <- listenerChan
        request, err := this.MessageHandler.HandleRequest(value)
        if err == nil {
            go this.do(request)
        }
    }
}

func (this *DownloaderService) do(u string) {
    if this.State == StopState {
        return
    }
    d,err := this.getDownloader()
    if err != nil {
        return
    }
    html,err := d.Request(u)
    defer d.Relase()
    logger.Info("downloader id: %s download url: %s.", d.ID, u)
    if err != nil {
        return
    }
    this.response(u, html)
    return
}

func (this *DownloaderService) response(u, html string) error {
    resp := msg.DownloadResponse{URL:u, Html:html}
    resp, err := this.MessageHandler.HandleResponse(resp)
    if err != nil {
        return err
    }
    respJson,err := json.Marshal(resp)
    if err != nil {
        return err
    }
    this.EventPublisher <- string(respJson)
    return nil
}

func (this *DownloaderService) getDownloader() (dr *downloader.Downloader, err error) {
    findResult := false
    for _,d := range this.Downloaders {
        if d.State == downloader.FreeState {
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

func CreateDownloaderService() (downloaderService *DownloaderService) {
    downloaderService = &DownloaderService{}
    downloaderService.Downloaders = make(map[string]*downloader.Downloader, 0)
    downloaderService.EventPublisher = make(chan string)
    downloaderService.MessageHandler = &msg.DownloaderMessageHandler{}
    return
}