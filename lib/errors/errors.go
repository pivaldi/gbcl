package liberrors

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var isTerm bool = isatty.IsTerminal(os.Stdout.Fd())
var quiet bool = os.Getenv("QUIET") != "" || IsTestMode()
var envDebugLevel = os.Getenv("DEBUG_LEVEL")

const (
	colorBlack = iota + 30
	colorRed
	colorGreen
	colorYellow
	colorBlue
	colorMagenta
	colorCyan
	colorWhite

	// colorBold     = 1
	// colorDarkGray = 90
)

func IsTestMode() bool {
	return flag.Lookup("test.v") != nil
}

// HandleError is a fansy log error handler
func HandleError(err error) (b bool) {
	if err != nil {
		// notice that we're using 1, so it will actually log the where
		// the error happened, 0 = this function, we don't want that.
		pc, filename, line, _ := runtime.Caller(1)

		log.Error().Stack().Err(err).
			Str("package", runtime.FuncForPC(pc).Name()).
			Str("file", fmt.Sprintf("%s:%d", filename, line)).
			Send()

		// log.Error().Err(err).Stack().Send()

		b = true
		panic(err)
	}

	return
}

// HandleErrorExit is a fansy log error handler and exit on error
func HandleErrorExit(err error) (b bool) {
	if err != nil {
		HandleError(err)
		os.Exit(1)
	}

	return
}

// HandleErrorPanic is a fansy log error handler and panic on error
func HandleErrorPanic(err error) (b bool) {
	if err != nil {
		HandleError(err)
		panic(err)
	}

	return
}

// colorize returns the string s wrapped in ANSI code c, unless disabled is true.
func colorize(s any, c int, disabled bool) string {
	if disabled {
		return fmt.Sprintf("%s", s)
	}

	return fmt.Sprintf("\x1b[%dm%v\x1b[0m", c, s)
}

// InitZerolog inits the Zerolog package
func InitLog() {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	if quiet {
		zerolog.SetGlobalLevel(zerolog.Disabled)
	} else {
		debugL := zerolog.DebugLevel
		if envDebugLevel != "" {
			n, err := strconv.Atoi(envDebugLevel)
			if err != nil {
				n = 0
				log.Debug().Msg("DEBUG_LEVEL is not a number. Set to 0…")
			}

			debugL = zerolog.Level(n)
		}

		zerolog.SetGlobalLevel(debugL)
	}

	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		NoColor:    !isTerm,
		TimeFormat: "2006/01/02 - 15:04:05.999",
		// See https://github.com/rs/zerolog/issues/157
		// and https://github.com/rs/zerolog/pull/416
		FormatMessage: func(i any) string {
			if i == nil {
				return ""
			}

			return "[" + colorize(i, colorBlue, !isTerm) + "]"
		},
	}).With().Timestamp().Caller().Logger()
}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
