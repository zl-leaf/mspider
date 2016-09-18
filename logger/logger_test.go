package logger
import(
    "testing"
)

func TestInfo(t *testing.T) {
    SetLogPath("mspider.log")
    Info("test")
    Info("test2")
    Info("test %d", 3)
    Error("test4")
}