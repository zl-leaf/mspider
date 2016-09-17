package logger
import(
    "testing"
    "time"
    "os"
    "log"
)

func TestInfo(t *testing.T) {
    SetLogPath("mspider.log")
    Info("test")
    Info("test2")
    Info("test %d", 3)
    Info("test4")
}

func TestGoLog(t *testing.T) {
    logName := time.Now().Format("2006-01-02") + ".log"
    file, _ := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
    defer file.Close()
    logger := log.New(file, "[INFO]", log.Llongfile)

    logger.Printf("test")
    logger.Printf("test2")
    logger.Printf("test %d", 3)
    logger.Printf("test4")
}