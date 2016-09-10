package spider
import(
    "regexp"
    "strconv"
    "math/rand"
)

func autoID() string {
    prefix := "spider"
    id := prefix + strconv.Itoa(rand.Int())
    return id
}

func GetRedirectURL(html string) []string {
    redirects := make([]string, 0)

    // 去除注释
    notesRegexp := regexp.MustCompile(`(<\!\-\-)[\s\S]*?(\-\->)`)
    html = string(notesRegexp.ReplaceAll([]byte(html), []byte("")))

    hrefRegexp := regexp.MustCompile(`<a.*?href=\"(.*?[^\"])\".*?>.*?</a>`)
    match := hrefRegexp.FindAllStringSubmatch(html, -1)
    if match != nil {
        for _,m := range match {
            redirects = append(redirects, m[1])
        }
    }
    return redirects
}