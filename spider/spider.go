package spider
import(
    "regexp"
    "fmt"
    "net/url"
)

type SpiderHeart interface {
    StartURLs() []string
    Rules() []string
    Parse(url string, data []byte) error
}

type Spider struct {
    ID string
    URL string
    Data []byte
    Heart SpiderHeart
}

func New(heart SpiderHeart) (spider *Spider, err error) {
    spiderID := autoID()
    spider = &Spider{ID:spiderID, Heart:heart}
    return
}

func (this *Spider) Rules() []string {
    return this.Heart.Rules()
}

func (this *Spider) StartURLs() []string {
    return this.Heart.StartURLs()
}

func (this *Spider) Do(u string, data []byte) error {
    if _,err := url.Parse(u); err != nil {
        return fmt.Errorf("url:%s is illegal")
    }
    this.URL = u
    this.Data = data
    return this.Heart.Parse(u, data)
}

func (this *Spider) Redirects() []string {
    redirects := make([]string, 0)
    hrefs, err := GetRedirectURL(this.Data)
    if err != nil {
        return redirects
    }
    baseURL, _ := url.Parse(this.URL)
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

func (this *Spider) MatchRules(u string) bool {
    result := false
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
    return result
}