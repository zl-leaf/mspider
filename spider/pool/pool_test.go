package pool
import(
    "testing"
    "github.com/zl-leaf/mspider/spider"
)

func TestPool(t *testing.T) {
    hearta := &spider.Heart{
        StartURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        Rules : []spider.Rule{
            spider.Rule{Match:"jobbole.*"},
            },
        Parse: func(param spider.Param) error {
            return nil
        },
    }
    sa,_ := spider.New(hearta)

    heartb := &spider.Heart{
        StartURLs : []string{"http://baidu.com"},
        Rules : []spider.Rule{
            spider.Rule{Match:"baidu.*", ContentType:"image"},
            },
        Parse: func(param spider.Param) error {
            return nil
        },
    }
    sb,_ := spider.New(heartb)

    pool := New()

    pool.Put(sa)
    pool.Put(sb)

    a, err := pool.Get(spider.Param{URL:"jobbole.com"})
    if err != nil {
        t.Error(err)
        return
    }
    if a.ID != sa.ID {
        t.Errorf("spider should got %s, but got %s", sa.ID, a.ID)
    }

    b, err := pool.Get(spider.Param{URL:"baidu.com", ContentType:"image/jpeg"})
    if err != nil {
        t.Error(err)
        return
    }

    if b.ID != sb.ID {
        t.Errorf("spider should got %s, but got %s", sb.ID, b.ID)
    }

    if _, err := pool.Get(spider.Param{URL:"test.com"}); err == nil {
        t.Errorf("test.com should got error, but not now")
    }
}