package mspider
import(
    "os"
    "time"
    "github.com/zl-leaf/mspider/config"
    "github.com/zl-leaf/mspider/engine"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
)

const(
    defaultLogPath = "./"
    defaultSpiderInterval = int(time.Second)
    defaultDownloaderNum = 1
)

type MSpider struct {
    Engine *engine.Engine
}

func New() (m *MSpider, err error) {
    m = &MSpider{}
    m.init()
    return
}

func (this *MSpider) init() {
    this.Engine = &engine.Engine{}
    this.Engine.Init()
}

func (this *MSpider) Load(mConfig *config.Config) error {
    // load logger
    if mConfig.LogPath != "" {
        logger.SetLogPath(mConfig.LogPath)
    } else {
        logger.SetLogPath(defaultLogPath)
    }

    // load scheduler service
    if mConfig.SpiderInterval > 0 {
        this.Engine.SchedulerService.Scheduler.Interval = mConfig.SpiderInterval
    } else {
        this.Engine.SchedulerService.Scheduler.Interval = defaultSpiderInterval
    }

    // load download service
    var downloaderNum = defaultDownloaderNum
    if mConfig.DownloaderNum > 0 {
        downloaderNum = mConfig.DownloaderNum
    }
    this.Engine.DownloaderService.DownloaderPool.Total = downloaderNum
    for i := 0; i < downloaderNum; i++ {
        d,_ := downloader.New()
        if !this.Engine.DownloaderService.DownloaderPool.Put(d) {
            break
        }
    }
    return nil
}

func (this *MSpider) RegisterSpider(s *spider.Spider) {
    this.Engine.AddSpider(s)
}

func (this *MSpider) SetSchedulerValidate(validate func(request string) error) {
    this.Engine.SchedulerService.Validate = validate
}

func (this *MSpider) SetDownloaderValidator(validate func(request string) error) {
    this.Engine.DownloaderService.Validate = validate
}

func (this *MSpider) SetSpiderValidate(validate func(request spider.Param) error) {
    this.Engine.SpiderService.Validate = validate
}

func (this *MSpider)Start() {
    this.Engine.Start()
}

func (this *MSpider)Stop() {
    this.Engine.Stop()
    os.Exit(0)
}