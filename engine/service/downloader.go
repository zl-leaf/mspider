package service
import(
    "fmt"
    "time"
    "github.com/zl-leaf/mspider/engine/msg"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/logger"
)

const(
    stopDownloaderWait = 1
    getDownloaderRetryNum = 3
    getDownloaderRetryWait = 1
)

type DownloaderService struct {
    Downloaders map[string]*downloader.Downloader
    EventPublisher chan msg.DownloadResult
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
            logger.Info(logger.SYSTEM, "downloader id: %s wait for stop", d.ID)
            if d.State != downloader.FreeState {
                for {
                    time.Sleep(time.Duration(stopDownloaderWait) * time.Second)
                    if d.State == downloader.FreeState {
                        break
                    }
                }
            }
            logger.Info(logger.SYSTEM, "downloader id: %s has stop", d.ID)
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
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }

        go this.do(request)
    }
}

func (this *DownloaderService) do(u string) {
    if this.State == StopState {
        return
    }
    d,err := this.getDownloader()
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
        return
    }
    html,err := d.Request(u)
    defer d.Relase()
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
        return
    }
    logger.Info(logger.SYSTEM, "downloader id: %s download url: %s.", d.ID, u)
    if this.State == StopState {
        return
    }
    err = this.response(u, html)
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
    }
    return
}

func (this *DownloaderService) response(u, html string) error {
    response := msg.DownloadResult{URL:u, Html:html}
    response, err := this.MessageHandler.HandleResponse(response)
    if err != nil {
        return err
    }

    this.EventPublisher <- response
    return nil
}

func (this *DownloaderService) getDownloader() (dr *downloader.Downloader, err error) {
    findResult := false
    for i := 0; i < getDownloaderRetryNum; i++ {
        for _,d := range this.Downloaders {
            if d.State == downloader.FreeState {
                findResult = true
                dr = d
                break
            }
        }
        if findResult {
            break
        } else {
            time.Sleep(time.Duration(getDownloaderRetryWait) * time.Second)
        }
    }

    if !findResult {
        err = fmt.Errorf("can not find free downloader")
    }
    return
}

func CreateDownloaderService() (downloaderService *DownloaderService) {
    downloaderService = &DownloaderService{}
    downloaderService.Downloaders = make(map[string]*downloader.Downloader, 0)
    downloaderService.EventPublisher = make(chan msg.DownloadResult)
    downloaderService.MessageHandler = &msg.DownloaderMessageHandler{}
    return
}