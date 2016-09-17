package engine
import (
    "github.com/zl-leaf/mspider/engine/service"
    "github.com/zl-leaf/mspider/scheduler"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
)

type Engine struct {
    mSchedulerService *service.SchedulerService
    mDownloaderService *service.DownloaderService
    mSpiderService *service.SpiderService
}

func (this *Engine) Init() {
    ch := make(chan string)
    this.mSchedulerService = service.CreateSchedulerService(ch)
    this.mDownloaderService = service.CreateDownloaderService(ch)
    this.mSpiderService = service.CreateSpiderService(ch)

    this.mSchedulerService.AddListener(this.mSpiderService)
    this.mDownloaderService.AddListener(this.mSchedulerService)
    this.mSpiderService.AddListener(this.mDownloaderService)
}

func (this *Engine) SetScheduler(s *scheduler.Scheduler) {
    this.mSchedulerService.SetScheduler(s)
}

func (this *Engine) AddDownloader(d *downloader.Downloader) {
    this.mDownloaderService.AddDownloader(d)
    logger.Info("add Downloader, id %s.", d.ID())
}

func (this *Engine) AddSpider(s *spider.Spider) {
    this.mSpiderService.AddSpider(s)
    logger.Info("add spider, id %s.", s.ID())
}

func (this *Engine) Start() {
    this.mSchedulerService.Start()
    this.mDownloaderService.Start()
    this.mSpiderService.Start()
}

func (this *Engine) Stop() {
    this.mSchedulerService.Stop()
    this.mDownloaderService.Stop()
    this.mSpiderService.Stop()
}