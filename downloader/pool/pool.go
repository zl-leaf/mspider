package pool
import(
    "container/list"
    "time"
    "github.com/zl-leaf/mspider/downloader"
)

const(
    timeWait = 1
)

type Pool struct {
    Total int
    data *list.List
}

func New(total int) *Pool {
    pool := &Pool{Total:total, data:list.New()}
    return pool
}

func (this *Pool) Get() *downloader.Downloader {
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

func (this *Pool) Put(d *downloader.Downloader) bool {
    if this.Total > 0 && this.data.Len() >= this.Total {
        return false
    }
    this.data.PushBack(d)
    return true
}