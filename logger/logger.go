package logger
import(
    "os"
    "log"
    "time"
)

type Logger struct {
    LogPath string
}

var mLogger *Logger = &Logger{}

func (this *Logger) Info(format string, v ...interface{}) {
    logName := time.Now().Format("2006-01-02") + ".log"
    file, _ := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
    defer file.Close()
    log.SetOutput(file)
    log.SetPrefix("[INFO]")

    log.Printf(format, v ...)

}