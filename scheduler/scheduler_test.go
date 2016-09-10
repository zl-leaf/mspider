package scheduler
import(
    "testing"
)

var testURL string = "https://google.com"
var testScheduler *Scheduler

func TestNew(t *testing.T) {
    s,err := New()
    if err != nil {
        t.Errorf("create scheduler error.")
        t.Log(err)
    }
    testScheduler = s
}

func TestAdd(t *testing.T) {
    testScheduler.Add(testURL)
    testScheduler.Add("https://github.com")
    testScheduler.Add("https://baidu.com")
    if testScheduler.queue.Len() != 3 {
        t.Errorf("scheduler queue length error. got %d", testScheduler.queue.Len())
    }
}

func TestGet(t *testing.T) {
    url, err := testScheduler.Head()
    if err != nil {
        t.Error(err)
    }
    if url != testURL {
        t.Errorf("scheduler get error. Should get %s, but got %s", testURL, url)
    }
}