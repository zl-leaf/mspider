package logger
import(
    "fmt"
)

const(
    SYSTEM = "system"
    USER = "user"
)

func SetLogPath(LogPath string) {
    systemLogger.LogPath = LogPath
    userLogger.LogPath = LogPath
}

func Info(level string, format string, v ...interface{}) {
    logger,_ := getLogger(level)
    logger.Info(format, v...)
}

func Error(level string, format string, v ...interface{}) {
    logger,_ := getLogger(level)
    logger.Error(format, v...)
}

func getLogger(level string) (logger *Logger, err error) {
    switch level {
    case SYSTEM:
        logger = systemLogger
    case USER:
        logger = userLogger
    default:
        err = fmt.Errorf("logger level:%s, is error", level)
    }
    return
}