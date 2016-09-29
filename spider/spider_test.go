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

func (this *TestSpiderHeart) Parse(url, content string) error {
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

func TestGetRedirectURL(t *testing.T) {
    testHtml := `<!doctype html>
<html>
    <head>
        <meta a="b">
    </head>
    <body>
        <p><!-- this is a comment -->
        This is some text.
        </p>
        <div></div>
        <h1 class="header"></h1>
        <a href="aURL">testa</a>
        <h2 class="header"></h2>
        <a href="bURL">testb</a>
    </body>
</html>`
    redirects, err := GetRedirectURL(testHtml)
    if err != nil {
        t.Error(err)
        return
    }
    if len(redirects) != 2 {
        t.Errorf("redirects len error, got %d", len(redirects))
        return
    }

    if redirects[0] != "aURL" || redirects[1] != "bURL" {
        t.Errorf("redirects should 0:aURL, 1:bURL, but got 0:%s, 1:%s", redirects[0], redirects[1])
    }
}