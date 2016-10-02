## mspider
mspider是一个网络爬虫的框架（参考scrapy），通过自定义爬虫，可以抓取各个html
## Quick Start
###### 下载与安装
    go get github.com/zl-leaf/mspider

###### 创建文件 `main.go`
```go
package main
import (
    "time"
    "github.com/zl-leaf/mspider/config"
    "github.com/zl-leaf/mspider/spider"
    "github.com/zl-leaf/mspider/logger"
    "github.com/zl-leaf/mspider"
)

func Parse(param spider.Param) error {
    // TODO
    return nil
}

func Callback(param spider.Param) error {
    logger.Info(logger.SYSTEM, "url:%s call callback function", param.URL)
    return nil
}

func main() {
    mspider,_ := mspider.New()
    c := &config.Config{DownloaderNum:2}
    mspider.Load(c)

    heart := &spider.Heart{
        StartURLs : []string{"http://myurl.com"},
        Rules : []spider.Rule{
            spider.Rule{Match:"myurl.*", ContentType:"html"},
            spider.Rule{Match:"myurl.*", ContentType:"image", Callback:Callback},
            },
        Parse: Parse,
    }
    spider,_ := spider.New(heart)
    mspider.RegisterSpider(spider)

    mspider.Start()

    time.Sleep(time.Duration(10) * time.Second)
    mspider.Stop()
}
```

###### 构建并运行
```bash
    go build main.go
    ./main
```
