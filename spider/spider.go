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
    if len(this.Heart.Rules) == 0 {
        return []Rule{Rule{Match:".*"}}
    }
    return this.Heart.Rules
}

func (this *Spider) StartURLs() []string {
    return this.Heart.StartURLs
}

func (this *Spider) Do(param Param) error {
    if _,err := url.Parse(param.URL); err != nil {
        return fmt.Errorf("url:%s is illegal", param.URL)
    }
    this.Param = param
    rule, ok := this.MatchRules(param)
    if !ok {
        return fmt.Errorf("url:%s can not find the spider rule", param.URL)
    }
    if rule.Callback != nil {
        return rule.Callback(param)
    }
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

func (this *Spider) MatchRules(param Param) (result Rule, match bool) {
    rules := this.Rules()
    for _,rule := range rules {
        matchURL,_ := regexp.MatchString(rule.Match, param.URL)
        matchContentType := false
        if rule.ContentType == "" || strings.Index(param.ContentType, rule.ContentType) >= 0 {
            matchContentType = true
        }
        if matchURL && matchContentType {
            result = rule
            match = true
            break
        }
    }
    return
}