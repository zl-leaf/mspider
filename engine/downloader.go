package engine
import(
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/downloader/pool"
    "github.com/zl-leaf/mspider/logger"
)

type DownloaderService struct {
    DownloaderPool *pool.Pool
    EventListener chan string
    EventPublisher chan spider.Param
    State IState
    Validate func(request string) error
}

func (this *DownloaderService) Start() error {
    nextState, err := this.State.Next(WorkingStateCode)
    if err != nil {
        return err
    }
    this.State = nextState
    go this.listen(this.EventListener)
    return nil
}

func (this *DownloaderService) Stop() error {
    nextState, err := this.State.Next(FreeStateCode)
    if err != nil {
        return err
    }
    this.State = nextState
    return nil
}

func (this *DownloaderService) listen(listenerChan chan string) {
    for {
        request := <- listenerChan
        if this.Validate != nil {
            if err := this.Validate(request); err != nil {
                logger.Error(logger.SYSTEM, err.Error())
                continue
            }
        }

        d := this.DownloaderPool.Get()
        go this.do(request, d)
    }
}

func (this *DownloaderService) do(u string, d *downloader.Downloader) {
    if this.State.Code() != WorkingStateCode {
        return
    }
    result,err := d.Request(u)
    defer this.DownloaderPool.Put(d)
    if err != nil {
        logger.Error(logger.SYSTEM, err.Error())
        return
    }
    logger.Info(logger.SYSTEM, "downloader id: %s download url: %s.", d.ID, u)
    if this.State.Code() != WorkingStateCode {
        return
    }
    response := spider.Param{URL:u, Data:result.Data, ContentType:result.ContentType}
    this.EventPublisher <- response
    return
}

func CreateDownloaderService() (downloaderService *DownloaderService) {
    downloaderService = &DownloaderService{}
    downloaderService.DownloaderPool = pool.New(0)
    downloaderService.EventListener = make(chan string, 10)
    downloaderService.State = createState(FreeStateCode)
    return
}