package tstlog

import (
	"log"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// 此服务不打印日志
var (
// _globalMu sync.RWMutex
// zap_L, _  = NewZapConsolLogger()
// zap_sugar = zap_L.Sugar()
)

// zlogger, _ := zap.NewProduction()
// defer zlogger.Sync() // flushes buffer, if any

func std_Index(ctx *fasthttp.RequestCtx) {
	// log.Fprint(ctx, "Welcome!\n")
	// log.Println("Welcome!")
}
func std_Hello(ctx *fasthttp.RequestCtx) {
	// fmt.Fprintf(ctx, "hello, %s!\n", ctx.UserValue("name"))
	// log.Printf("hello, %s!", ctx.UserValue("name"))
}

func StartHttpSrv_std() {
	router := fasthttprouter.New()

	router.GET("/", std_Index)
	router.GET("/hello/:name", std_Hello)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))

}
