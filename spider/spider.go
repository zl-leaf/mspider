package spider
import(
    "regexp"
)

const (
    FreeState = 0
    WorkingState = 1
)

type ISpider interface {
    ID() string
    Do(u string,content string) error
    Relase() error
    Parse() error
    Rules() []string
    StartURLs() []string
    Redirects() []string
    State() int
    MatchRules(u string) (bool,error)
}

type Spider struct {
    id string
    startURLs []string
    rules []string
    url string
    html string
    state int
}


func (this *Spider) Init(id string, startURLs []string, rules []string) error {
    spiderID := id
    if spiderID == "" {
        spiderID = autoID()
    }
    this.id = id
    this.startURLs = startURLs
    if len(rules) > 0 {
        this.rules = rules
    } else {
        this.rules = []string{"*"}
    }
    this.state = FreeState
    return nil
}

func (this *Spider) ID() string {
    return this.id
}

func (this *Spider) Rules() []string {
    return this.rules
}

func (this *Spider) StartURLs() []string {
    return this.startURLs
}

func (this *Spider) Do(u string, content string) error {
    this.url = u
    this.html = content
    this.state = WorkingState
    return nil
}

func (this *Spider) Relase() error {
    this.state = FreeState
    return nil
}

func (this *Spider) Parse() error {
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
    for _,rule := range this.rules {
        if r,_ := regexp.MatchString(rule, u); r {
            result = true
            break
        }
    }
    return
}