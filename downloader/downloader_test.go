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
    downloaderID := "testDownloader"
    downloader,err := New(downloaderID)
    if err != nil {
        t.Errorf("create downloader error.")
        t.Log(err)
    }
    if downloader.ID() != downloaderID {
        t.Errorf("downloader ID error. got %s", downloader.ID)
    }

    url := "http://hao.jobbole.com/python-scrapy/"
    html,err := downloader.Request(url)
    if err != nil {
        t.Error(err)
    }

    t.Log(html)
}