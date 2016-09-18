package msg

type DownloadResponse struct {
    URL string `json:"url"`
    Html string `json:"html"`
}

type ISchedulerMessageHandler interface {
    HandleRequest(req string) (string, error)
    HandleResponse(resp string) (string, error)
}

type IDownloaderMessageHandler interface {
    HandleRequest(req string) (string, error)
    HandleResponse(resp DownloadResponse) (DownloadResponse, error)
}

type ISpiderMessageHandler interface {
    HandleRequest(req DownloadResponse) (DownloadResponse, error)
    HandleResponse(resp []string) ([]string, error)
}

type SchedulerMessageHandler struct {}

func (this *SchedulerMessageHandler) HandleRequest(req string) (value string, err error) {
    value = req
    return
}

func (this *SchedulerMessageHandler) HandleResponse(resp string) (value string, err error) {
    value = resp
    return
}

type DownloaderMessageHandler struct {}

func (this *DownloaderMessageHandler) HandleRequest(req string) (value string, err error) {
    value = req
    return
}

func (this *DownloaderMessageHandler) HandleResponse(resp DownloadResponse) (value DownloadResponse, err error) {
    value = resp
    return
}

type SpiderMessageHandler struct {}

func (this *SpiderMessageHandler) HandleRequest(req DownloadResponse) (value DownloadResponse, err error) {
    value = req
    return
}

func (this *SpiderMessageHandler) HandleResponse(resp []string) (value []string, err error) {
    value = resp
    return
}