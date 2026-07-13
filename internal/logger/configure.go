package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func Configure(cmd *cobra.Command) {
	colorOutput, _ := cmd.Flags().GetBool("color")

	logLevel := getLogLevel()
	zerolog.SetGlobalLevel(logLevel)

	if colorOutput {
		out := zerolog.ConsoleWriter{Out: os.Stdout}
		log.Logger = log.Output(out)
	}
}

func getLogLevel() zerolog.Level {
	logLevelStr := strings.TrimSpace(os.Getenv("LOG_LEVEL"))
	if len(logLevelStr) == 0 {
		return zerolog.InfoLevel
	}
	switch strings.ToUpper(logLevelStr) {
	case "TRACE":
		return zerolog.TraceLevel
	case "DEBUG":
		return zerolog.DebugLevel
	case "INFO":
		return zerolog.InfoLevel
	case "ERROR":
		return zerolog.ErrorLevel
	case "WARN":
		return zerolog.WarnLevel
	case "FATAL":
		return zerolog.FatalLevel
	default:
		return zerolog.InfoLevel
	}
}
