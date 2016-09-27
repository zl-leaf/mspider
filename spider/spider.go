package spider
import(
    "regexp"
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
    redirects := GetRedirectURL(this.Html)
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