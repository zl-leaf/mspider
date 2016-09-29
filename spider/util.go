package spider
import(
    "strings"
    "strconv"
    "math/rand"
    "golang.org/x/net/html"
    "container/list"
)

func autoID() string {
    prefix := "spider"
    id := prefix + strconv.Itoa(rand.Int())
    return id
}

func GetRedirectURL(content string) (redirects []string, err error) {
    doc, err := html.Parse(strings.NewReader(content))
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