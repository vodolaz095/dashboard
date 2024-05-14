package zerologger

import (
	"io"
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/journald"
	"github.com/rs/zerolog/log"
	"github.com/vodolaz095/dashboard/config"
)

// ExtractZerologLevel получает уровень логгирования в совместимом с zerolog формате
func ExtractZerologLevel(level string) zerolog.Level {
	switch level {
	case TraceLevel:
		return zerolog.TraceLevel
	case DebugLevel:
		return zerolog.DebugLevel
	case InfoLevel:
		return zerolog.InfoLevel
	case WarnLevel:
		return zerolog.WarnLevel
	case ErrorLevel:
		return zerolog.ErrorLevel
	case FatalLevel:
		return zerolog.FatalLevel
	default:
		return zerolog.DebugLevel // TODO - может быть info???
	}
}

func Configure(params config.Log) {
	var outputsEnabled []io.Writer

	if params.ToJournald {
		outputsEnabled = append(outputsEnabled, journald.NewJournalDWriter())
	} else {
		outputsEnabled = append(outputsEnabled, zerolog.ConsoleWriter{
			Out:        os.Stdout, // https://12factor.net/ru/logs
			TimeFormat: "15:04:05",
		})
	}
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return file + ":" + strconv.Itoa(line)
	}
	sink := zerolog.New(zerolog.MultiLevelWriter(outputsEnabled...)).
		With().Timestamp().Caller().
		Logger().Level(ExtractZerologLevel(params.Level))
	log.Logger = sink
	return
}
