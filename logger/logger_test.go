package logger
import(
    "testing"
    "os"
)

func TestInfo(t *testing.T) {
    logPath := "log"

    os.Mkdir(logPath, 0777)
    SetLogPath(logPath)
    Info(SYSTEM, "test")
    Info(SYSTEM, "test2")
    Info(USER, "test %d", 3)
    Error(USER, "test4")
    os.RemoveAll(logPath)
}