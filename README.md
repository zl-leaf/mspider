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
    "github.com/zl-leaf/mspider"
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

func (this *DemoSpiderHeart)Parse(url, content string) error {
    // TODO
    return nil
}

func main() {
    mspider,_ := mspider.New()
    c := &config.Config{DownloaderNum:2}
    mspider.Load(c)

    heart := &DemoSpiderHeart{
        startURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        rules : []string{"jobbole.*"},
    }
    spider,_ := spider.New("", heart)
    mspider.RegisterSpider(spider)

    mspider.Start()

    time.Sleep(time.Duration(10) * time.Second)// 10秒后停止抓取
    mspider.Stop()

}
```

###### 设置消息组件（在自定义的组件中进行数据验证或格式设置）
```go
type SchedulerMessageHandler struct {}
func (this *SchedulerMessageHandler) HandleRequest(req string) (value string, err error) {
    value = req
    return
}
func (this *SchedulerMessageHandler) HandleResponse(resp string) (value string, err error) {
    value = resp
    return
}

type DownloaderMessageHandler struct {}
func (this *DownloaderMessageHandler) HandleRequest(req string) (value string, err error) {
    value = req
    return
}
func (this *DownloaderMessageHandler) HandleResponse(resp DownloadResult) (value DownloadResult, err error) {
    value = resp
    return
}

type SpiderMessageHandler struct {}
func (this *SpiderMessageHandler) HandleRequest(req DownloadResult) (value DownloadResult, err error) {
    value = req
    return
}
func (this *SpiderMessageHandler) HandleResponse(resp []string) (value []string, err error) {
    value = resp
    return
}

mspider.SetSchedulerMessageHandler(&SchedulerMessageHandler{})
mspider.SetDownloaderMessageHandler(&DownloaderMessageHandler{})
mspider.SetSpiderMessageHandler(&SpiderMessageHandler{})
```

###### 构建并运行
```bash
    go build main.go
    ./main
```
