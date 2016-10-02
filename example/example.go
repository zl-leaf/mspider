package main
import (
    "time"
    "github.com/zl-leaf/mspider/config"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider"
)

func Parse(param spider.Param) error {
    // TODO
    return nil
}

func main() {
    mspider,_ := mspider.New()
    c := &config.Config{DownloaderNum:2}
    mspider.Load(c)

    heart := &spider.Heart{
        StartURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        Rules : []string{"jobbole.*"},
        Parse: Parse,
    }
    spider,_ := spider.New(heart)
    mspider.RegisterSpider(spider)

    mspider.Start()

    time.Sleep(time.Duration(10) * time.Second)
    mspider.Stop()
}