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
    mEngine *engine.Engine
}

func New() (m *MSpider, err error) {
    m = &MSpider{}
    m.init()
    return
}

func (this *MSpider) init() {
    this.mEngine = &engine.Engine{}
    this.mEngine.Init()
}

func (this *MSpider) Load(mConfig *config.Config) error {
    s,_ := scheduler.New()
    if mConfig.SpiderInterval > 0 {
        s.Interval = mConfig.SpiderInterval
    } else {
        s.Interval = int(time.Second)
    }
    this.mEngine.SetScheduler(s)

    for i := 0; i < mConfig.DownloaderNum; i++ {
        d,_ := downloader.New()
        this.mEngine.AddDownloader(d)
    }

    if mConfig.LogPath != "" {
        logger.SetLogPath(mConfig.LogPath)
    } else {
        logger.SetLogPath("./")
    }
    return nil
}

func (this *MSpider) RegisterSpider(s *spider.Spider) {
    this.mEngine.AddSpider(s)
}

func (this *MSpider)Start() {
    this.mEngine.Start()
}

func (this *MSpider)Stop() {
    this.mEngine.Stop()
    os.Exit(0)
}