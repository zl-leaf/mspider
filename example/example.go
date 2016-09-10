package main
import (
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/engine"
    "os"
    "log"
)

type DemoSpider struct {
    spider.Spider
}


func (this *DemoSpider)Parse() error {
    file, _ := os.Create("content.log")
    logger := log.New(file, "", log.LstdFlags|log.Llongfile)
    redirects := this.Redirects()
    for _,u := range redirects {
        logger.Println(u)
    }
    return nil
}

func main() {
    e := &engine.Engine{}
    e.Init()
    e.Load()

    s := &DemoSpider{}
    s.Init("", []string{"http://hao.jobbole.com/python-scrapy/"}, []string{"jobbole.*"})
    e.AddSpider(s)

    e.Start()

    ch := make(chan string)
    <-ch
}