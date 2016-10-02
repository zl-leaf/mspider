package spider
import(
    "bytes"
    "strings"
    "net/http"
    "strconv"
    "errors"
    "math/rand"
    "container/list"
    "golang.org/x/net/html"
)

func autoID() string {
    prefix := "spider"
    id := prefix + strconv.Itoa(rand.Int())
    return id
}

func IsHtml(data []byte) bool{
    contentType := strings.ToLower(http.DetectContentType(data))
    if strings.Index(contentType,"html" ) == -1  {
        return false
    }
    return true
}

func GetRedirectURL(data []byte) (redirects []string, err error) {
    if !IsHtml(data) {
        err = errors.New("it is not html")
        return
    }
    doc, err := html.Parse(bytes.NewReader(data))
    if err != nil {
        return
    }
    q := list.New()
    q.PushBack(doc)
    for {
        if q.Len() == 0 {
            break
        }
        e := q.Front()
        q.Remove(e)

        node := e.Value.(*html.Node)
        for c := node.FirstChild; c != nil; c = c.NextSibling {
            if c.Type == html.ElementNode {
                q.PushBack(c)
            }
        }

        if node.Type == html.ElementNode && node.Data == "a" {
            for _,a := range node.Attr {
                if a.Key == "href" {
                    redirects = append(redirects, a.Val)
                }
            }
        }
    }
    return
}