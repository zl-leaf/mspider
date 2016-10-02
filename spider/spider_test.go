package spider
import(
    "testing"
)

var testSpider *Spider

func TestNew(t *testing.T) {
    heart := &Heart{
        StartURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        Rules : []Rule{
            Rule{Match:"jobbole.*"},
            Rule{Match:"test.*", ContentType:"html"},
            },
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
    testaURL := "http://www.jobbole.com";
    matResult := testSpider.MatchRules(Param{URL:testaURL})
    if !matResult {
        t.Errorf("url can not match spider rules, url:%s", testaURL)
    }

    matResult = testSpider.MatchRules(Param{URL:testaURL, ContentType:"text/html; charset=utf8"})
    if !matResult {
        t.Errorf("url can not match spider rules, url:%s and contentType:%s", testaURL, "text/html; charset=utf8")
    }

    testbURL := "http://test.com"
    matResult = testSpider.MatchRules(Param{URL:testbURL, ContentType:"text/html; charset=utf8"})
    if !matResult {
        t.Errorf("url can not match spider rules, url:%s and contentType:%s", testbURL, "text/html; charset=utf8")
    }

    matResult = testSpider.MatchRules(Param{URL:testbURL, ContentType:"image/jpeg"})
    if matResult {
        t.Errorf("url should not match spider rules, url:%s and contentType:%s", testbURL, "image/jpeg")
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