package main

import (
	"fmt"

	blg4go "github.com/YoungPioneers/blog4go"
	"github.com/valyala/fasthttp"

	"Utilgo"
)

var g_strListenAddr string

func init() {
}

func ServeHTTP(ctx *fasthttp.RequestCtx) {
	fmt.Printf("ctx=%+v\n", *ctx)
	blg4go.Infof("ctx=%+v\n", *ctx)
	oneBody := Utilgo.Bytes2str(ctx.PostBody())
	blg4go.Info("oneBody=", oneBody)

	switch string(ctx.Path()) {
	case "/":
		break
	default:
		ctx.Response.SetBody([]byte("Error:2001"))
		blg4go.Error("main|ServeHTTP|path|", string(ctx.Path()), " is not support!")
	}

}

func StartRecvUpload() {
	g_strListenAddr = ":9089"

	defer fmt.Println("listen server done!")

	h := ServeHTTP
	if err := fasthttp.ListenAndServe(g_strListenAddr, h); err != nil {
		fmt.Println("listen failed! err=", err)
		panic(err)
	}

}
