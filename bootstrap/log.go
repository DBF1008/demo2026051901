package bootstrap

import (
	"gofly/global"
	"gofly/utils/gf"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	level   zapcore.Level
	options []zap.Option
)

func InitializeLog() *zap.Logger {
	createRootDir()

	setLogLevel()

	if global.App.Config.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}

	return zap.New(getZapCore(), options...)
}

func createRootDir() {
	if ok, _ := gf.PathExists(global.App.Config.Log.RootDir); !ok {
		_ = os.Mkdir(global.App.Config.Log.RootDir, os.ModePerm)
	}
}

func setLogLevel() {
	switch global.App.Config.Log.Level {
	case "debug":
		level = zap.DebugLevel
		options = append(options, zap.AddStacktrace(level))
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
		options = append(options, zap.AddStacktrace(level))
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
}

func getZapCore() zapcore.Core {
	var encoder zapcore.Encoder

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(time.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = func(l zapcore.Level, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(global.App.Config.App.Env + "." + l.String())
	}

	if global.App.Config.Log.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	return zapcore.NewCore(encoder, getLogWriter(), level)
}

func getLogWriter() zapcore.WriteSyncer {
	tiemstr := time.Now().Format("200601/02")
	file := &lumberjack.Logger{
		Filename:   global.App.Config.Log.RootDir + "/" + tiemstr + ".log",
		MaxSize:    global.App.Config.Log.MaxSize,
		MaxBackups: global.App.Config.Log.MaxBackups,
		MaxAge:     global.App.Config.Log.MaxAge,
		Compress:   global.App.Config.Log.Compress,
	}

	return zapcore.AddSync(file)
}
