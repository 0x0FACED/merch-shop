package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/0x0FACED/merch-shop/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	log *zap.Logger

	cfg config.LoggerConfig
}

func New(cfg config.LoggerConfig) *ZapLogger {
	logDir := "logs"
	if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
		panic("Error creating logs dir: " + err.Error())
	}

	currentDate := time.Now().Format(time.DateOnly)
	logFileName := fmt.Sprintf("%s-%s.log", "merch-shop", currentDate)
	logFilePath := filepath.Join(logDir, logFileName)

	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Error opening log file with err: " + err.Error())
	}

	cEnc := zapcore.NewConsoleEncoder(cmdConfig())
	fEnc := zapcore.NewJSONEncoder(fileConfig())

	level := level(cfg.LogLevel)

	core := zapcore.NewTee(
		zapcore.NewCore(cEnc, zapcore.AddSync(os.Stdout), zapcore.Level(level)),
		zapcore.NewCore(fEnc, zapcore.AddSync(file), zapcore.Level(level)),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))

	return &ZapLogger{
		log: logger,
		cfg: cfg,
	}
}

// просто вернет конфиг для запа для writer console
func cmdConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// просто вернет конфиг для запа для writer file
func fileConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		NameKey:       "logger",
		CallerKey:     "caller",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime:    zapcore.RFC3339TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
}

func level(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	case "panic":
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("[2006-01-02 | 15:04:05]"))
}

func (z *ZapLogger) Info(msg string, args ...any) {
	z.log.Sugar().Infow(msg, args...)
}

func (z *ZapLogger) Debug(wrappedMsg string, fields ...any) {
	z.log.Sugar().Debugw(wrappedMsg, fields...)
}

func (z *ZapLogger) Error(wrappedMsg string, fields ...any) {
	z.log.Sugar().Errorw(wrappedMsg, fields...)
}

func (z *ZapLogger) Fatal(wrappedMsg string, fields ...any) {
	z.log.Sugar().Fatalw(wrappedMsg, fields...)
}
