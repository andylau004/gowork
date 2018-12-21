package tstlog

import (
	"fmt"
	"os"
	"time"

	zap "go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// A Syncer is a spy for the Sync portion of zapcore.WriteSyncer.
type Syncer struct {
	err    error
	called bool
}

// SetError sets the error that the Sync method will return.
func (s *Syncer) SetError(err error) {
	s.err = err
}

// Sync records that it was called, then returns the user-supplied error (if
// any).
func (s *Syncer) Sync() error {
	s.called = true
	return s.err
}

// Called reports whether the Sync method was called.
func (s *Syncer) Called() bool {
	return s.called
}

// A Discarder sends all writes to ioutil.Discard.
type Discarder struct{ Syncer }

// Write implements io.Writer.
func (d *Discarder) Write(b []byte) (int, error) {
	// return ioutil.Discard.Write(b)
	return os.Stderr.Write(b)

}

//zap.NewDevelopment() 包含代码中文件信息
//zap.NewProduction() 去除了文件信息
// 输出json格式日志,排除掉此库
func tst_newproduction() {
	fmt.Println("tst_newproduction beg")
	defer fmt.Println("tst_newproduction end\n\n")

	var i8 int8 = 10
	var str = "string"
	any := struct {
		I int `json:"int"`
		S string
	}{
		I: 1,
		S: "str",
	}

	pl, _ := zap.NewProduction()

	pl.With(zap.Namespace("namespace")).Named("name").Warn("NewProduction name", zap.Any("any", any))
	//pl.Fatal("NewProduction")
	go func() {
		defer func() {
			if err := recover(); err != nil {
				// log.Println(err)
			}
		}()
		pl.Panic("NewProduction", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	}()
	pl.DPanic("NewProduction", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	pl.Error("NewProduction", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	pl.Warn("NewProduction", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	pl.Info("NewProduction", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	pl.With(zap.Int8("i8", i8)).Info("NewProduction", zap.Any("any", any), zap.String("str", str))
	pl.Info("NewProduction", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str), zap.Namespace("namespace"))
	pl.Info("NewProduction", zap.Namespace("namespace"), zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	pl.Debug("NewProduction", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
}

//zap.NewDevelopment() 包含代码中文件信息
func tst_newdevelopment() {
	fmt.Println("tst_newdevelopment beg-------------------------------------------")
	defer fmt.Println("tst_newdevelopment end-------------------------------------------\n\n")

	var i8 int8 = 10
	var str = "string"
	any := struct {
		I int `json:"int"`
		S string
	}{
		I: 1,
		S: "str",
	}

	dl, _ := zap.NewDevelopment()
	//dl.Fatal("NewDevelopment")
	// go func() {
	//     defer func() {
	//         if err := recover(); err != nil {
	//             // log.Println(err)
	//         }
	//     }()
	//     dl.Panic("NewDevelopment", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	// }()
	// go func() {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			// log.Println(err)
	// 		}
	// 	}()
	// 	dl.DPanic("NewDevelopment", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	// 	dl.Info("aaaaaaaaaaaaaaaa")
	// }()
	dl.Error("NewDevelopment", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	dl.Warn("NewDevelopment", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	dl.Info("NewDevelopment", zap.Int8("i8", i8), zap.Any("any", any), zap.String("str", str))
	dl.Debug("NewDevelopment", zap.Int8("Debug_i8", i8), zap.Any("Debug_any", any), zap.String("str", str))

	// abc := "bbbbbbbbbbbbbbbbbbbbb"
	// dl.Info("key=", abc)
	// dl.Debugf("Sugar i8=%d str=%s", i8, str)
	// dl.Infof("Sugar i8=%d str=%s", i8, str)
	// dl.Errorf("Sugar i8=%d str=%s", i8, str)

	time.Sleep(time.Second)
}

func newZapLogger(lvl zapcore.Level) *zap.Logger {
	// ec := zap.NewProductionEncoderConfig()
	ec := zap.NewDevelopmentEncoderConfig()
	ec.EncodeDuration = zapcore.NanosDurationEncoder
	ec.EncodeTime = zapcore.EpochNanosTimeEncoder
	enc := zapcore.NewJSONEncoder(ec)
	return zap.New(zapcore.NewCore(
		enc,
		&Discarder{},
		// ioutil.Discard,
		lvl,
	))
}

func tst_example_1() {
	// Using zap's preset constructors is the simplest way to get a feel for the
	// package, but they don't allow much customization.
	logger := zap.NewExample() // or NewProduction, or NewDevelopment
	defer logger.Sync()

	const url = "http://example.com"

	// In most circumstances, use the SugaredLogger. It's 4-10x faster than most
	// other structured logging packages and has a familiar, loosely-typed API.
	sugar := logger.Sugar()
	sugar.Infow("Failed to fetch URL.",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)
	// output:
	// {"level":"info","msg":"Failed to fetch URL.","url":"http://example.com","attempt":3,"backoff":"1s"}
	// {"level":"info","msg":"Failed to fetch URL: http://example.com"}

}

func tst_sugar_printf_console() {
	fmt.Println("tst_sugar_printf_console beg-------------------------------------------")
	defer fmt.Println("tst_sugar_printf_console end-------------------------------------------\n\n")

	var i8 int8 = 10
	var str = "string"

	dl, _ := zap.NewDevelopment()
	dls := dl.Sugar()

	// go func() {
	// 	defer func() {
	// 		if err := recover(); err != nil {
	// 			log.Println("recover err=", err)
	// 		}
	// 	}()
	// 	dls.Panic("Sugar 1111111111111err=", i8, ", 1111111111111errstr=", str)
	// }()

	abc := "bbbbbbbbbbbbbbbbbbbbb"
	dls.Info("Starting Server")

	dls.Info("key=", abc, ", xxx=", abc)
	dls.Debugf("Sugar i8=%d str=%s", i8, str)
	dls.Infof("Sugar i8=%d str=%s", i8, str)
	dls.Errorf("Sugar i8=%d str=%s", i8, str)

	time.Sleep(time.Second)

	// go func() {
	//     defer func() {
	//         if err := recover(); err != nil {
	//             // log.Println(err)
	//         }
	//     }()
	//     dls.DPanic("Sugar NewDevelopment", "i8", i8, "any", any, "str", str)
	// }()
	// dls.Error("Sugar NewDevelopment", "i8", i8, "any", any, "str", str, "end")
	// dls.Warn("Sugar NewDevelopment", "i8", i8, "any", any, "str", str, "end")
	// dls.Info("Sugar NewDevelopment", "i8", i8, "any", any, "str", str, "end")
	// dls.Debug("Sugar NewDevelopment", "i8", i8, "any", any, "str", str, "end")
}

func tst_dummy_sugar() {
	// logger := newZapLogger(zap.ErrorLevel).Sugar()
	logger := newZapLogger(zap.InfoLevel).Sugar()
	defer logger.Sync()

	var i_index int = 0
	for i_index < 1000 {
		logger.Info(GetMessage(0))
		i_index++
	}
	//output:
	//{"L":"INFO","T":1544439938604779280,"M":"Test logging, but use a somewhat realistic message length. (#0)"}
}

func Tst_zap_fun_1() {
	{
		// zlogger, _ := zap.NewProduction()
		// sugar_zap := zlogger.Sugar()

		// var y int64 = 0
		// for i := 0; i < 10000; i++ {
		// 	t := time.Now()
		// 	// sugar.Infow("this is a dummy log", "Dummy", dummyLog)
		// 	sugar_zap.Info("this is a dummy log val----------------")
		// 	// sugar_zap.Infof("this is a dummy log val=%v", y)
		// 	y += time.Since(t).Nanoseconds()
		// }
	}
	{
		logger, _ := zap.NewProduction()
		// logger.Debug("This is a DEBUG message")
		logger.Info("NewProduction    This is an INFO message")
	}
	fmt.Println("")
	{
		logger, _ := zap.NewDevelopment()
		// logger.Debug("This is a DEBUG message")
		logger.Info("NewDevelopment    This is an INFO message")
	}

	return

	// tst_example_1()
	// return

	// tst_dummy_sugar()
	// return

	// tst_7()
	// tst_newdevelopment()
	tst_sugar_printf_console()
}

func tst_7() {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()

	sugar.Infow("failed to fetch URL",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("failed to fetch URL: %s", "http://example.com")
	sugar.Info("failed to fetch URL:", "http://example.com")

	// output:
	//{"level":"info","msg":"failed to fetch URL","url":"http://example.com","attempt":3,"backoff":"1s"}
	//{"level":"info","msg":"failed to fetch URL: http://example.com"}
	//{"level":"info","msg":"failed to fetch URL:http://example.com"}

	{
		// logger := zap.NewExample()
		// logger.Debug("This is a DEBUG message")
		// logger.Info("This is an INFO message")
		// // logger.Infof("This is an INFO message, val=%d", 123)
		// // output:{"level":"debug","msg":"This is a DEBUG message"}
		// //        {"level":"info","msg":"This is an INFO message"}
	}

	{
		logger, _ := zap.NewDevelopment()
		// logger.Debug("This is a DEBUG message")
		logger.Info("NewDevelopment    This is an INFO message")
		//output:2018-12-10T14:51:28.708+0800	INFO	tstlog-zap/fun_zap.go:162	NewDevelopment    This is an INFO message
	}
}
