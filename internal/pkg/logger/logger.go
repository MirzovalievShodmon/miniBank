package logger

import (
	"github.com/rs/zerolog"
	"os"
)

func InitLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).Level(zerolog.InfoLevel)

	return logger
}

//logger.Info().
//	Str("name", "john").
//	Int("age", 22).
//	Bool("registered", true).
//	Msg("new signup!")
//
//logger.Info().
//	Str("name", "john").
//	Int("age", 22).
//	Bool("registered", true).
//	Send()
