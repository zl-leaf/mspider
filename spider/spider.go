package spider
import(
    "regexp"
    "fmt"
    "strings"
    "net/url"
)

type Rule struct {
    Match string
    ContentType string
    Callback func(param Param) error
}

type Param struct {
    URL string
    Data []byte
    ContentType string
}

type Heart struct {
    StartURLs []string
    Rules []Rule
    Parse func(param Param) error
}

type Spider struct {
    ID string
    Param Param
    Heart *Heart
}

func New(heart *Heart) (spider *Spider, err error) {
    spiderID := autoID()
    spider = &Spider{ID:spiderID, Heart:heart}
    return
}

func (this *Spider) Rules() []Rule {
    return this.Heart.Rules
}

func (this *Spider) StartURLs() []string {
    return this.Heart.StartURLs
}

func (this *Spider) Do(param Param) error {
    if _,err := url.Parse(param.URL); err != nil {
        return fmt.Errorf("url:%s is illegal")
    }
    this.Param = param
    return this.Heart.Parse(param)
}

func (this *Spider) Redirects() []string {
    redirects := make([]string, 0)
    hrefs, err := GetRedirectURL(this.Param.Data)
    if err != nil {
        return redirects
    }
    baseURL, _ := url.Parse(this.Param.URL)

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

func (this *Spider) MatchRules(param Param) bool {
    result := false
    rules := this.Rules()
    if len(rules) == 0 {
        result = true
    } else {
        for _,rule := range rules {
            matchURL,_ := regexp.MatchString(rule.Match, param.URL)
            matchContentType := false
            if rule.ContentType == "" || strings.Index(param.ContentType, rule.ContentType) >= 0 {
                matchContentType = true
            }
            if matchURL && matchContentType {
                result = true
                break
            }
        }
    }
    return result
}