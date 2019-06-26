
package bench

import (
	// "encoding/json"
	// "io/ioutil"
	// "log"
	// "os"
	// "time"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
    "gopkg.in/natefinch/lumberjack.v2"

	"go.uber.org/zap/internal/ztest"
	
)

// logpath 日志文件路径
// loglevel 日志级别
func initLogger(logpath string, loglevel string) *zap.Logger {

    hook := lumberjack.Logger{
        Filename:   logpath, // 日志文件路径
        MaxSize:    1024, // megabytes
        MaxBackups: 3,    // 最多保留3个备份
        MaxAge:     7,    //days
        Compress:   true, // 是否压缩 disabled by default
    }
    w := zapcore.AddSync(&hook)

	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
    var level zapcore.Level
    switch loglevel {
    case "debug":
        level = zap.DebugLevel
    case "info":
        level = zap.InfoLevel
    case "error":
        level = zap.ErrorLevel
    default:
        level = zap.InfoLevel
    }
	encoderConfig := zap.NewProductionEncoderConfig()
	// 时间格式
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    core := zapcore.NewCore(
        zapcore.NewConsoleEncoder(encoderConfig),
        w,
        level,
    )

    logger := zap.New(core)
    logger.Info("DefaultLogger init success")

    return logger
}

func newZapLogger(lvl zapcore.Level) *zap.Logger {
	ec := zap.NewProductionEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewJSONEncoder(ec)
	return zap.New(zapcore.NewCore(
		enc,
		&ztest.Discarder{},
		lvl,
	))
}
func newSampledLogger(lvl zapcore.Level) *zap.Logger {
	return zap.New(zapcore.NewSampler(
		newZapLogger(zap.DebugLevel).Core(),
		100*time.Millisecond,
		10, // first
		10, // thereafter
	))
}


func BenchmarkZapTextPositive(b *testing.B) {
	// stream := &blackholeStream{}

	
	logger, _ := zap.NewDevelopment()
	b.ResetTimer()

	// var i_count int = 0

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("The quick brown fox jumps over the lazy dog")

			// logger.Info("This is an INFO message with fields", 
            // 			zap.String("region", "us-west"), 
            // 			zap.Int("id", 2))

// 			slogger := logger.Sugar()
// slogger.Info("Info() uses sprint")
// slogger.Infof("Infof() uses %s", "sprintf")
// // slogger.Infow("Infow() allows tags", "name", "Legolas", "type", 1)
// output: 2018-12-08T13:01:28.746+0800	INFO	go-loggers-bench/zap_test.go:77	Info() uses sprint
//         2018-12-08T13:01:28.746+0800	INFO	go-loggers-bench/zap_test.go:78	Infof() uses sprintf

			//  newlogger, _ := zap.NewProduction()
			//  newlogger.Info( "aaaaaaaaaaaaaaaaaaa" )
			// output:{"level":"info","ts":1544244850.8146913,"caller":"go-loggers-bench/zap_test.go:82","msg":"aaaaaaaaaaaaaaaaaaa"}

			// i_count++
			// if i_count == 1 {
			// 	break
			// }

		}
	})

	// if stream.WriteCount() != uint64(b.N) {
	// 	b.Fatalf("Log write count")
	// }


    // logger := initLogger("./out.log", "debug")
    // for i := 0; i < 6; i++ {
    //     logger.Info(fmt.Sprint("test log ", i), zap.Int("line", 47))
    //     logger.Debug(fmt.Sprint("debug log ", i), zap.ByteString("level", data))
    //     logger.Info(fmt.Sprint("Info log ", i), zap.String("level", `{"a":"4","b":"5"}`))
    //     logger.Warn(fmt.Sprint("Info log ", i), zap.String("level", `{"a":"7","b":"8"}`))
	// 	logger.Info("test log----------------------------")
    // }

	// loglevel := "debug"
	
	// w := zapcore.AddSync(&stream)

	// // 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// // debug->info->warn->error
    // var level zapcore.Level
    // switch loglevel {
    // case "debug":
    //     level = zap.DebugLevel
    // case "info":
    //     level = zap.InfoLevel
    // case "error":
    //     level = zap.ErrorLevel
    // default:
    //     level = zap.InfoLevel
    // }
	// encoderConfig := zap.NewProductionEncoderConfig()
	// // 时间格式
    // encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
    // core := zapcore.NewCore(
    //     zapcore.NewConsoleEncoder(encoderConfig),
    //     w,
    //     level,
    // )

    // logger := zap.New(core)
    // logger.Info("DefaultLogger init success")


	// if stream.WriteCount() != uint64(b.N) {
	// 	b.Fatalf("Log write count")
	// }
}


// func testEncoderConfig() EncoderConfig {
// 	return EncoderConfig{
// 		MessageKey:     "msg",
// 		LevelKey:       "level",
// 		NameKey:        "name",
// 		TimeKey:        "ts",
// 		CallerKey:      "caller",
// 		StacktraceKey:  "stacktrace",
// 		LineEnding:     "\n",
// 		EncodeTime:     EpochTimeEncoder,
// 		EncodeLevel:    LowercaseLevelEncoder,
// 		EncodeDuration: SecondsDurationEncoder,
// 		EncodeCaller:   ShortCallerEncoder,
// 	}
// }

// func BenchmarkJSONLogMarshalerFunc(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		enc := NewJSONEncoder(testEncoderConfig())
// 		enc.AddObject("nested", ObjectMarshalerFunc(func(enc ObjectEncoder) error {
// 			enc.AddInt64("i", int64(i))
// 			return nil
// 		}))
// 	}
// }


