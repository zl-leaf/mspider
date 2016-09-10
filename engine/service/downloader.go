package service
import(
    "fmt"
    "encoding/json"
    "github.com/zl-leaf/mspider/downloader"
)

type DownloadResponse struct {
    URL string `json:"url"`
    Html string `json:"html"`
}

type DownloaderService struct {
    mDownloaders map[string]*downloader.Downloader
    engineListener chan string
    eventPublisher chan string
    listeners []chan string
}

func (this *DownloaderService) EventPublisher() chan string {
    return this.eventPublisher
}

func (this *DownloaderService) Start() error {
    for _,listener := range this.listeners {
        go this.listen(listener)
    }
    return nil
}

func (this *DownloaderService) AddListener(s IService) error {
    this.listeners = append(this.listeners, s.EventPublisher())
    return nil
}

func (this *DownloaderService) AddDownloader(d *downloader.Downloader) {
    this.mDownloaders[d.ID()] = d
}

func (this *DownloaderService) listen(listener chan string) {
    for {
        value := <- listener
        this.do(value)
    }
}

func (this *DownloaderService) do(u string) {
    d,err := this.getDownloader()
    if err != nil {
        return
    }
    html,err := d.Request(u)
    if err != nil {
        return
    }
    resp := DownloadResponse{URL:u, Html:html}
    respJson,err := json.Marshal(resp)
    if err == nil {
        this.eventPublisher <- string(respJson)
    }
    d.Relase()
    return
}

func (this *DownloaderService) getDownloader() (dr *downloader.Downloader, err error) {
    findResult := false
    for _,d := range this.mDownloaders {
        if d.State() == downloader.FreeState {
            findResult = true
            dr = d
            break
        }
    }
    if !findResult {
        err = fmt.Errorf("can not find suitable downloader")
    }
    return
}

func CreateDownloaderService(engineListener chan string) (downloaderService *DownloaderService) {
    downloaderService = &DownloaderService{}
    downloaderService.mDownloaders = make(map[string]*downloader.Downloader, 0)
    downloaderService.engineListener = engineListener
    downloaderService.eventPublisher = make(chan string)
    downloaderService.listeners = make([]chan string, 0)
    return
}