package pool
import(
    "container/list"
    "time"
    "github.com/zl-leaf/mspider/downloader"
)

const(
    timeWait = 1
)

type DownloaderPool struct {
    data *list.List
}

func New() *DownloaderPool {
    pool := &DownloaderPool{data:list.New()}
    return pool
}

func (this *DownloaderPool) Get() *downloader.Downloader {
    var result *downloader.Downloader
    for {
        if this.data.Len() > 0 {
            e := this.data.Front()
            this.data.Remove(e)
            result = e.Value.(*downloader.Downloader)
            break
        }
        time.Sleep(time.Duration(timeWait) * time.Second)
    }
    return result
}

func (this *DownloaderPool) Put(d *downloader.Downloader) {
    this.data.PushBack(d)
}