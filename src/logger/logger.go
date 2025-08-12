package logger

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var Suger *zap.SugaredLogger

type Config struct {
	WriteBoth string
}

// ZapConfig
// @Date: 2023-01-09 16:37:05
// @Description: zap日志配置结构体
type ZapConfig struct {
	Prefix     string         `yaml:"prefix" mapstructure:"prefix"`
	TimeFormat string         `yaml:"timeFormat" mapstructure:"timeFormat"`
	Level      string         `yaml:"level" mapstructure:"level"`
	Caller     bool           `yaml:"caller" mapstructure:"caller"`
	StackTrace bool           `yaml:"stackTrace" mapstructure:"stackTrace"`
	Writer     string         `yaml:"writer" mapstructure:"writer"`
	Encode     string         `yaml:"encode" mapstructure:"encode"`
	LogFile    *LogFileConfig `yaml:"logFile" mapstructure:"logFile"`
}

// LogFileConfig
// @Date: 2023-01-09 16:38:45
// @Description: 日志文件配置结构体
type LogFileConfig struct {
	MaxSize  int      `yaml:"maxSize" mapstructure:"maxSize"`
	BackUps  int      `yaml:"backups" mapstructure:"backups"`
	Compress bool     `yaml:"compress" mapstructure:"compress"`
	Output   []string `yaml:"output" mapstructure:"output"`
	Errput   []string `yaml:"errput" mapstructure:"errput"`
}

//	func InitLogger() {
//		logger, err := zap.NewProduction()
//		if err != nil {
//			panic("Fail to Init Logger")
//		}
//		Logger = logger
//		Suger = logger.Sugar()
//	}

func InitZap() {
	config := &ZapConfig{
		Prefix:     "ZapLogTest",
		TimeFormat: "2006/01/02 - 15:04:05.00000",
		Level:      "debug",
		Caller:     true,
		StackTrace: false,
		Writer:     "both",
		Encode:     "console",
		LogFile: &LogFileConfig{
			MaxSize:  20,
			BackUps:  5,
			Compress: true,
			Output:   []string{"./log/output.log"},
			Errput:   []string{},
		},
	}
	// 构建编码器
	encoder := zapEncoder(config)
	// 构建日志级别
	levelEnabler := zap.InfoLevel
	// 最后获得Core和Options
	subCore, options := tee(config, encoder, levelEnabler)
	// 创建Logger
	logger := zap.New(subCore, options...)
	Suger = logger.Sugar()
	Logger = logger

}

func zapEncoder(config *ZapConfig) zapcore.Encoder {
	// 新建一个配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "Time",
		LevelKey:      "Level",
		NameKey:       "Logger",
		CallerKey:     "Caller",
		MessageKey:    "Message",
		StacktraceKey: "StackTrace",
		LineEnding:    zapcore.DefaultLineEnding,
		FunctionKey:   zapcore.OmitKey,
	}
	// 自定义时间格式
	// encoderConfig.EncodeTime = CustomTimeFormatEncoder
	// 日志级别大写
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	// 秒级时间间隔
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	// 简短的调用者输出
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	// 完整的序列化logger名称
	encoderConfig.EncodeName = zapcore.FullNameEncoder
	// 最终的日志编码 json或者console
	switch config.Encode {
	case "json":
		{
			return zapcore.NewJSONEncoder(encoderConfig)
		}
	case "console":
		{
			return zapcore.NewConsoleEncoder(encoderConfig)
		}
	}
	// 默认console
	return zapcore.NewConsoleEncoder(encoderConfig)
}

// 将所有合并
func tee(cfg *ZapConfig, encoder zapcore.Encoder, levelEnabler zapcore.LevelEnabler) (core zapcore.Core, options []zap.Option) {
	sink := zapWriteSyncer(cfg)
	return zapcore.NewCore(encoder, sink, levelEnabler), buildOptions(cfg, levelEnabler)
}

// 构建Option
func buildOptions(cfg *ZapConfig, levelEnabler zapcore.LevelEnabler) (options []zap.Option) {
	if cfg.Caller {
		options = append(options, zap.AddCaller())
	}

	if cfg.StackTrace {
		options = append(options, zap.AddStacktrace(levelEnabler))
	}
	return
}

func zapWriteSyncer(cfg *ZapConfig) zapcore.WriteSyncer {
	syncers := make([]zapcore.WriteSyncer, 0, 2)
	// 如果开启了日志控制台输出，就加入控制台书写器
	// if cfg.Writer == config.WriteBoth || cfg.Writer == config.WriteConsole {
	syncers = append(syncers, zapcore.AddSync(os.Stdout))
	// }

	// 如果开启了日志文件存储，就根据文件路径切片加入书写器
	// if cfg.Writer == config.WriteBoth || cfg.Writer == config.WriteFile {
	// 添加日志输出器
	for _, path := range cfg.LogFile.Output {
		logger := &lumberjack.Logger{
			Filename:   path,                 //文件路径
			MaxSize:    cfg.LogFile.MaxSize,  //分割文件的大小
			MaxBackups: cfg.LogFile.BackUps,  //备份次数
			Compress:   cfg.LogFile.Compress, // 是否压缩
			LocalTime:  true,                 //使用本地时间
		}
		syncers = append(syncers, zapcore.Lock(zapcore.AddSync(logger)))
		// }
	}
	return zap.CombineWriteSyncers(syncers...)
}

// func zapLevelEnabler(cfg *ZapConfig) zapcore.LevelEnabler {
// 	switch cfg.Level {
// 	case config.DebugLevel:
// 		return zap.DebugLevel
// 	case config.InfoLevel:
// 		return zap.InfoLevel
// 	case config.ErrorLevel:
// 		return zap.ErrorLevel
// 	case config.PanicLevel:
// 		return zap.PanicLevel
// 	case config.FatalLevel:
// 		return zap.FatalLevel
// 	}
// 	// 默认Debug级别
// 	return zap.DebugLevel
// }

type ZapBotLogger struct{}

func (l *ZapBotLogger) Println(v ...any) {
	Suger.Info(v...)
}

func (l *ZapBotLogger) Printf(format string, v ...any) {
	Suger.Infof(format, v...)
}
