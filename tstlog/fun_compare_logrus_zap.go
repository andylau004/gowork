package tstlog

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Sirupsen/logrus"
	// "github.com/apex/log"
	zap "go.uber.org/zap"

	"github.com/rs/zerolog"
)

type dummy struct {
	Foo string `json:"foo"`
	Bar string `json:"bar"`
}

func Tst_compare_1() {
	// dummyLog := dummy{
	// 	Foo: "foo",
	// 	Bar: "bar",
	// }

	var x int64 = 0
	for i := 0; i < 10000; i++ {
		t := time.Now()
		// logrus.WithField("Dummy", dummyLog).Infoln("this is a dummy log")
		logrus.Infoln("this is a dummy log")
		x += time.Since(t).Nanoseconds()
	}

	// 如果是NewProduction则生成日志格式是json
	// zlogger, _ := zap.NewProduction()
	// 如果是NewDevelopment则声称日志格式是我们常见普通输出,你可以单独测试这俩方式
	// json方式确实快些
	zlogger, _ := zap.NewDevelopment()
	sugar_zap := zlogger.Sugar()

	var y int64 = 0
	// var val int = 123
	for i := 0; i < 10000; i++ {
		t := time.Now()
		// sugar.Infow("this is a dummy log", "Dummy", dummyLog)
		sugar_zap.Info("this is a dummy log")
		// sugar_zap.Infof("this is a dummy log va=%v", val)
		y += time.Since(t).Nanoseconds()
	}

	var z int64 = 0
	for i := 0; i < 10000; i++ {
		t := time.Now()
		// dummyStr, _ := json.Marshal(dummyLog)
		// log.Printf("this is a dummy log: %s\n", string(dummyStr))
		// log.Info("")
		log.Println("this is a dummy log")
		z += time.Since(t).Nanoseconds()
	}

	var t_zero int64 = 0
	{
		// logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		for i := 0; i < 10000; i++ {
			t := time.Now()
			logger.Info().Msg("this is a dummy log")
			t_zero += time.Since(t).Nanoseconds()
		}
	}

	fmt.Println("=====================")
	fmt.Printf("Logrus: %5d ns per request \n", x/10000)
	fmt.Printf("Zap:    %5d ns per request \n", y/10000)
	fmt.Printf("StdLog: %5d ns per request \n", z/10000)
	fmt.Printf("zerolog: %5d ns per request \n", t_zero/10000)

}
