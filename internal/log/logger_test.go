package log

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name: "基本設定",
			config: Config{
				Enabled:       true,
				Level:         InfoLevel,
				Format:        TextFormat,
				ConsoleOutput: true,
			},
			wantErr: false,
		},
		{
			name: "無効ロギング",
			config: Config{
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "JSONフォーマット",
			config: Config{
				Enabled:       true,
				Level:         DebugLevel,
				Format:        JSONFormat,
				ConsoleOutput: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := NewZapLogger(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, logger)

			// ロギング動作の確認（エラーが起きなければOK）
			logger.Info("test info message")
			logger.Debug("test debug message")
			logger.Warn("test warning message")
			logger.Error("test error message")
		})
	}
}

func TestFileOutput(t *testing.T) {
	// テスト用の一時ディレクトリを作成
	tempDir, err := os.MkdirTemp("", "logger-test")
	require.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Failed to remove temporary directory %s: %v", tempDir, err)
		}
	}()

	logFile := filepath.Join(tempDir, "test.log")

	config := Config{
		Enabled:       true,
		Level:         InfoLevel,
		File:          logFile,
		MaxSize:       1,
		MaxAge:        1,
		MaxFiles:      2,
		Format:        TextFormat,
		ConsoleOutput: false,
	}

	// ロガーを初期化
	logger, err := NewZapLogger(config)
	require.NoError(t, err)

	// ログを出力
	logger.Info("file output test",
		zap.String("key1", "value1"),
		zap.Int("key2", 123),
	)
	logger.Warn("warning message")

	// バッファをフラッシュ
	require.NoError(t, logger.Sync())

	// ファイルが作成されたことを確認
	_, err = os.Stat(logFile)
	assert.NoError(t, err)

	// ファイル内容を読み込み
	content, err := os.ReadFile(logFile)
	require.NoError(t, err)

	// 期待するログメッセージが含まれているか確認
	contentStr := string(content)
	assert.Contains(t, contentStr, "file output test")
	assert.Contains(t, contentStr, "key1")
	assert.Contains(t, contentStr, "value1")
	assert.Contains(t, contentStr, "warning message")
}

func TestGlobalLogger(t *testing.T) {
	// テスト前にグローバルロガーをリセット（他のテストの影響を受けないように）
	mu.Lock()
	globalLogger = nil
	mu.Unlock()

	// デフォルト設定でのロガー取得
	logger := GetLogger()
	assert.NotNil(t, logger)

	// カスタム設定でグローバルロガーを初期化
	customConfig := Config{
		Enabled:       true,
		Level:         DebugLevel,
		Format:        JSONFormat,
		ConsoleOutput: true,
	}

	err := InitGlobalLogger(customConfig)
	assert.NoError(t, err)

	// グローバルロガーが初期化されたことを確認
	logger = GetLogger()
	assert.NotNil(t, logger)

	// ロギング動作の確認（エラーが起きなければOK）
	Info("global logger info test")
	Debug("global logger debug test")
	Warn("global logger warn test")

	// グローバルロガーをフラッシュ
	// 標準出力へのSyncはエラーを返す場合があるため、結果は無視する
	_ = Sync()
}

func TestLoggerWithInvalidConfig(t *testing.T) {
	// Zapロガーは基本的に設定エラーでパニックしないため、
	// 実装上、設定が無効でもエラーは返さない
	// 代わりにデフォルト値やNOPロガーを使用する

	config := Config{
		Enabled: true,
		Level:   "invalid", // 無効なレベル
		Format:  "invalid", // 無効なフォーマット
	}

	// 設定を検証
	err := config.Validate()
	assert.Error(t, err)

	// 実際のところ、Validateが失敗しても、NewZapLoggerは
	// デフォルト値を使用してロガーを作成する
	logger, err := NewZapLogger(Config{
		Enabled:       true,
		Level:         "invalid_level", // 内部でデフォルトのInfoLevelに変換
		Format:        TextFormat,
		ConsoleOutput: true,
	})

	// エラーは発生せず、デフォルト値を使ってロガーが作成される
	assert.NoError(t, err)
	assert.NotNil(t, logger)
}

// テスト用のWriteSyncerを実装
type testWriteSyncer struct {
	logs []string
}

func (ws *testWriteSyncer) Write(p []byte) (int, error) {
	ws.logs = append(ws.logs, string(p))
	return len(p), nil
}

func (ws *testWriteSyncer) Sync() error {
	return nil
}

func TestLogLevels(t *testing.T) {
	// 異なるログレベルでのフィルタリングをテスト

	// テスト用のWriteSyncer
	ws := &testWriteSyncer{}

	// カスタムのコア作成
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "msg",
		LevelKey:   "level",
		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,
	}
	encoder := zapcore.NewConsoleEncoder(encoderConfig)

	// InfoLevelのロガーを作成
	infoLevel := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	core := zapcore.NewCore(encoder, zapcore.AddSync(ws), infoLevel)
	logger := zap.New(core)

	// 各レベルでログを出力
	logger.Debug("debug message") // 出力されないはず
	logger.Info("info message")   // 出力される
	logger.Warn("warn message")   // 出力される
	logger.Error("error message") // 出力される

	// 結果確認
	assert.Equal(t, 3, len(ws.logs))

	// WarnLevelに変更してテスト
	infoLevel.SetLevel(zapcore.WarnLevel)

	// 再度テスト
	logger.Debug("debug message again") // 出力されない
	logger.Info("info message again")   // 出力されない
	logger.Warn("warn message again")   // 出力される
	logger.Error("error message again") // 出力される

	// 最終的には5つのログ（Info以上が3つ + Warn以上が2つ）
	assert.Equal(t, 5, len(ws.logs))
	assert.Contains(t, ws.logs[0], "info message")
	assert.Contains(t, ws.logs[1], "warn message")
	assert.Contains(t, ws.logs[2], "error message")
	assert.Contains(t, ws.logs[3], "warn message again")
	assert.Contains(t, ws.logs[4], "error message again")
}
