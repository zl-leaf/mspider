package downloader
import(
    "strconv"
    "math/rand"
    "strings"
    "net/http"
)

func autoID() string {
    prefix := "downloader"
    id := prefix + strconv.Itoa(rand.Int())
    return id
}

func checkContentType(content []byte) bool {
    contentType := strings.ToLower(http.DetectContentType(content))
    if strings.Index(contentType,"text/html" ) == -1  {
        return false
    }
    return true
}