package logger

import (
	"io"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	once sync.Once
)

func (l *logConfig) newZeroLogger() zerolog.Logger {

	var log zerolog.Logger

	once.Do(func() {

		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		logLevel := zerolog.InfoLevel

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
			FieldsExclude: []string{
				"user_agent",
				"git_revision",
				"go_version",
			},
		}

		fileLogger := &lumberjack.Logger{
			Filename:   l.LogFile,
			MaxSize:    5, //
			MaxBackups: 10,
			MaxAge:     14,
			Compress:   true,
		}

		output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)

		var gitRevision string

		buildInfo, ok := debug.ReadBuildInfo()
		if ok {
			for _, v := range buildInfo.Settings {
				if v.Key == "vcs.revision" {
					gitRevision = v.Value
					break
				}
			}
		}

		log = zerolog.New(output).
			Level(zerolog.Level(logLevel)).
			With().
			Timestamp().
			Str("git_revision", gitRevision).
			Str("go_version", buildInfo.Main.Version).
			Logger()

		zerolog.DefaultContextLogger = &log
	})

	return log
}

// func (l *logConfig) newZeroLogger() *zerolog.Logger {
// 	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
// 	zerolog.TimeFieldFormat = time.RFC3339Nano

// 	logLevel := zerolog.InfoLevel

// 	var output io.Writer = zerolog.ConsoleWriter{
// 		Out:        os.Stdout,
// 		TimeFormat: time.RFC3339,
// 		FieldsExclude: []string{
// 			"user_agent",
// 			"git_revision",
// 			"go_version",
// 		},
// 	}

// 	if l.Environment != "development" {
// 		fileLogger := &lumberjack.Logger{
// 			Filename:   null.ValueFromNull[string](l.OutputFile),
// 			MaxSize:    5, //
// 			MaxBackups: 10,
// 			MaxAge:     14,
// 			Compress:   true,
// 		}

// 		output = zerolog.MultiLevelWriter(os.Stderr, fileLogger)
// 	}

// 	var gitRevision string

// 	buildInfo, ok := debug.ReadBuildInfo()
// 	if ok {
// 		for _, v := range buildInfo.Settings {
// 			if v.Key == "vcs.revision" {
// 				gitRevision = v.Value
// 				break
// 			}
// 		}
// 	}

// 	fmt.Println()

// 	log := zerolog.New(output).
// 		Level(zerolog.Level(logLevel)).
// 		With().
// 		Timestamp().
// 		Str("git_revision", gitRevision).
// 		Str("go_version", buildInfo.Main.Version).
// 		Logger()

// 	zerolog.DefaultContextLogger = &log

// 	return &log
// }
