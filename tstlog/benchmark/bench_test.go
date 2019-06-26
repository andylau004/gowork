package benchmarks

import (
	// "io/ioutil"
	// "log"

	"fmt"
	"os"
	"testing"

	blg4go "github.com/YoungPioneers/blog4go"
	"github.com/cihub/seelog"

	"github.com/Sirupsen/logrus"
	"github.com/rs/zerolog"
	"go.uber.org/zap"

	"../../tstlog"

)

// fy@fy-ubuntu:~/WorkDir/mygitwork/gocode/tstlog-zap/benchmark$
// go test  -benchmem  -bench .  > 1

// -benchmem 可以提供每次操作分配内存的次数，以及每次操作分配的字节数。
// -benchtime 可以控制benchmark的运行时间 如果想让测试运行的时间更长，可以通过 -benchtime 指定，比如3秒

// go test bench_test.go -benchmem -bench="Benchmark_WriteLogFile*" > 1
// go test bench_test.go -benchmem -bench="Benchmark_Fun*" > 1
// go test bench_test.go -benchmem -bench="Benchmark_xxx*" > 1
// go test bench_2_test.go -benchmem -bench="Benchmark_file2*"

// go test bench_test.go -benchmem -bench="Benchmark_WriteLogFile*"

// 此函数不做压测使用，只加载下seelog, blg4go实例，防止import报错
func Benchmark_xxx_1(b *testing.B) {
	b.Run("xxx_1", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				fmt.Println("1111111111111111111111")
			}
		})
	})

	seelog.Info("Hello route usr:", "a")
	blg4go.Infof("Hello route usr:%s", "aaa")
	tstlog.GetMessage(tstlog.GenerateRangeNum(0, 1000))
}

func Benchmark_WriteLogFile_1(b *testing.B) {

	{ // 非json方式
		dls_zap_nojson := tstlog.NewZapLogger_nojson()
		defer dls_zap_nojson.Sync() // flushes buffer, if any
		b.Run("zap_fun--Benchmark_WriteLogFile--Text", func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					str := tstlog.GetMessage(tstlog.GenerateRangeNum(0, 1000))
					dls_zap_nojson.Infof("Hello route usr:%s", str)
				}
			})
		})
	}

	{ // 非json方式
		// tstlog.Init_blog4goFun()
		b.Run("blog4go--Benchmark_WriteLogFile--Text", func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					str := tstlog.GetMessage(tstlog.GenerateRangeNum(0, 1000))
					blg4go.Infof("Hello route usr:%s", str)
					// blg4go.Info("------------------------0 tttttttttttttttttttttt 0------------------------")
				}
			})
		})
	}

	{ // 非json方式
		// // tstlog.Init_blog4goFun()
		// b.Run("blog4go--aaaaaaaaaaaBenchmark_WriteLogFile--Text", func(b *testing.B) {
		// 	b.ResetTimer()
		// 	b.RunParallel(func(pb *testing.PB) {
		// 		for pb.Next() {
		// 			// str := tstlog.GetMessage(tstlog.GenerateRangeNum(0, 1000))
		// 			// blg4go.Infof("Hello route usr:%s", str)
		// 			// blg4go.Info("xxxxxxxxxxxxxxxxx----------------")
		// 			// blg4go.Info("ttttttttttttttttaaaa");
		// 			// blg4go.Info("twoeritowjrtjwrjtwe ");
		// 		}
		// 	})
		// })
	}

	{ // 非json方式
		// tstlog.Init_seelogFun()
		b.Run("seelog--Benchmark_WriteLogFile--Text", func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					str := tstlog.GetMessage(tstlog.GenerateRangeNum(0, 1000))
					seelog.Infof("Hello route usr:%s", str)
				}
			})
		})
	}

}

func Benchmark_Fun_1(b *testing.B) {
	b.Run("logrus_fun", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logrus.Infoln("this is a dummy log")
			}
		})
	})
	{ // 非json方式
		zlogger, _ := zap.NewDevelopment()
		defer zlogger.Sync() // flushes buffer, if any
		sugar_zap := zlogger.Sugar()
		b.Run("sugar_zap_fun--Text", func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					sugar_zap.Info("this is a dummy log")
				}
			})
		})
	}

	{ // json方式
		zlogger, _ := zap.NewProduction()
		defer zlogger.Sync() // flushes buffer, if any
		sugar_zap := zlogger.Sugar()
		b.Run("sugar_zap_fun--json", func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					sugar_zap.Info("this is a dummy log")
				}
			})
		})
	}

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	b.Run("zerolog_fun", func(b *testing.B) {
		// defer logger.Flush()
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				logger.Info().Msg("this is a dummy log")
			}
		})
	})

	// b.Run("log_fun", func(b *testing.B) {
	// 	// defer logger.Flush()
	// 	b.ResetTimer()
	// 	b.RunParallel(func(pb *testing.PB) {
	// 		for pb.Next() {
	// 			log.Println("this is a dummy log")
	// 		}
	// 	})
	// })

}

// func Benchmark_Fun(b *testing.B) {
// 	b.Logf("Logging at a disabled level without any structured context.")

// 	b.Run("seelog_fun", func(b *testing.B) {
// 		stream := &blackholeStream{}
// 		logger, err := see_log.LoggerFromWriterWithMinLevelAndFormat(stream, see_log.ErrorLvl, "%Time %Level %Msg")
// 		if err != nil {
// 			b.Fatal(err)
// 		}
// 		b.ResetTimer()
// 		b.RunParallel(func(pb *testing.PB) {
// 			defer logger.Flush()
// 			for pb.Next() {
// 				logger.Info(getMessage(0))
// 			}
// 		})
// 	})
// 	b.Run("Zap.Sugar", func(b *testing.B) {
// 		// logger := newZapLogger(zap.ErrorLevel).Sugar()
// 		dl, _ := zap.NewDevelopment()
// 		logger := dl.Sugar()
// 		defer logger.Sync()

// 		b.ResetTimer()
// 		b.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				logger.Info(getMessage(0))
// 			}
// 		})
// 	})
// 	b.Run("Zap.SugarFormatting", func(b *testing.B) {
// 		// logger := newZapLogger(zap.ErrorLevel).Sugar()
// 		dl, _ := zap.NewDevelopment()
// 		logger := dl.Sugar()
// 		defer logger.Sync()

// 		b.ResetTimer()
// 		b.RunParallel(func(pb *testing.PB) {
// 			for pb.Next() {
// 				logger.Infof("%v %v %v %s %v %v %v %v %v %s\n", fakeFmtArgs()...)
// 			}
// 		})
// 	})

// }
