package log

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"refina-web-bff/config/env"
	"refina-web-bff/internal/utils/data"

	"github.com/sirupsen/logrus"
)

// ApacheStyleFormatter — custom formatter with color support
type ApacheStyleFormatter struct {
	NoColors bool
}

func (f *ApacheStyleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	var levelColor string
	resetColor := "\x1b[0m"

	if !f.NoColors {
		switch entry.Level {
		case logrus.DebugLevel:
			levelColor = "\x1b[36m"
		case logrus.InfoLevel:
			levelColor = "\x1b[32m"
		case logrus.WarnLevel:
			levelColor = "\x1b[33m"
		case logrus.ErrorLevel:
			levelColor = "\x1b[31m"
		case logrus.FatalLevel, logrus.PanicLevel:
			levelColor = "\x1b[35m"
		default:
			levelColor = "\x1b[0m"
		}
	}

	timestamp := entry.Time.Format("02/Jan/2006:15:04:05 -0700")
	level := strings.ToUpper(entry.Level.String())

	if !f.NoColors {
		fmt.Fprintf(b, "[%s] %s%s%s: %s", timestamp, levelColor, level, resetColor, entry.Message)
	} else {
		fmt.Fprintf(b, "[%s] %s: %s", timestamp, level, entry.Message)
	}

	if len(entry.Data) > 0 {
		fmt.Fprintf(b, " - ")
		fieldCount := 0
		for key, value := range entry.Data {
			if fieldCount > 0 {
				fmt.Fprintf(b, ", ")
			}
			var valueStr string
			switch v := value.(type) {
			case string:
				if strings.ContainsAny(v, " ,=") {
					valueStr = fmt.Sprintf(`"%s"`, v)
				} else {
					valueStr = v
				}
			default:
				valueStr = fmt.Sprintf("%v", v)
			}
			fmt.Fprintf(b, "%s: %s", key, valueStr)
			fieldCount++
		}
	}
	b.WriteByte('\n')
	return b.Bytes(), nil
}

var logger *logrus.Logger

func SetupLogger() {
	logger = logrus.New()
	logger.SetOutput(os.Stdout)

	isProduction := env.Cfg.Server.Mode != data.DEVELOPMENT_MODE
	logger.SetFormatter(&ApacheStyleFormatter{NoColors: isProduction})

	if isProduction {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(logrus.DebugLevel)
	}
}

func getLogger() *logrus.Logger {
	if logger == nil {
		SetupLogger()
	}
	return logger
}

func Info(msg string, fields map[string]any) {
	getLogger().WithFields(fields).Info(msg)
}

func Warn(msg string, fields map[string]any) {
	getLogger().WithFields(fields).Warn(msg)
}

func Error(msg string, fields map[string]any) {
	getLogger().WithFields(fields).Error(msg)
}

func Fatal(msg string, fields map[string]any) {
	getLogger().WithFields(fields).Fatal(msg)
}

func Debug(msg string, fields map[string]any) {
	getLogger().WithFields(fields).Debug(msg)
}

func WithRequestID(requestID string) *logrus.Entry {
	return getLogger().WithField("request_id", requestID)
}

func GetInstance() *logrus.Logger {
	return getLogger()
}

func Ms(d time.Duration) float64 {
	return float64(d.Nanoseconds()) / 1e6
}
