package logger

import (
    "context"
    "os"
    "time"
    "github.com/sirupsen/logrus"
)

// Logger wraps logrus with additional functionality
type Logger struct {
    *logrus.Logger
}

// Fields type for structured logging
type Fields = logrus.Fields

// Global logger instance
var Log *Logger

// InitializeLogger sets up the global logger with configuration
func InitializeLogger(level string, format string) error {
    logger := logrus.New()
    
    // Set log level
    logLevel, err := logrus.ParseLevel(level)
    if err != nil {
        logLevel = logrus.InfoLevel
    }
    logger.SetLevel(logLevel)
    
    // Set formatter
    switch format {
    case "json":
        logger.SetFormatter(&logrus.JSONFormatter{
            TimestampFormat: time.RFC3339,
            FieldMap: logrus.FieldMap{
                logrus.FieldKeyTime:  "timestamp",
                logrus.FieldKeyLevel: "level",
                logrus.FieldKeyMsg:   "message",
                logrus.FieldKeyFunc:  "caller",
            },
        })
    default:
        logger.SetFormatter(&logrus.TextFormatter{
            FullTimestamp:   true,
            TimestampFormat: time.RFC3339,
        })
    }
    
    // Set output
    logger.SetOutput(os.Stdout)
    
    // Enable caller reporting
    logger.SetReportCaller(true)
    
    Log = &Logger{Logger: logger}
    return nil
}

// WithFields adds fields to log entry
func (l *Logger) WithFields(fields Fields) *logrus.Entry {
    return l.Logger.WithFields(logrus.Fields(fields))
}

// WithContext adds context information to log entry
func (l *Logger) WithContext(ctx context.Context) *logrus.Entry {
    entry := l.Logger.WithContext(ctx)
    
    // Add request ID if available
    if requestID := GetRequestID(ctx); requestID != "" {
        entry = entry.WithField("request_id", requestID)
    }
    
    // Add user ID if available (for authenticated requests)
    if userID := GetUserID(ctx); userID != "" {
        entry = entry.WithField("user_id", userID)
    }
    
    return entry
}

// WithError adds error information to log entry
func (l *Logger) WithError(err error) *logrus.Entry {
    return l.Logger.WithError(err)
}

// Request logging helpers
func (l *Logger) LogRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration) {
    fields := Fields{
        "method":      method,
        "path":        path,
        "status_code": statusCode,
        "duration_ms": duration.Milliseconds(),
    }
    
    entry := l.WithContext(ctx).WithFields(fields)
    
    switch {
    case statusCode >= 500:
        entry.Error("HTTP request completed with server error")
    case statusCode >= 400:
        entry.Warn("HTTP request completed with client error")
    default:
        entry.Info("HTTP request completed")
    }
}

// Database operation logging
func (l *Logger) LogDBOperation(ctx context.Context, operation, table string, duration time.Duration, err error) {
    fields := Fields{
        "operation":   operation,
        "table":       table,
        "duration_ms": duration.Milliseconds(),
    }
    
    entry := l.WithContext(ctx).WithFields(fields)
    
    if err != nil {
        entry.WithError(err).Error("Database operation failed")
    } else {
        entry.Debug("Database operation completed")
    }
}

// Cache operation logging
func (l *Logger) LogCacheOperation(ctx context.Context, operation, key string, hit bool, duration time.Duration) {
    fields := Fields{
        "operation":   operation,
        "cache_key":   key,
        "cache_hit":   hit,
        "duration_ms": duration.Milliseconds(),
    }
    
    l.WithContext(ctx).WithFields(fields).Debug("Cache operation completed")
}

// Security event logging
func (l *Logger) LogSecurityEvent(ctx context.Context, event string, details Fields) {
    allFields := Fields{
        "security_event": event,
    }
    
    // Merge additional details
    for k, v := range details {
        allFields[k] = v
    }
    
    l.WithContext(ctx).WithFields(allFields).Warn("Security event detected")
}

// Context key types
type contextKey string

const (
    RequestIDKey contextKey = "request_id"
    UserIDKey    contextKey = "user_id"
)

// GetRequestID extracts request ID from context
func GetRequestID(ctx context.Context) string {
    if id, ok := ctx.Value(RequestIDKey).(string); ok {
        return id
    }
    return ""
}

// GetUserID extracts user ID from context
func GetUserID(ctx context.Context) string {
    if id, ok := ctx.Value(UserIDKey).(string); ok {
        return id
    }
    return ""
}

// SetRequestID adds request ID to context
func SetRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, RequestIDKey, id)
}

// SetUserID adds user ID to context
func SetUserID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, UserIDKey, id)
}