package spider
import(
    "testing"
)

var testSpider *Spider

func TestNew(t *testing.T) {
    heart := &Heart{
        StartURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        Rules : []string{"jobbole.*"},
        Parse: func(param Param) error {
            return nil
        },
    }
    s, err := New(heart)
    if err != nil {
        t.Error(err)
    }

    testSpider = s
}

func TestRules(t *testing.T) {
    testURL := "www.jobbole.com";
    matResult := testSpider.MatchRules(testURL)
    if !matResult {
        t.Errorf("url can not match spider rules, url %s", testURL)
    }
}

func TestGetRedirectURL(t *testing.T) {
    param := Param{}
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
        <a href="a.html">testa</a>
        <h2 class="header"></h2>
        <a href="b.html">testb</a>
    </body>
</html>`
    param.URL = "http://www.test.com"
    param.Data = []byte(testHtml)
    param.ContentType = "text/html; charset=utf8"
    testSpider.Do(param)
    redirects := testSpider.Redirects()
    if len(redirects) != 2 {
        t.Errorf("redirects len error, got %d", len(redirects))
        return
    }

    aURL := "http://www.test.com/a.html"
    bURL := "http://www.test.com/b.html"
    if redirects[0] != aURL || redirects[1] != bURL {
        t.Errorf("redirects should 0:%s, 1:%s, but got 0:%s, 1:%s", aURL, bURL, redirects[0], redirects[1])
    }
}