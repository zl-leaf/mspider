package msg

type ISchedulerValidator interface {
    Validate(request string) error
}

type IDownloaderValidator interface {
    Validate(request string) error
}

type ISpiderValidator interface {
    Validate(request SpiderRequest) error
}