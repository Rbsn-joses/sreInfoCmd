package customLogger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

func InitLogger(logLevel string) *logrus.Logger {
	// Níveis de log suportados
	supportedLevels := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,

		"panic": logrus.PanicLevel,
	}

	// Verifica se o nível de log fornecido é válido
	level, ok := supportedLevels[logLevel]
	if !ok {
		level = logrus.InfoLevel // Nível padrão caso inválido
		logrus.Warnf("Invalid log level: %s. Defaulting to INFO.", logLevel)
	}

	// Cria um novo logger
	logger := logrus.New()

	// Configura o formatador para incluir informações do chamador
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (string, string) {
			filename := filepath.Base(frame.File)
			return fmt.Sprintf(" %s:%d ", filename, frame.Line), ""
		},
	})

	// Configura o nível de log, saída e relatório de chamador
	logger.SetLevel(level)
	logger.SetOutput(os.Stdout)
	logger.SetReportCaller(true)

	return logger
}
