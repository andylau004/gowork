

package main

import (
	// "context"
	"fmt"
	"time"

	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	blg4go "github.com/YoungPioneers/blog4go"
)

func GetCurrentPath() (string, error) {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return "", err
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return "", errors.New(`error: Can't find "/" or "\".`)
	}
	return string(path[0 : i+1]), nil
}

func dummyfun1() {

	time.Sleep( time.Second )
	fmt.Println( "" )

}

// add by fy 2019-5-21
// 切换日志库为Blg4go
func init_blog4goFun() {
	err := blg4go.NewWriterFromConfigAsFile("blog4go_config.xml")
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// defer blg4go.Close()

	// blg4go.SetTimeRotated(true)
}

func init() {
	init_blog4goFun()
}




func TstBlg4Fun() {

	index := 0
	for i := 0; i < 1000 * 10000; i ++ {

		time.Sleep( time.Second )

		index ++

		blg4go.Info( "test blg4go info ---index=", index )

		curPath, _ := GetCurrentPath()
	
		blg4go.Infof( "test blg4go infof PrintPwd=%s ---index=%d\n", curPath, index )
	
		blg4go.Errorf( "test blg4go Errorf PrintPwd=%s ---index=%d\n", curPath, index )
	
	}

}


