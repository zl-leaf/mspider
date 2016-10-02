package msg

type SpiderRequest struct {
    URL string
    Data []byte
    ContentType string
}