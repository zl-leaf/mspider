package pool
import(
    "testing"
    "github.com/zl-leaf/mspider/downloader"
)

func TestPool(t *testing.T) {
    da,_ := downloader.New()
    ida := da.ID
    db,_ := downloader.New()
    idb := db.ID

    p := New(0)
    p.Put(da)
    p.Put(db)

    a := p.Get()
    b := p.Get()

    if a.ID != ida || b.ID != idb {
        t.Errorf("should got ida:%s and idb:%s, but got ida:%s and idb%s", ida, idb, a.ID, b.ID)
    }

    p.Total = 1
    if !p.Put(a) || p.Put(b) {
        t.Error("pool shoul only put one downloader success, but it is not now")
    }
}