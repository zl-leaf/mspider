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
    Parse() error
}

type Spider struct {
    id string
    url string
    html string
    state int
    heart SpiderHeart
}

func New(id string, heart SpiderHeart) (spider *Spider, err error) {
    spiderID := id
    if spiderID == "" {
        spiderID = autoID()
    }

    spider = &Spider{id:spiderID, state:FreeState, heart:heart}
    return
}

func (this *Spider) ID() string {
    return this.id
}

func (this *Spider) Rules() []string {
    return this.heart.Rules()
}

func (this *Spider) StartURLs() []string {
    return this.heart.StartURLs()
}

func (this *Spider) Do(u string, content string) error {
    this.url = u
    this.html = content
    this.state = WorkingState
    return this.heart.Parse()
}

func (this *Spider) Relase() error {
    this.state = FreeState
    return nil
}

func (this *Spider) Redirects() []string {
    redirects := GetRedirectURL(this.html)
    return redirects
}

func (this *Spider) State() int {
    return this.state
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