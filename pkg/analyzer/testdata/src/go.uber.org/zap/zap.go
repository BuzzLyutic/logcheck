// Пакет zap это минимальная заглушка для тестирования logcheck analyzer.
package zap

// Field это заглушка для zap.Field.
type Field struct{}

// Logger это заглушка для *zap.Logger.
type Logger struct{}

// NewNop возвращает no-op Logger.
func NewNop() *Logger { return &Logger{} }

func (l *Logger) Debug(msg string, fields ...Field) {}
func (l *Logger) Info(msg string, fields ...Field)  {}
func (l *Logger) Warn(msg string, fields ...Field)  {}
func (l *Logger) Error(msg string, fields ...Field) {}
func (l *Logger) DPanic(msg string, fields ...Field) {}
func (l *Logger) Panic(msg string, fields ...Field) {}
func (l *Logger) Fatal(msg string, fields ...Field) {}

// SugaredLogger это заглушка для *zap.SugaredLogger.
type SugaredLogger struct{}

// Sugar конвертирует Logger в SugaredLogger.
func (l *Logger) Sugar() *SugaredLogger { return &SugaredLogger{} }

func (s *SugaredLogger) Debugw(msg string, keysAndValues ...interface{}) {}
func (s *SugaredLogger) Infow(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Warnw(msg string, keysAndValues ...interface{})  {}
func (s *SugaredLogger) Errorw(msg string, keysAndValues ...interface{}) {}
