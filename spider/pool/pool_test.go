package pool
import(
    "testing"
    "github.com/zl-leaf/mspider/spider"
)

type TestSpiderHeart struct {
    startURLs []string
    rules []string
}

func (this *TestSpiderHeart) StartURLs() []string {
    return this.startURLs
}

func (this *TestSpiderHeart) Rules() []string {
    return this.rules
}

func (this *TestSpiderHeart) Parse(url, content string) error {
    return nil
}

func TestPool(t *testing.T) {
    hearta := &TestSpiderHeart{
        startURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        rules : []string{"jobbole.*"},
    }
    sa,_ := spider.New(hearta)

    heartb := &TestSpiderHeart{
        startURLs : []string{"http://baidu.com"},
        rules : []string{"baidu.*"},
    }
    sb,_ := spider.New(heartb)

    pool := New()

    pool.Put(sa)
    pool.Put(sb)

    a, err := pool.Get("jobbole.com")
    if err != nil {
        t.Error(err)
        return
    }
    if a.ID != sa.ID {
        t.Errorf("spider should got %s, but got %s", sa.ID, a.ID)
    }

    b, err := pool.Get("baidu.com")
    if err != nil {
        t.Error(err)
        return
    }

    if b.ID != sb.ID {
        t.Errorf("spider should got %s, but got %s", sb.ID, b.ID)
    }

    if _, err := pool.Get("test.com"); err == nil {
        t.Errorf("test.com should got error, but not now")
    }
}