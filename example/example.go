package main
import (
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/engine"
)

type DemoSpiderHeart struct {
    startURLs []string
    rules []string
}

func (this *DemoSpiderHeart) StartURLs() []string {
    return this.startURLs
}

func (this *DemoSpiderHeart) Rules() []string {
    return this.rules
}

func (this *DemoSpiderHeart)Parse() error {
    // TODO
    return nil
}

func main() {
    e := &engine.Engine{}
    e.Init()
    e.Load()

    heart := &DemoSpiderHeart{
        startURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        rules : []string{"jobbole.*"},
    }
    spider,_ := spider.New("", heart)
    e.AddSpider(spider)

    e.Start()

    ch := make(chan string)
    <-ch
}