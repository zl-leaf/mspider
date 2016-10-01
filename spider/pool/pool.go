package pool
import(
    "time"
    "fmt"
    "github.com/zl-leaf/mspider/spider"
)

const(
    timeWait = 1
)

type SpiderPool struct {
    data map[string]*spider.Spider
    state map[string]bool
}

func New() *SpiderPool {
    pool := new(SpiderPool)
    pool.data = make(map[string]*spider.Spider)
    pool.state = make(map[string]bool)
    return pool
}

func (this *SpiderPool) Get(u string) (targetSpider *spider.Spider, err error) {
    for {
        matchResult := false
        for id, s := range this.data {
            if match := s.MatchRules(u); match {
                free, _ := this.state[id]
                matchResult = true
                if free {
                    targetSpider = s
                    this.state[id] = false
                    break
                }
            }
        }
        if !matchResult {
            err = fmt.Errorf("can not find suitable spider for url %s", u)
            break
        }
        if targetSpider != nil {
            break
        }
        time.Sleep(time.Duration(timeWait) * time.Second)
    }
    return
}

func (this *SpiderPool) Put(s *spider.Spider) {
    if _, ok := this.data[s.ID]; !ok {
        this.data[s.ID] = s
    }
    this.state[s.ID] = true
}

func (this *SpiderPool) All() map[string]*spider.Spider {
    return this.data
}

func (this *SpiderPool) States() map[string]bool {
    return this.state
}