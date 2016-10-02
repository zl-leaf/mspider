package spider
import(
    "testing"
    "errors"
)

var testSpider *Spider

func Callback(param Param) error {
    return errors.New("test call back ok")
}

func TestNew(t *testing.T) {
    heart := &Heart{
        StartURLs : []string{"http://hao.jobbole.com/python-scrapy"},
        Rules : []Rule{
            Rule{Match:"jobbole.*"},
            Rule{Match:"test.*", ContentType:"html"},
            Rule{Match:"callback.*", Callback:Callback},
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
    _, ok := testSpider.MatchRules(Param{URL:testaURL})
    if !ok {
        t.Errorf("url can not match spider rules, url:%s", testaURL)
    }

    _, ok = testSpider.MatchRules(Param{URL:testaURL, ContentType:"text/html; charset=utf8"})
    if !ok {
        t.Errorf("url can not match spider rules, url:%s and contentType:%s", testaURL, "text/html; charset=utf8")
    }

    testbURL := "http://test.com"
    _, ok = testSpider.MatchRules(Param{URL:testbURL, ContentType:"text/html; charset=utf8"})
    if !ok {
        t.Errorf("url can not match spider rules, url:%s and contentType:%s", testbURL, "text/html; charset=utf8")
    }

    _, ok = testSpider.MatchRules(Param{URL:testbURL, ContentType:"image/jpeg"})
    if ok {
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

func TestCallback(t *testing.T) {
    err := testSpider.Do(Param{URL:"http://jobbole.com"})
    if err != nil {
        t.Error("spider can not call the parse function")
    }

    testURL := "http://callback.com"
    err = testSpider.Do(Param{URL:testURL})
    if err == nil || err.Error() != "test call back ok" {
        t.Error("spider can not call the callback function")
    }
}