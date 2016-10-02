package downloader
import(
    "strconv"
    "math/rand"
)

func autoID() string {
    prefix := "downloader"
    id := prefix + strconv.Itoa(rand.Int())
    return id
}