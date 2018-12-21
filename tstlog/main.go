package tstlog

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"

	blg4go "github.com/YoungPioneers/blog4go"
)

var (
	_messages = fakeMessages(1000)
)

// logpath 日志文件路径
// loglevel 日志级别
func Init_ZapLogger_json(logpath string, loglevel string) *zap.Logger {
	hook := lumberjack.Logger{
		Filename:   logpath, // 日志文件路径
		MaxSize:    1024,    // megabytes
		MaxBackups: 10,      // 最多保留3个备份
		MaxAge:     7,       //days
		Compress:   true,    // 是否压缩 disabled by default
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
	// encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// encoderConfig.EncodeTime = zapcore.EpochMillisTimeEncoder
	core := zapcore.NewCore(
		// zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		level,
	)

	logger := zap.New(core)
	// logger.Info("DefaultLogger init success")
	return logger
}

type Test struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func tst_444() {
	t := &Test{
		Name: "xiaoming",
		Age:  12,
	}
	data, err := json.Marshal(t)
	if err != nil {
		fmt.Println("marshal is failed,err: ", err)
	}

	// 历史记录日志名字为：all-2018-11-15T07-45-51.763.log，服务重新启动，日志会追加，不会删除
	logger := Init_ZapLogger_json("./out.log", "debug")
	for i := 0; i < 6; i++ {
		logger.Info(fmt.Sprint("test log ", i), zap.Int("line", 47))
		logger.Debug(fmt.Sprint("debug log ", i), zap.ByteString("level", data))
		logger.Info(fmt.Sprint("Info log ", i), zap.String("level", `{"a":"4","b":"5"}`))
		logger.Warn(fmt.Sprint("Info log ", i), zap.String("level", `{"a":"7","b":"8"}`))
		logger.Info("test log----------------------------")
		// logger.Infof("test log----------------------------%d", 123)
	}
}

func tst_2() {
	logger := Init_ZapLogger_json("out.log", "info")
	logger.Info("test log", zap.Int("line", 47))
	logger.Info("test log----------------------------")
	// logger.Info("key1=", "val1", ", key2=", 1234)
	// logger.Infow("failed to fetch URL: %s", "http://example.com")
}

func tst_1() {
	sugar := zap.NewExample().Sugar()
	defer sugar.Sync()
	sugar.Infow("failed to fetch URL",
		"url", "http://example.com",
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("failed to fetch URL: %s", "http://example.com")
	// output:
	// {"level":"info","msg":"failed to fetch URL","url":"http://example.com","attempt":3,"backoff":"1s"}
	// {"level":"info","msg":"failed to fetch URL: http://example.com"}

}

func tst_5() {
	// fmt.Printf("\n*** Using the Example logger\n\n")

	logger := zap.NewExample()
	// logger.Debug("This is a DEBUG message")
	// logger.Info("This is an INFO message")
	// logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))
	// logger.Warn("This is a WARN message")
	// logger.Error("This is an ERROR message")
	// // logger.Fatal("This is a FATAL message")  // would exit if uncommented
	// logger.DPanic("This is a DPANIC message")
	// //logger.Panic("This is a PANIC message")   // would exit if uncommented

	// fmt.Println()

	fmt.Printf("*** Using the Development logger\n\n")

	logger, _ = zap.NewDevelopment()
	logger.Debug("This is a DEBUG message")
	logger.Info("This is an INFO message")
	// logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))
	// logger.Warn("This is a WARN message")
	// logger.Error("This is an ERROR message")
	// logger.Fatal("This is a FATAL message")   // would exit if uncommented
	// logger.DPanic("This is a DPANIC message") // would exit if uncommented
	//logger.Panic("This is a PANIC message")    // would exit if uncommented

	// output:
	// 2018-12-11T19:55:55.645+0800	DEBUG	tstlog/main.go:128	This is a DEBUG message
	// 2018-12-11T19:55:55.645+0800	INFO	tstlog/main.go:129	This is an INFO message

	fmt.Println()

	{
		// fmt.Printf("*** Using the Production logger\n\n")

		// logger, _ = zap.NewProduction()
		// logger.Debug("This is a DEBUG message")
		// logger.Info("This is an INFO message")
		// logger.Info("This is an INFO message with fields", zap.String("region", "us-west"), zap.Int("id", 2))
		// logger.Warn("This is a WARN message")
		// logger.Error("This is an ERROR message")
		// // logger.Fatal("This is a FATAL message")   // would exit if uncommented
		// logger.DPanic("This is a DPANIC message")
		// // logger.Panic("This is a PANIC message")   // would exit if uncommented
	}

}

func fakeMessages(n int) []string {
	messages := make([]string, n)
	for i := range messages {
		messages[i] = fmt.Sprintf("Test logging, but use a somewhat realistic message length. (#%v)", i)
	}
	return messages
}

func GetMessage(iter int) string {
	return _messages[iter%1000]
}

func tst_6() {
	fmt.Println("GetMessage(1)=", GetMessage(999))

	url := "www.baidu.com"
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("failed to fetch URL",
		// Structured context as loosely typed key-value pairs.
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("Failed to fetch URL: %s", url)

}

var gopherType string

// GenerateRangeNum 生成一个区间范围的随机数
func GenerateRangeNum(min, max int) int {
	randNum := rand.Intn(max - min)
	randNum = randNum + min
	// if randNum > 900 {
	// 	fmt.Printf("%v ", randNum)
	// }
	return randNum
}

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func tst_tmpfun() {

	for i := 0; i < 1000; i++ {
		// fmt.Println(" --------------------------beg")
		// defer fmt.Println(" --------------------------end")
		// fmt.Println(time.Now().Unix())                //获取当前秒
		// fmt.Println(time.Now().UnixNano())            //获取当前纳秒
		a := makeTimestamp()
		fmt.Println("a= ", a)
		// fmt.Println("ms=", time.Now().UnixNano()/1e6) //将纳秒转换为毫秒
		// fmt.Println(time.Now().UnixNano() / 1e9)      //将纳秒转换为秒

		// c := time.Unix(time.Now().UnixNano()/1e9, 0) //将秒转换为 time 类型
		// fmt.Println(c.String())                      //输出当前英文时间戳格式

		time.Sleep(time.Second)
	}

}

func HelloBlg() {
	str := "aaaaaaaaaaaaaaaa"
	blg4go.Infof("Hello route usr:%s", str)

	var exit_chan chan int
	// HelloBlg()
	<-exit_chan

}

func main() {
	fmt.Println("main in")
	defer fmt.Println("main out")

	// tst_tmpfun()
	// return
	// Tst_timeFormat()

	// t := time.Now()
	// fmt.Println(t.Format("2006-01-02 15:04:05"))

	// return

	rand.Seed(time.Now().Unix())

	// for i := 0; i < 100; i++ {
	// 	GenerateRangeNum(0, 1000)
	// 	if i > 0 && i%32 == 0 {
	// 		fmt.Println("")
	// 	}
	// }
	// return

	// tst_sugar_printf_console()
	// return

	// tst_5()
	// return

	// time.Sleep(time.Second)
	// InitBlog4go("blg.log", 1, 7)
	// blg4go.Error("something1---------------------------")

	// blg4go.Info("something1---------------------------")
	// return

	// fmt.Println(fmt.Sprint("debug log ", 1))
	// return
	// // onename := "aaaa"
	// // fmt.Println(fmt.Sprintf("%s", onename))
	// // return
	// tst_2()
	// return

	namePtr := flag.String("l", "l", "日志类型: zap使用zap方式打印;std使用标准方式打印;")
	flag.Parse()

	args := flag.Args()
	fmt.Println("name:", *namePtr)
	fmt.Println("args:", args)
	fmt.Println("os.Args=", os.Args)

	// logrus_logger.Info("Hello route usr:", "aaaaaaaaaaaaaaa")
	// return

	StartHttpSrv_zap()
	return

	if *namePtr == "zap" {
		fmt.Println("startup zap server")
		StartHttpSrv_zap()
	} else if *namePtr == "std" {
		fmt.Println("startup std server")
		StartHttpSrv_std()
	}

	return

	Tst_zap_fun_1()
	return

	Tst_compare_1()
	return
	// Tst_zerolog_1()

	//zap.NewDevelopment() 包含代码中文件信息
	//zap.NewProduction() 去除了文件信息

	// tst_newproduction()
	Tst_zap_fun_1()
	return

	Tst_zerolog_entry()
	return

	// tst_sugar()
}
