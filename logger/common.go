package logger

func SetLogPath(LogPath string) {
    mLogger.LogPath = LogPath
}

func Info(format string, v ...interface{}) {
    mLogger.Info(format, v...)
}

func Error(format string, v ...interface{}) {
    mLogger.Error(format, v...)
}