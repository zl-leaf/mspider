package engine
import (
    "github.com/zl-leaf/mspider/engine/service"
    "github.com/zl-leaf/mspider/scheduler"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
)

type Engine struct {
    SchedulerService *service.SchedulerService
    DownloaderService *service.DownloaderService
    SpiderService *service.SpiderService
}

func (this *Engine) Init() {
    this.SchedulerService = service.CreateSchedulerService()
    this.DownloaderService = service.CreateDownloaderService()
    this.SpiderService = service.CreateSpiderService()

    this.SchedulerService.Listener = this.SpiderService
    this.DownloaderService.Listener = this.SchedulerService
    this.SpiderService.Listener = this.DownloaderService
}

func (this *Engine) SetScheduler(s *scheduler.Scheduler) {
    this.SchedulerService.SetScheduler(s)
}

func (this *Engine) AddDownloader(d *downloader.Downloader) {
    this.DownloaderService.Pool.Put(d)
    logger.Info(logger.SYSTEM, "add Downloader, id %s.", d.ID)
}

func (this *Engine) AddSpider(s *spider.Spider) {
    this.SpiderService.Pool.Put(s)
    logger.Info(logger.SYSTEM, "add spider, id %s.", s.ID)
}

func (this *Engine) Start() {
    this.SchedulerService.Start()
    this.DownloaderService.Start()
    this.SpiderService.Start()
}

func (this *Engine) Stop() {
    this.SchedulerService.Stop()
    this.DownloaderService.Stop()
    this.SpiderService.Stop()
}