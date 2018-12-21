package tstlog

import (
	// "os"

	"os"

	"github.com/rs/zerolog"
	z_log "github.com/rs/zerolog"
)

func newZerolog() zerolog.Logger {
	return zerolog.New(os.Stdout).With().Timestamp().Logger()
}
func Tst_zerolog_2() {
	logger := newZerolog()
	logger.Info().Msgf("%v %v %s\n", 1, 2, 3, "12345")
	//output:{"level":"info","time":"2018-12-10T15:06:02+08:00","message":"1 2 %!s(int=3)\n%!(EXTRA string=12345)"}

}

// 仅支持json格式日志,排除掉
func Tst_zerolog_1() {
	// zlog.Info().Msg("hello world")
	logger := z_log.New(os.Stderr).With().Timestamp().Logger()

	logger.Info().Str("foo", "bar").Msg("hello world")
	// output:{"level":"info","foo":"bar","time":"2018-12-10T13:40:54+08:00","message":"hello world"}
}

func Tst_zerolog_entry() {
	Tst_zerolog_2()
}
