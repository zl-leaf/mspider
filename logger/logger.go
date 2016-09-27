package logger
import(
    "os"
    "log"
    "path/filepath"
)

type Logger struct {
    Level string
    LogPath string
}

var systemLogger *Logger = &Logger{Level:SYSTEM}
var userLogger *Logger = &Logger{Level:USER}

func (this *Logger) Info(format string, v ...interface{}) {
    actualPath := filepath.Join(this.LogPath, this.LogName())
    file, _ := openLogFile(actualPath)
    defer file.Close()
    log.SetOutput(file)
    log.SetPrefix("[INFO]")

    log.Printf(format, v ...)
}

func (this *Logger) Error(format string, v ...interface{}) {
    actualPath := filepath.Join(this.LogPath, this.LogName())
    file, _ := openLogFile(actualPath)
    defer file.Close()
    log.SetOutput(file)
    log.SetPrefix("[ERROR]")

    log.Printf(format, v ...)
}

func (this *Logger) LogName() string {
    return this.Level + ".log"
}

func openLogFile(path string) (*os.File, error) {
    file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
    return file, err
}