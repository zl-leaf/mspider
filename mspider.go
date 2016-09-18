package mspider
import(
    "os"
    "time"
    "github.com/zl-leaf/mspider/config"
    "github.com/zl-leaf/mspider/engine"
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

    if mConfig.LogPath != "" {
        logger.SetLogPath(mConfig.LogPath)
    } else {
        logger.SetLogPath("./")
    }
    return nil
}

func (this *MSpider) RegisterSpider(s *spider.Spider) {
    this.Engine.AddSpider(s)
}

func (this *MSpider)Start() {
    this.Engine.Start()
}

func (this *MSpider)Stop() {
    this.Engine.Stop()
    os.Exit(0)
}