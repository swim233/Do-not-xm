package logger

import "strings"

type TelegramBotApiLoggerAdapter struct {
	logger              *Logger
	loglevelStrProvider func(int) string
	logLevel            int
}

// Printf(format string, v ...interface{})	// 格式化打印
// Println(v ...interface{})								// 自动fmt到string后打印
func (l *TelegramBotApiLoggerAdapter) Printf(format string, v ...interface{}) {
	l.logger.log(l.logLevel, l.loglevelStrProvider(l.logLevel), format, v...)
}

func (l *TelegramBotApiLoggerAdapter) Println(v ...interface{}) {
	builder := strings.Builder{}
	for _, item := range v {
		builder.WriteString(item.(string))
	}
	l.logger.log(l.logLevel, l.loglevelStrProvider(l.logLevel), "%s", builder.String())
}

func (l *TelegramBotApiLoggerAdapter) SetLogger(logger *Logger) {
	l.logger = logger
}

func (l *TelegramBotApiLoggerAdapter) SetLogLevel(level int) {
	l.logLevel = level
}

func (l *TelegramBotApiLoggerAdapter) SetLogLevelStrProvider(provider func(int) string) {
	l.loglevelStrProvider = provider
}
