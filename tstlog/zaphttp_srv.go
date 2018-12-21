package tstlog

import (
	"fmt"
	"log"
	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	zap "go.uber.org/zap"

	blg4go "github.com/YoungPioneers/blog4go"

	seelog "github.com/cihub/seelog"

	"github.com/Sirupsen/logrus"
)

const (
	LEVEL_TRACE = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARNING
	LEVEL_ERROR
	LEVEL_CRITICLE
	LEVEL_UNKNOWN
)

// 测试命令
// wrk -t4 -c1000 -d30s -T30s --latency http://192.168.166.41:8081/hello/a111111111111111111

var (
	// 打屏幕
	// zap_L, _  = NewZapConsolLogger()
	// zap_sugar = zap_L.Sugar()

	// zap 非json方式 打屏幕
	// zap_nojson, _  = zap.NewDevelopment()
	// dls_zap_nojson = zap_nojson.Sugar()

	// zap 文本方式 写文件
	// Requests/sec:  27363.89
	// Transfer/sec:      2.43MB
	dls_zap_nojson = NewZapLogger_nojson()

	// zap json方式 写日志
	// Requests/sec:  27913.92
	// Transfer/sec:      2.48MB
	zap_logger_json = Init_ZapLogger_json("./logs/zap.log", "info")

	// blog4go 文本方式 写日志
	// Requests/sec:  38255.16
	// Transfer/sec:      3.39M

	// seelog 文本方式 写日志
	// Requests/sec:  18735.70
	// Transfer/sec:      1.66MB

	// logrus json方式 写日志
	// Requests/sec:  15866.41
	// Transfer/sec:      1.41MB
	logrus_logger = logrus.New()
)

func init() {
	Init_blog4goFun()
	Init_seelogFun()
	init_logrusFun()
}

func init_zapFun() {
}

func init_logrusFun() {
	file, err := os.OpenFile("./logs/logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logrus_logger.Out = file
	} else {
		fmt.Println("Failed to log to file, using default stderr")
		fmt.Println("init_logrusFun failed!!!")
		os.Exit(1)
		return
	}
	// logrus_logger.Info("A group of walrus emerges from the ocean")
}

func Init_seelogFun() {
	see_logger, err := seelog.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		fmt.Println("Init_seelogFun failed!!!")
		seelog.Critical("err parsing config log file", err)
		os.Exit(1)
		return
	}
	seelog.ReplaceLogger(see_logger)
}

func Init_blog4goFun() {
	// fmt.Println("Init_blog4goFun beg ")
	err := blg4go.NewWriterFromConfigAsFile("blog4go_config.xml")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
		fmt.Println("Init_blog4goFun failed!!!")
	}
	// defer blg4go.Close()
	// fmt.Println("Init_blog4goFun end ")
}

// zlogger, _ := zap.NewProduction()
// defer zlogger.Sync() // flushes buffer, if any

func NewZapConsolLogger() (*zap.Logger, error) {
	return zap.NewProduction()
	// return &Logger{
	// 	core:        zapcore.NewNopCore(),
	// 	errorOutput: zapcore.AddSync(ioutil.Discard),
	// 	addStack:    zapcore.FatalLevel + 1,
	// }
}

func Index(ctx *fasthttp.RequestCtx) {
	// fmt.Fprint(ctx, "Welcome!\n")
	// zap_sugar.Info("Welcome Index")
	// zap_logger_json.Info("Welcome Index")
	// blg4go.Info("Welcome Index")
}
func Hello(ctx *fasthttp.RequestCtx) {
	// str := fmt.Sprintf("%v", ctx.UserValue("name"))

	str := GetMessage(GenerateRangeNum(0, 1000))
	// fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))

	// zap_sugar.Infof("hello! usr:%s", str)
	// dls_zap_nojson.Infof("Hello route usr:%s", str)
	zap_logger_json.Info("Hello route", zap.String("usr", str))

	// blg4go.Infof("Hello route usr:%s", str)

	// seelog.Info("Hello route usr:", str)

	// logrus_logger.Info("Hello route usr:", str)
}

func StartHttpSrv_zap() {
	router := fasthttprouter.New()

	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	log.Fatal(fasthttp.ListenAndServe(":8081", router.Handler))
}

func InitBlog4go(filename string, level int, maxdays int) {
	// level = LEVEL_DEBUG
	// if !path.IsAbs(filename) {
	// 	workingdir, _ := gofile.WorkDir()
	// 	if workingdir != "" {
	// 		filename = path.Join(workingdir, filename)
	// 	}
	// }
	// diroflog := path.Dir(filename)
	// os.MkdirAll(diroflog, 0777)

	// blog4go.NewBaseFileWriter(filename, true)
	// //兼容旧库loglevel 7
	// if level > 5 {
	// 	level = 1
	// }
	// fmt.Println("filename=", filename)
	// fmt.Println("diroflog=", diroflog)
	// blog4go.SetLevel(blog4go.Levels[level])

}
