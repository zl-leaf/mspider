package mspider
import(
    "os"
    "time"
    "github.com/zl-leaf/mspider/config"
    "github.com/zl-leaf/mspider/engine"
    "github.com/zl-leaf/mspider/engine/msg"
    "github.com/zl-leaf/mspider/scheduler"
    "github.com/zl-leaf/mspider/downloader"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
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
    if mConfig.LogPath != "" {
        logger.SetLogPath(mConfig.LogPath)
    } else {
        logger.SetLogPath("./")
    }

    s,_ := scheduler.New()
    if mConfig.SpiderInterval > 0 {
        s.Interval = mConfig.SpiderInterval
    } else {
        s.Interval = int(time.Second)
    }
    this.Engine.SetScheduler(s)

    for i := 0; i < mConfig.DownloaderNum; i++ {
        d,_ := downloader.New()
        this.Engine.AddDownloader(d)
    }
    return nil
}

func (this *MSpider) RegisterSpider(s *spider.Spider) {
    this.Engine.AddSpider(s)
}

func (this *MSpider) SetSchedulerValidator(validator msg.ISchedulerValidator) {
    this.Engine.SchedulerService.Validator = validator
}

func (this *MSpider) SetDownloaderValidator(validator msg.IDownloaderValidator) {
    this.Engine.DownloaderService.Validator = validator
}

func (this *MSpider) SetSpiderValidator(validator msg.ISpiderValidator) {
    this.Engine.SpiderService.Validator = validator
}

func (this *MSpider)Start() {
    this.Engine.Start()
}

func (this *MSpider)Stop() {
    this.Engine.Stop()
    os.Exit(0)
}