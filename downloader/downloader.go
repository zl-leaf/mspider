package downloader
import(
    "net/http"
    "io/ioutil"
    "fmt"
)

const (
    retryNum = 2
    retryWait = 1
)

type Downloader struct {
    ID string
}

type Result struct {
    Data []byte
    ContentType string
}

func New() (downloader *Downloader, err error) {
    downloaderID := autoID()
    downloader = &Downloader{ID:downloaderID}
    return
}

func (this *Downloader) Request(u string) (result Result, err error) {
    var resp *http.Response
    downloadSuccess := false
    for i := 0; i < retryNum; i++ {
        resp, err = this.download(u)
        if err == nil {
            downloadSuccess = true
            break
        }
    }
    if !downloadSuccess {
        return
    }
    defer resp.Body.Close()

    statusCode := resp.StatusCode
    if statusCode != 200 {
        err = fmt.Errorf("get url:%s error, response statusCode not 200", u)
        return
    }

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        err = fmt.Errorf("get url:%s error, content IO read error", u)
        return
    }
    result.Data = b
    result.ContentType = http.DetectContentType(b)
    return
}

func (this *Downloader) download(u string) (resp *http.Response, err error) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", u, nil)
    if err != nil {
        return
    }
    req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko`)
    resp, err = client.Do(req)
    if err != nil {
        return
    }
    return
}