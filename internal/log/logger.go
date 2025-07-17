package log

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// グローバルロガーインスタンス
	globalLogger *zap.Logger
	// 複数のゴルーチンからアクセスするためのミューテックス
	mu sync.RWMutex
)

// Logger はアプリケーションで使用するロガーインターフェース
type Logger interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Fatal(msg string, fields ...zapcore.Field)
}

// NewZapLogger は設定に基づいてzapロガーを初期化する関数
func NewZapLogger(config Config) (*zap.Logger, error) {
	// ロギングが無効の場合はNopLoggerを返す
	if !config.Enabled {
		return zap.NewNop(), nil
	}

	// ログレベルの設定
	var level zapcore.Level
	switch config.Level {
	case DebugLevel:
		level = zapcore.DebugLevel
	case InfoLevel:
		level = zapcore.InfoLevel
	case WarnLevel:
		level = zapcore.WarnLevel
	case ErrorLevel:
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	atomicLevel := zap.NewAtomicLevelAt(level)

	// エンコーダー設定
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// カラー出力の設定
	if config.Color && config.ConsoleOutput {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	var encoder zapcore.Encoder
	if config.Format == JSONFormat {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 出力先の設定
	var cores []zapcore.Core

	// ファイル出力
	if config.File != "" {
		fileWriteSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.File,
			MaxSize:    config.MaxSize,
			MaxAge:     config.MaxAge,
			MaxBackups: config.MaxFiles,
			Compress:   true,
		})
		cores = append(cores, zapcore.NewCore(encoder, fileWriteSyncer, atomicLevel))
	}

	// コンソール出力
	if config.ConsoleOutput {
		var consoleEncoder zapcore.Encoder
		// コンソール出力用に別のエンコーダーを用意（カラー設定を適用するため）
		if config.Color {
			consoleConfig := encoderConfig
			consoleConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
			if config.Format == JSONFormat {
				consoleEncoder = zapcore.NewJSONEncoder(consoleConfig)
			} else {
				consoleEncoder = zapcore.NewConsoleEncoder(consoleConfig)
			}
		} else {
			consoleEncoder = encoder
		}
		consoleWriteSyncer := zapcore.AddSync(os.Stdout)
		cores = append(cores, zapcore.NewCore(consoleEncoder, consoleWriteSyncer, atomicLevel))
	}

	// コアの作成
	var core zapcore.Core
	if len(cores) > 0 {
		core = zapcore.NewTee(cores...)
	} else {
		// 出力先がない場合はNOPコア（何もしない）
		core = zapcore.NewNopCore()
	}

	// ロガーのオプション設定
	options := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	}

	// 詳細ログが有効な場合はスタックトレースを追加
	if config.Verbose {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// ロガーの作成
	logger := zap.New(core, options...)

	return logger, nil
}

// InitGlobalLogger はグローバルロガーを初期化する
func InitGlobalLogger(config Config) error {
	logger, err := NewZapLogger(config)
	if err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	// 既存のグローバルロガーがある場合は、リソースを解放
	if globalLogger != nil {
		_ = globalLogger.Sync()
	}

	globalLogger = logger
	return nil
}

// GetLogger はグローバルロガーを取得する
// グローバルロガーが初期化されていない場合はデフォルト設定で初期化する
func GetLogger() *zap.Logger {
	mu.RLock()
	if globalLogger != nil {
		defer mu.RUnlock()
		return globalLogger
	}
	mu.RUnlock()

	// グローバルロガーがない場合は初期化
	mu.Lock()
	defer mu.Unlock()

	if globalLogger == nil {
		var err error
		globalLogger, err = NewZapLogger(DefaultConfig())
		if err != nil {
			// 初期化エラーの場合は何もしないロガーを返す
			globalLogger = zap.NewNop()
		}
	}

	return globalLogger
}

// Debug はデバッグレベルのログを出力する
func Debug(msg string, fields ...zap.Field) {
	GetLogger().Debug(msg, fields...)
}

// Info は情報レベルのログを出力する
func Info(msg string, fields ...zap.Field) {
	GetLogger().Info(msg, fields...)
}

// Warn は警告レベルのログを出力する
func Warn(msg string, fields ...zap.Field) {
	GetLogger().Warn(msg, fields...)
}

// Error はエラーレベルのログを出力する
func Error(msg string, fields ...zap.Field) {
	GetLogger().Error(msg, fields...)
}

// Fatal は致命的なエラーレベルのログを出力する（その後、プロセスを終了する）
func Fatal(msg string, fields ...zap.Field) {
	GetLogger().Fatal(msg, fields...)
}

// Sync はバッファされたログをフラッシュする
func Sync() error {
	mu.RLock()
	defer mu.RUnlock()

	if globalLogger != nil {
		return globalLogger.Sync()
	}
	return nil
}
