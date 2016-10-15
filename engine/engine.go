package engine
import (
    "github.com/zl-leaf/mspider/scheduler"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
)

type Engine struct {
    SchedulerService *SchedulerService
    DownloaderService *DownloaderService
    SpiderService *SpiderService
}

func (this *Engine) Init() {
    this.SchedulerService = CreateSchedulerService()
    this.DownloaderService = CreateDownloaderService()
    this.SpiderService = CreateSpiderService()

    this.SpiderService.EventPublisher = this.SchedulerService.EventListener
    this.SchedulerService.EventPublisher = this.DownloaderService.EventListener
    this.DownloaderService.EventPublisher = this.SpiderService.EventListener
}

func (this *Engine) SetScheduler(s *scheduler.Scheduler) {
    this.SchedulerService.SetScheduler(s)
}

func (this *Engine) AddDownloader(d *downloader.Downloader) {
    this.DownloaderService.DownloaderPool.Put(d)
    logger.Info(logger.SYSTEM, "add Downloader, id %s.", d.ID)
}

func (this *Engine) AddSpider(s *spider.Spider) {
    this.SpiderService.SpiderPool.Put(s)
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