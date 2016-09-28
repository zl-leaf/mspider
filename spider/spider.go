package spider
import(
    "regexp"
    "fmt"
    "net/url"
)

const (
    FreeState = 0
    WorkingState = 1
)

type SpiderHeart interface {
    StartURLs() []string
    Rules() []string
    Parse(url, content string) error
}

type Spider struct {
    ID string
    URL string
    Html string
    State int
    Heart SpiderHeart
}

func New(heart SpiderHeart) (spider *Spider, err error) {
    spiderID := autoID()
    spider = &Spider{ID:spiderID, State:FreeState, Heart:heart}
    return
}

func (this *Spider) Rules() []string {
    return this.Heart.Rules()
}

func (this *Spider) StartURLs() []string {
    return this.Heart.StartURLs()
}

func (this *Spider) Do(u string, content string) error {
    if _,err := url.Parse(u); err != nil {
        return fmt.Errorf("url:%s is illegal")
    }
    this.URL = u
    this.Html = content
    this.State = WorkingState
    return this.Heart.Parse(u, content)
}

func (this *Spider) Relase() error {
    this.State = FreeState
    return nil
}

func (this *Spider) Redirects() []string {
    hrefs := GetRedirectURL(this.Html)
    baseURL, _ := url.Parse(this.URL)
    redirects := make([]string, 0)
    for _, href := range hrefs {
        hrefURL, err := url.Parse(href)
        if err != nil {
            continue
        }
        redirect := baseURL.ResolveReference(hrefURL).String()
        redirects = append(redirects, redirect)
    }
    return redirects
}

func (this *Spider) MatchRules(u string) (result bool, err error) {
    rules := this.Rules()
    if len(rules) == 0 {
        result = true
    } else {
        for _,rule := range rules {
            if r,_ := regexp.MatchString(rule, u); r {
                result = true
                break
            }
        }
    }
    return
}