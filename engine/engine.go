package engine
import (
    "github.com/zl-leaf/mspider/engine/service"
    "github.com/zl-leaf/mspider/scheduler"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/spider"
)

type Engine struct {
    mSchedulerService *service.SchedulerService
    mDownloaderService *service.DownloaderService
    mSpiderService *service.SpiderService
}

func (this *Engine)Init() {
    mScheduler,_ := scheduler.New()
    ch := make(chan string)
    this.mSchedulerService = service.CreateSchedulerService(ch, mScheduler)
    this.mDownloaderService = service.CreateDownloaderService(ch)
    this.mSpiderService = service.CreateSpiderService(ch)

    this.mSchedulerService.AddListener(this.mSpiderService)
    this.mDownloaderService.AddListener(this.mSchedulerService)
    this.mSpiderService.AddListener(this.mDownloaderService)
}

func (this *Engine)Load() {
    d,_ := downloader.New("")
    this.mDownloaderService.AddDownloader(d)
}

func (this *Engine)AddSpider(s *spider.Spider) {
    this.mSpiderService.AddSpider(s)
}

func (this *Engine)Start() {
    this.mSchedulerService.Start()
    this.mDownloaderService.Start()
    this.mSpiderService.Start()
}