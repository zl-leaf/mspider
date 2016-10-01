package service
import(
    "github.com/zl-leaf/mspider/engine/msg"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/downloader/pool"
    "github.com/zl-leaf/mspider/logger"
)

const(
    stopDownloaderWait = 1
    getDownloaderRetryNum = 3
    getDownloaderRetryWait = 1
)

type DownloaderService struct {
    Pool *pool.DownloaderPool
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
    return nil
}

func (this *DownloaderService) listen(listenerChan chan string) {
    for {
        value := <- listenerChan
        request, err := this.MessageHandler.HandleRequest(value)
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }

        d := this.Pool.Get()
        if err != nil {
            logger.Error(logger.SYSTEM, err.Error())
            continue
        }
        go this.do(request, d)
    }
}

func (this *DownloaderService) do(u string, d *downloader.Downloader) {
    if this.State == StopState {
        return
    }
    html,err := d.Request(u)
    defer this.Pool.Put(d)
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

func CreateDownloaderService() (downloaderService *DownloaderService) {
    downloaderService = &DownloaderService{}
    downloaderService.Pool = pool.New()
    downloaderService.EventPublisher = make(chan msg.DownloadResult)
    downloaderService.MessageHandler = &msg.DownloaderMessageHandler{}
    return
}