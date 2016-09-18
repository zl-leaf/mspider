package spider
import(
    "testing"
)

type TestSpiderHeart struct {
    startURLs []string
    rules []string
}

func (this *TestSpiderHeart) StartURLs() []string {
    return this.startURLs
}

func (this *TestSpiderHeart) Rules() []string {
    return this.rules
}

func (this *TestSpiderHeart) Parse() error {
    return nil
}

func TestNewSpider(t *testing.T) {
    heart := &TestSpiderHeart{
        startURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        rules : []string{"jobbole.*"},
    }
    testSpider,err := New(heart)

    if err != nil {
        t.Error(err)
    }

    if len(testSpider.StartURLs()) == 0 {
        t.Errorf("spider starturls length error, got 0")
    }

    if len(testSpider.Rules()) == 0 {
        t.Errorf("spider rules length error, got 0")
    }
}

func TestRules(t *testing.T) {
    heart := &TestSpiderHeart{
        startURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        rules : []string{"jobbole.*"},
    }
    testSpider,_ := New(heart)
    testURL := "www.jobbole.com";
    matResult,_ := testSpider.MatchRules(testURL)
    if !matResult {
        t.Errorf("url can not match spider rules, url %s", testURL)
    }
}