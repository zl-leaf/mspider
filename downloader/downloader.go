package downloader
import(
    "net/http"
    "io/ioutil"
    "errors"
    "strings"
)

const (
    FreeState = 0
    WorkingState = 1
)

type Downloader struct {
    id string
    state int
}

func New(id string) (downloader *Downloader, err error) {
    downloaderID := id
    if downloaderID == "" {
        downloaderID = autoID()
    }
    downloader = &Downloader{id:downloaderID}
    return
}

func (this *Downloader)ID() string {
    return this.id
}

func (this *Downloader) State() int {
    return this.state
}

func (this *Downloader)Request(u string) (html string, err error) {
    this.state = WorkingState
    client := &http.Client{}
    req, err := http.NewRequest("GET", u, nil)
    req.Header.Add("User-Agent", `Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; rv:11.0) like Gecko`)
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()

    statusCode := resp.StatusCode
    if statusCode != 200 {
        err = errors.New("response statusCode not 200")
        return
    }

    b, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        err = errors.New("contentType IO read error")
        return
    }
    contentType := strings.ToLower(http.DetectContentType(b))
    if strings.Index(contentType,"text/html" ) == -1  {
        err = errors.New("contentType not html, got " + contentType)
        return
    }

    html = string(b)
    return
}

func (this *Downloader) Relase() error {
    this.state = FreeState
    return nil
}