package spider
import(
    "testing"
)

type TestSpider struct {
    Spider
}

var testSpider *TestSpider = &TestSpider{}

func TestNewSpider(t *testing.T) {
    spiderID := "testSpider"
    err := testSpider.Init(spiderID, []string{"http://hao.jobbole.com/python-scrapy"}, []string{"jobbole.*"})

    if err != nil {
        t.Error(err)
    }

    if testSpider.ID() != spiderID {
        t.Errorf("got spider error, should got %s, but got %s", spiderID, testSpider.ID())
    }
}

func TestRules(t *testing.T) {
    testURL := "www.jobbole.com";
    matResult,_ := testSpider.MatchRules(testURL)
    if !matResult {
        t.Errorf("url can not match spider rules, url %s", testURL)
    }
}