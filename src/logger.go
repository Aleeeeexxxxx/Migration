package src

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
)

const ZModule = "module"

/*
 * Gin middleware
 */

type ctxCorrelationID struct{}

func WithContext(ctx context.Context, logger *zap.Logger) *zap.Logger {
	co := ctx.Value(&ctxCorrelationID{})
	if co != nil {
		logger = logger.With(zap.String("correlation_id", co.(string)))
	}
	return logger
}

func AddCorrelationID(ctx *gin.Context) {
	newCtx := context.WithValue(ctx.Request.Context(), ctxCorrelationID{}, uuid.New().String())
	ctx.Request = ctx.Request.WithContext(newCtx)
}

/*
 * Default logger
 */

var defaultLogger *zap.Logger
var once sync.Once
var level *zap.AtomicLevel

func SetDefaultLoggerLevel(_level zapcore.Level) {
	level.SetLevel(_level)
}

func GetDefaultLogger() *zap.Logger {
	once.Do(func() {
		cfg := zap.NewProductionConfig()
		level = &cfg.Level

		defaultLogger, _ = cfg.Build()
	})
	return defaultLogger
}

/*
 * Customer logger for gorm
 */

const slowSQLThreshold = 1 // second

type CustomLogger struct {
	logger *zap.Logger
}

func NewCustomLogger() *CustomLogger {
	return &CustomLogger{
		logger: GetDefaultLogger().With(zap.String(ZModule, "gorm")),
	}
}

func (l *CustomLogger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l *CustomLogger) Info(_ context.Context, msg string, data ...interface{}) {
	l.logger.Sugar().Infof(msg, data...)
}

func (l *CustomLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	l.logger.Sugar().Warnf(msg, data...)
}

func (l *CustomLogger) Error(_ context.Context, msg string, data ...interface{}) {
	l.logger.Sugar().Errorf(msg, data...)
}

func (l *CustomLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.logger.Error(
			"SQL Error",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.String("duration", elapsed.String()),
			zap.Error(err),
		)
	} else if elapsed.Seconds() > slowSQLThreshold {
		l.logger.Warn(
			"slow sql",
			zap.String("sql", sql),
			zap.Int64("rows", rows),
			zap.Float64("elapsed", elapsed.Seconds()),
		)
	}
}
