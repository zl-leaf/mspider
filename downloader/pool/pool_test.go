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

    p := New()
    p.Put(da)
    p.Put(db)

    a := p.Get()
    b := p.Get()

    if a.ID != ida || b.ID != idb {
        t.Errorf("should got ida:%s and idb:%s, but got ida:%s and idb%s", ida, idb, a.ID, b.ID)
    }
}