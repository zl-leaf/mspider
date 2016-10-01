package downloader
import(
    "testing"
)

func TestAutoID(t *testing.T) {
    ida := autoID()
    idb := autoID()
    if ida == idb {
        t.Errorf("autoID conflict, ida %s, idb %s", ida, idb)
    }
}

func TestRequest(t *testing.T) {
    downloader,err := New()
    if err != nil {
        t.Errorf("create downloader error.")
        t.Log(err)
    }

    url := "http://hao.jobbole.com/python-scrapy/"
    _,err = downloader.Request(url)
    if err != nil {
        t.Error(err)
    }
}