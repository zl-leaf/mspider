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
    EventListener chan string
    EventPublisher chan msg.SpiderRequest
    State int
    Validator msg.IDownloaderValidator
}

func (this *DownloaderService) Start() error {
    this.State = WorkingState
    go this.listen(this.EventListener)
    return nil
}

func (this *DownloaderService) Stop() error {
    this.State = StopState
    return nil
}

func (this *DownloaderService) listen(listenerChan chan string) {
    for {
        request := <- listenerChan
        if this.Validator != nil {
            if err := this.Validator.Validate(request); err != nil {
                logger.Error(logger.SYSTEM, err.Error())
                continue
            }
        }

        d := this.Pool.Get()
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
    response := msg.SpiderRequest{URL:u, Html:html}
    this.EventPublisher <- response
    return
}

func CreateDownloaderService() (downloaderService *DownloaderService) {
    downloaderService = &DownloaderService{}
    downloaderService.Pool = pool.New()
    downloaderService.EventListener = make(chan string, 10)
    return
}