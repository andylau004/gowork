package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"tstfun/pool"

	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/pquerna/ffjson/ffjson"
	"golang.org/x/sys/unix"

	// "reflect"

	"Utilgo"
)

var (
	ONE_PIECE_SIZE = 4 * 1024

	TYPE_FILE   = 1
	TYPE_FOLDER = 2

	OPR_RENAME = 1
	OPR_DEL    = 2

	srvAddr = "172.17.0.2:80"
	strHost = "172.17.0.2"

	// srvAddr = "192.168.6.107:80"
	// strHost = "192.168.6.107"
)

func i64Tostring(val int64) string {
	strPId := strconv.FormatInt(val, 10)
	return strPId
}

func GetConnAndPubHead(cmdUrl string) (net.Conn, string) {
	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)

	var strHttpHead string
	strHttpHead += fmt.Sprintf("POST /cloudfile/v1/%s HTTP/1.1\r\n", cmdUrl)
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"
	strHttpHead += "skdId: 0\r\n"
	return conn, strHttpHead
}

func dummyFunc_http() {
	var retBuffer []byte
	url_Upfile := "http://172.17.0.2:80/cloudfile/v1/encryptupfile"

	req, err := http.NewRequest("POST", url_Upfile, bytes.NewBuffer(retBuffer))
	// req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func GetCurTId() string {
	strId := fmt.Sprintf("%d ", unix.Gettid())
	return strId
}
func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)

	return n
}

func GetGoroutineIDStr() string {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)

	u10 := strconv.FormatUint(n, 10)
	// fmt.Printf("%T, %v\n", s10, s10)

	return fmt.Sprintf("%v", u10)
}

type T struct{}

// func (t *T) A() {
// 	fmt.Println( "aaaaaaaaaa" )
// }
// func (t *T) B() {
// 	fmt.Println( "bbbbbbbbbb" )
// }
func (t T) A() {
	fmt.Println("aaaaaaaaaa")
}
func (t T) B() {
	fmt.Println("bbbbbbbbbb")
}

type Ter interface {
	A()
	B()
}

func identity(z *T) *T {
	return z
}
func ref(z T) *T {
	return &z
}

func tst15() {

	var val int
	val = 2
	workcount := 0
	for i := 0; i < 15; i++ {
		val *= 2
		workcount++
	}

	fmt.Println("val=", val)
	fmt.Println("workcount=", workcount)

	var other int
	other = 2 << 4
	fmt.Println("other=", other)

	var tmp1 int
	tmp1 = -1
	tmp1 = (-1 << 3)
	fmt.Println("tmp1=", tmp1)
}

// const srvAddr string = "172.17.0.2:444"

// connection pool
var g_connpoolObj pool.Pool

// < 连接对象, 使用次数 >
var g_mapUseCount sync.Map

func clientImpl(wg *sync.WaitGroup) {
	defer wg.Done()

	curId := unix.Gettid()

	v, err := g_connpoolObj.Get()
	if err != nil {
		fmt.Println(curId, " fatal error Get conn failed!!! err=", err)
		return
	}
	defer g_connpoolObj.Put(v)

	newConn := v.(net.Conn)
	fmt.Println(curId, " Get Conn Obj=", newConn)

	{
		newConn.Write([]byte("hello server\n"))
		// sendlen, err := newConn.Write([]byte("hello server\n"))
		// fmt.Println(curId, " sendlen=", sendlen, ", err=", err)
		respBuf := make([]byte, 2048)
		newConn.Read(respBuf)
		// recLen, err := newConn.Read(respBuf)
		// fmt.Println(curId, " recLen=", recLen, " respBuf=", string(respBuf))
	}

}

func clientFun() {
	var err error
	//factory 创建连接的方法
	factory := func() (interface{}, error) { return net.Dial("tcp", srvAddr) }

	//close 关闭连接的方法
	close := func(v interface{}) error { return v.(net.Conn).Close() }

	//创建一个连接池： 初始化5，最大连接30
	poolConfig := &pool.Config{
		InitialCap: 5,
		MaxCap:     30,
		Factory:    factory,
		Close:      close,
		//连接最大空闲时间，超过该时间的连接 将会关闭，可避免空闲时连接EOF，自动失效的问题
		IdleTimeout: 15 * time.Second,
	}

	g_connpoolObj, err = pool.NewChannelPool(poolConfig, &g_mapUseCount)
	if err != nil {
		fmt.Println("Create channelpool err=", err)
		return
	}
	lenpool := g_connpoolObj.Len()
	fmt.Println("conn pool len=", lenpool)

	tCount := 10
	wg := &sync.WaitGroup{}
	wg.Add(tCount)

	for i := 0; i < tCount; i++ {
		go clientImpl(wg)
	}
	wg.Wait()

	lenpool = g_connpoolObj.Len()
	fmt.Println("conn pool len=", lenpool)

	//释放连接池中的所有连接
	//p.Release()

	g_connpoolObj.UseCount()
}

func TstTcpConnPool() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)

	clientFun()
	fmt.Println("use ctrl + c exit")
	<-c
	fmt.Println("all work done")
}

func tst1() {
	fmt.Println("tst1 .............")
}
func tst2() {
	fmt.Println("tst2 .............")
}

func TstDefer() {
	i := 1

	{
		if i == 1 {
			defer tst2()
			fmt.Println("abcdefd")
		} else {
			// defer tst2()
		}
	}
	{
		if i == 1 {
			defer tst1()
			fmt.Println("fghijklmn")
		} else {
			// defer tst2()
		}

	}
	fmt.Println("12345677")

}

type TransferFileBean struct {
	FileBuffer []byte
	FileSize   int
}

var g_chFileBean chan *TransferFileBean = make(chan *TransferFileBean, 1000)

// var g_chFileBean chan TransferFileBean = make(chan TransferFileBean)

func ConsumeData(wg *sync.WaitGroup) {

	// for {
	// 	select {
	// 	// case recvBean := <-g_chFileBean:
	// 	// 	fmt.Println("recvBean buffer=", &byteBuffer1, ", recvBean len=", recvBean.)
	// 	}
	// }

	recvBean := <-g_chFileBean
	fmt.Printf("ConsumeData buffer1=%p len1=%d buffer1=%s\n",
		&(recvBean.FileBuffer), recvBean.FileSize, string(recvBean.FileBuffer))
	wg.Done()

	{
		recvBean := <-g_chFileBean
		// fmt.Println("recvBean buffer=", &(recvBean.FileBuffer), ", recvBean FileSize=", recvBean.FileSize)
		fmt.Printf("ConsumeData recvBean2 buffer=%p len2=%d buffer2=%s\n",
			&(recvBean.FileBuffer), recvBean.FileSize, string(recvBean.FileBuffer))
		wg.Done()
	}

}
func TstChByte() {

	wg := &sync.WaitGroup{}
	wg.Add(2)

	{
		byteBuffer1 := make([]byte, 1024)
		copy(byteBuffer1, []byte("12345"))
		fmt.Printf(" byteBuffer1=%p\n", &(byteBuffer1))

		// var oneBean TransferFileBean
		oneBean := &TransferFileBean{}
		oneBean.FileBuffer = byteBuffer1
		oneBean.FileSize = len(string(oneBean.FileBuffer))
		len1 := oneBean.FileSize

		fmt.Printf(" oneBean1.FileBuffer=%p , len1=%d\n", &(oneBean.FileBuffer), len1)

		// fmt.Println("before1")
		g_chFileBean <- oneBean
		// fmt.Println("over1")
	}

	{
		byteBuffer2 := make([]byte, 1024)
		copy(byteBuffer2, []byte("678910"))
		// len2 := len(byteBuffer2)
		fmt.Printf(" byteBuffer2=%p\n", &(byteBuffer2))

		// var oneBean TransferFileBean
		oneBean := &TransferFileBean{}
		oneBean.FileBuffer = byteBuffer2[0:]
		// oneBean.FileSize = len2
		oneBean.FileSize = len(oneBean.FileBuffer)
		len2 := oneBean.FileSize

		fmt.Printf(" oneBean2.FileBuffer=%p , len2=%d\n", &(oneBean.FileBuffer), len2)

		// fmt.Println("before2")
		g_chFileBean <- oneBean
		// fmt.Println("after2")
	}

	go ConsumeData(wg)
	wg.Wait()
}

func TstThreadId() {

	workCount := 20
	wg := &sync.WaitGroup{}
	wg.Add(workCount)

	for i := 0; i < workCount; i++ {

		go func() {
			defer wg.Done()
			strCurId := GetGoroutineIDStr()

			for {
				time.Sleep(2 * time.Second)
				fmt.Println(strCurId, " is working")
			}

		}()
	}

	wg.Wait()
}

type OneFile struct {
	UserId int64 `json:"userId"`

	FileId       int64  `json:"fileId"`
	OwnerId      int64  `json:"ownerId"`
	FileName     string `json:"fileName"`
	FileType     int    `json:"fileType"`
	FileSuffix   string `json:"fileSuffix"`
	FilePath     string `json:"filePath"`
	FileSize     int64  `json:"fileSize"`
	SecretKey    string `json:"secretKey"`
	UploaderId   int64  `json:"uploaderId"`
	UploaderName string `json:"uploaderName"`
	Md5Hash      string `json:"md5"`
	Sha2Hash     string `json:"sha2"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
	DeletedAt    int64  `json:"deletedAt"`
	// OrderNum int    `json:"orderNum"`
	// RoleId int    `json:"roleId"`
	// Members int    `json:"members"`
}
type SyncFileRequest struct {
	UserId    int64 `json:"userId"`
	TimeStamp int64 `json:"timeStamp"`
}

func GetFileSize(fileFullPath string) int64 {
	fileInfo, _ := os.Stat(fileFullPath)
	filesize := fileInfo.Size()
	fmt.Println("fileName=", fileFullPath, ", fileSize=", filesize) //返回的是字节
	return filesize
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s\n", err.Error())
		// os.Exit(1)
	}
}

func PostUploadEmptyPath() {
	// Empty Path
	upEmptyPath := "EmptyPath-1"
	upEmptyPath = "EmptyPath-12389"
	upEmptyPath = "EmptyPath-9981"
	upEmptyPath = "EmptyPath-6636/"
	upEmptyPath = "myspace/path1"

	// var fileBean OneFile
	// fileBean.UserId = 1234567
	// fileBean.OwnerId = fileBean.UserId

	// fileBean.FileName = upEmptyPath
	// fileBean.FileType = 2

	// fileBean.FileSuffix = ""
	// fileBean.FileSize = 0
	// fileBean.Md5Hash = "asdfasdfadfasd=="
	// fileBean.Sha2Hash = "tsasdf123901234=="
	// fileBean.SecretKey = "1234456"

	// retBuffer, _ := ffjson.Marshal(&fileBean)
	// fmt.Println(" upload FileJson =", string(retBuffer))
	conn, strHttpHead := GetConnAndPubHead("encryptupfile")

	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "sdkId: 0\r\n"
	// strHttpHead += "fileId: 9981\r\n"
	strHttpHead += "fileType: 2\r\n"
	strHttpHead += fmt.Sprintf("fileName: %s\r\n", upEmptyPath)
	strHttpHead += "\r\n"

	// send http head
	nSendLen, err := conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	respBuf := make([]byte, 2048)
	conn.Read(respBuf)
	fmt.Println("post empty folder respBuf=", string(respBuf))

	return
	fmt.Println(nSendLen)
}

func PostMoidfyFile() {
	upFileName := "myspace/1.go"
	fileSize := GetFileSize(upFileName)

	// var fileBean OneFile
	// fileBean.UserId = 1234567
	// fileBean.OwnerId = fileBean.UserId
	// // 文件ID=53；此文件已经存在，上传被修改后的已存在文件；
	// fileBean.FileId = 53

	// fileBean.FileName = "logs/filexixi"
	// fileBean.FileType = 1

	// fileBean.FileSuffix = "log"
	// // fileBean.FilePath = "./logs"
	// fileBean.FileSize = fileSize
	// fileBean.SecretKey = "1234456"
	// fileBean.Md5Hash = "123hh"
	// fileBean.Sha2Hash = "456hh"

	// retBuffer, _ := ffjson.Marshal(&fileBean)
	// fmt.Println(" upload FileJson =", string(retBuffer))

	var strHttpHead string
	strHttpHead += "POST /cloudfile/v1/encryptupfile HTTP/1.1\r\n"
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"

	strHttpHead += fmt.Sprintf("fileName: %s\r\n", upFileName)
	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "fileId: 99\r\n"
	// strHttpHead += fmt.Sprintf("Content-Length: %d\r\n", (int64)(fileSize))
	strHttpHead += fmt.Sprintf("fileSize: %d\r\n", fileSize)
	strHttpHead += "sha2: newsha2777\r\n"
	strHttpHead += "md5: newmd5777\r\n"
	strHttpHead += "secretKey: 6655789\r\n"
	strHttpHead += "fileType: 1\r\n"

	strHttpHead += fmt.Sprintf("Content-Length: %d\r\n", (int64)(fileSize))
	strHttpHead += "\r\n"

	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)
	// send http head
	var nSendLen int
	nSendLen, err = conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	// send json file metaInfo
	// nSendLen, err = conn.Write(Utilgo.Str2bytes(string(retBuffer)))
	// checkError(err)

	// send file content
	var sendSize int64
	sendCount := 0
	pieceSize := fileSize / (int64)(ONE_PIECE_SIZE)
	remainSize := fileSize - (int64)((int)(pieceSize)*(int)(ONE_PIECE_SIZE))
	{
		if remainSize != 0 {
			pieceSize++
		}
		//
		fp, err := os.Open(upFileName) // 获取文件指针
		if err != nil {
			fmt.Println("open file faild!!! filename=", upFileName, ", err=", err)
			return
		}
		defer fp.Close()

		readBuffer := make([]byte, ONE_PIECE_SIZE)
		for {
			// 注意这里要取bytesRead, 否则有问题
			bytesRead, err := fp.Read(readBuffer) // 文件内容读取到buffer中
			if bytesRead == 0 {
				// fmt.Println("last sendCount=", sendCount)
				// fmt.Println("bytesRead == 0")
				break
			}
			// fmt.Println("bytesRead=", bytesRead, ", err=", err)

			nSendLen, err = conn.Write(readBuffer[:bytesRead])
			checkError(err)
			// fmt.Println("nSendLen=", nSendLen, ", err=", err)
			sendCount++
			sendSize += (int64)(nSendLen)

			if err != nil {
				if err == io.EOF {
					err = nil
					break
				} else {
					fmt.Println("some error happened")
					return
				}
			}

		} // for end

	}

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	// fmt.Println("resp beg----------------------")
	fmt.Println("post modifyfile respBuf=", string(respBuf))
	// fmt.Println("resp end----------------------")

	return
	fmt.Println(nSendLen)
	fmt.Println("pieceSize=", pieceSize)
	fmt.Println("sendSize=", sendSize, ", fileSize=", fileSize)
}

type HttpResponse struct {
	RetCode   int   `json:"retCode"`
	RetFileId int64 `json:"retfileId"`
	TimeStamp int64 `json:"timeStamp"`
	OprTime   int64 `json:"oprTime"`
}

func PostUploadFile(parentId int64, upFileName string, wgWork *sync.WaitGroup) int64 {
	defer wgWork.Done()
	// endTag := strings.Index(upFileName, "1.go")
	// if endTag < 0 {
	// 	return
	// }

	fmt.Println("\n------------->PostUploadFile beg--------------")
	defer fmt.Println("------------->PostUploadFile end--------------\n")

	var fileSize int64
	if IsFile(upFileName) {
		fmt.Println("before  ----getfilesize")
		fileSize = GetFileSize(upFileName)
		fmt.Println("after  ----getfilesize")
		fmt.Println("upload onefile=", upFileName)
	}
	if IsDir(upFileName) {
		fmt.Println("upload onefolder=", upFileName)
	}

	strHttpHead := "POST /cloudfile/v1/encryptupfile HTTP/1.1\r\n"
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"
	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += fmt.Sprintf("fileName: %s\r\n", upFileName)
	// strHttpHead += "fileId: 33\r\n"
	strHttpHead += fmt.Sprintf("PID: %s\r\n", i64Tostring(parentId))

	tmpType := -1
	if IsDir(upFileName) {
		// ret := PostCheckFolderExist(1234567, upFileName)
		// if ret == -3 {
		// 	fmt.Println("post dir: ", upFileName, " exist!!!")
		// 	return
		// }
		// if ret == -4 {
		// 	// fmt.Println("post dir: ", upFileName, " no found!!!")
		// }
		strHttpHead += "fileType: 2\r\n"
		tmpType = 2
	} else if IsFile(upFileName) {
		tmpMd5 := "88f6e0c05304eb77bd52cebf8e942a23545055fee17cad822839a0db118dc26md5"
		tmpSha2 := "88f6e0c05304eb77bd52cebf8e942a23545055fee17cad822839a0db118dc26aa"

		ret := PostCheckFileExist(tmpMd5, tmpSha2, fileSize)
		if ret == -3 {
			fmt.Println("post file: ", upFileName, " exist!!!")
			return (int64)(ret)
		}
		if ret == -4 {
			// fmt.Println("post file: ", upFileName, " no found!!!")
		}
		strHttpHead += fmt.Sprintf("Content-Length: %d\r\n", (int64)(fileSize))
		strHttpHead += fmt.Sprintf("fileSize: %d\r\n", fileSize)
		strHttpHead += fmt.Sprintf("sha2: %s\r\n", tmpSha2)
		strHttpHead += fmt.Sprintf("md5: %s\r\n", tmpMd5)
		strHttpHead += "secretKey: aasdfa123\r\n"
		strHttpHead += "fileType: 1\r\n"
		tmpType = 1
	} else {
		fmt.Println("unknown ...file type...")
		return -1
	}
	strHttpHead += "\r\n"

	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)

	// send http head
	var nSendLen int
	nSendLen, err = conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	// send json file metaInfo
	// nSendLen, err = conn.Write(Utilgo.Str2bytes(string(retBuffer)))
	// checkError(err)

	// send file content
	var sendSize int64
	sendCount := 0
	pieceSize := fileSize / (int64)(ONE_PIECE_SIZE)
	remainSize := fileSize - (int64)((int)(pieceSize)*(int)(ONE_PIECE_SIZE))
	if IsFile(upFileName) {
		if remainSize != 0 {
			pieceSize++
		}
		fp, err := os.Open(upFileName) // 获取文件指针
		if err != nil {
			fmt.Println("open file faild!!! filename=", upFileName, ", err=", err)
			return -1
		}
		defer fp.Close()

		readBuffer := make([]byte, ONE_PIECE_SIZE)
		for {
			// 注意这里要取bytesRead, 否则有问题
			bytesRead, err := fp.Read(readBuffer) // 文件内容读取到buffer中
			if bytesRead == 0 {
				// fmt.Println("last sendCount=", sendCount)
				// fmt.Println("bytesRead == 0")
				break
			}
			// fmt.Println("bytesRead=", bytesRead, ", err=", err)

			nSendLen, err = conn.Write(readBuffer[:bytesRead])
			checkError(err)
			// fmt.Println("nSendLen=", nSendLen, ", err=", err)
			sendCount++
			sendSize += (int64)(nSendLen)

			if err != nil {
				if err == io.EOF {
					err = nil
					break
				} else {
					fmt.Println("some error happened")
					return -1
				}
			}

		} // for end
	}

	respBuf := make([]byte, 20*2048)
	recvLen, errR := conn.Read(respBuf)
	// fmt.Println("recvLen=", recvLen)
	// fmt.Println("resp beg----------------------")
	if tmpType == 2 {
		fmt.Println("post dir response=", string(respBuf))
	} else if tmpType == 1 {
		fmt.Println("post file response=", string(respBuf))
	} else {
		fmt.Println("post unknown type!!! response=", string(respBuf))
		return -1
	}
	// fmt.Println("resp end----------------------")

	tmpbytes := make([]byte, recvLen)
	copy(tmpbytes, respBuf)

	pos := strings.Index(string(tmpbytes), "\r\n\r\n")
	startPos := pos + len("\r\n\r\n")

	tmpstr := tmpbytes[startPos:]
	// fmt.Println("len(tmpstr)=", len(tmpstr))

	var repsHttp HttpResponse
	err1 := ffjson.Unmarshal([]byte(tmpstr), &repsHttp)
	// fmt.Printf(" tmpstr---->[%s]", tmpstr)
	// for i, ch := range tmpstr {
	// 	fmt.Println(i, ch) //ch的类型为rune 默认utf-8编码，一个汉字三个字节
	// }
	// fmt.Println("err1=", err1, " tmpstr---->", string(tmpstr))
	// fmt.Printf("resh=%+v\n", repsHttp)

	return repsHttp.RetFileId

	fmt.Println("err1=", err1, " tmpstr---->", string(tmpstr))
	fmt.Println(errR)
	fmt.Println("pieceSize=", pieceSize)
	fmt.Println("sendSize=", sendSize, ", fileSize=", fileSize)

	fmt.Println(nSendLen)
	fmt.Println(tmpType)
	return -1
}

func TstSyncRequest() {
	// var syncReq SyncFileRequest
	// syncReq.UserId = 1234567
	// syncReq.TimeStamp = 468
	// // syncReq.TimeStamp = 4680000
	// retBuffer, _ := ffjson.Marshal(&syncReq)
	// fmt.Println(" syncreq =", string(retBuffer))
	conn, strHttpHead := GetConnAndPubHead("syncfiles")

	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "timeStamp: 2\r\n"
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err := conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	// send json file metaInfo
	// nSendLen, err = conn.Write(Utilgo.Str2bytes(string(retBuffer)))
	// checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("sync user file/path respBuf=", string(respBuf))

	return
	fmt.Println(nSendLen)
}

type PullFileRequest struct {
	UserId int64 `json:"userId"`
	SdkId  int64 `json:"sdkId"`
	FileId int64 `json:"fileId"`
}

func TstPullFileRequest() {
	conn, strHttpHead := GetConnAndPubHead("pullfilelist")

	// strHttpHead += fmt.Sprintf("fileName: %s\r\n", upFileName)
	// strHttpHead += "userId: 1234567\r\n"
	// strHttpHead += "fileId: 99\r\n"
	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "fileId: 33\r\n"
	// strHttpHead += fmt.Sprintf("Content-Length: %d\r\n", (int64)(fileSize))
	// strHttpHead += fmt.Sprintf("fileSize: %d\r\n", fileSize)
	// strHttpHead += "sha2: newsha2777\r\n"
	// strHttpHead += "md5: newmd5777\r\n"
	// strHttpHead += "secretKey: 6655789\r\n"
	strHttpHead += "fileType: 1\r\n"
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err := conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	{
		saveFileName := "./save-pull"
		os.Remove(saveFileName)

		f, err := os.OpenFile(saveFileName, os.O_WRONLY|os.O_CREATE, 0666)
		// f, err := os.OpenFile("", os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Println("write file create failed. err: " + err.Error())
			return
		}
		defer f.Close()

		var endHead int
		bJson := true

		writeBytes := 0
		for {
			RecvBuf := make([]byte, 10*1024)
			// fmt.Println("----------------------before Read")
			recvLen, err := conn.Read(RecvBuf)
			findPos := strings.Index(string(RecvBuf), "HTTP/1.1 404 ")
			if findPos >= 0 {
				fmt.Println("no found file 404------------------!!!, RecvBuf=", string(RecvBuf))
				return
			} else {
				fmt.Println("findpos=", findPos, ", RecvBuf=", string(RecvBuf))
			}
			// fmt.Println("----------------------after Read, recvLen=", recvLen)

			if err != nil {
				fmt.Println("after Read, err=", err)
				if err != io.EOF {
					//Error Handler
					fmt.Println("err != io.EOF err=", err)
				}
				break
			}

			if recvLen > 0 {
				// fmt.Println("recvLen=", recvLen)
				if bJson {
					// fmt.Println("1111111111 RecvBuf=", string(RecvBuf))
					endHead = strings.Index(string(RecvBuf), "\r\n\r\n")
					endHead = endHead + len("\r\n\r\n")
					// fmt.Println("pos=", endHead)
					// fmt.Println("333333333333 xtts=", string(RecvBuf[endHead:]))

					responseHead := RecvBuf[0:endHead]
					fmt.Println("responseHead=", string(responseHead))

					tmpLen := recvLen - endHead
					// fmt.Println("endHead=", endHead, " tmpLen=", tmpLen)

					_, err = f.Write(RecvBuf[endHead:recvLen])
					writeBytes += tmpLen

					bJson = false
				} else {
					// fmt.Println("2222222222 RecvBuf=", string(RecvBuf))
					_, err = f.Write(RecvBuf[0:recvLen])
					writeBytes += recvLen
				}

				// fmt.Println("end write file, recvLen=", recvLen, "writeBytes=", writeBytes, ", err=", err)
			} else {
				fmt.Println("recvLen == 0, err=", err)
				break
			}

		}
		fmt.Println("end write file, writeBytes=", writeBytes)
	}

	return
	fmt.Println(nSendLen)
}

func PostRenameFile() {
	// tmpstr := "{\"userId\":1234567,\"fileId\":2,\"ownerId\":1234567,\"newName\":\"logs/rename1\",\"fileType\":1,\"oprType\":1}"
	// fmt.Println("len data=", len(tmpstr))
	// return
	conn, strHttpHead := GetConnAndPubHead("handlefile")

	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "fileId: 50\r\n"
	strHttpHead += fmt.Sprintf("newName: %s\r\n", "myspace/path3/guoqing--file333")
	strHttpHead += "fileType: 1\r\n"
	strHttpHead += "oprType: 1\r\n"
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err := conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	// send json file
	// nSendLen, err = conn.Write(Utilgo.Str2bytes(string(requestJson)))
	// checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	fmt.Println("resp beg----------------------")
	fmt.Println("rename file respBuf=", string(respBuf))
	fmt.Println("resp end----------------------")

	return
	fmt.Println(nSendLen)
}

func PostMvFile() {
	conn, strHttpHead := GetConnAndPubHead("mv")

	// 删除用户9999的71目录
	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "mvFileIds: 91\r\n"
	strHttpHead += "targetId: 50\r\n"
	strHttpHead += "\r\n"

	// send http head
	nSendLen, err := conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	fmt.Println("mvfile resp beg----------------------")
	fmt.Println("mvfile  respBuf=", string(respBuf))
	fmt.Println("mvfile resp end----------------------")

	return
	fmt.Println(nSendLen)
}

func PostDelFile() {
	// var Delreq ReDelRequest
	// Delreq.UserId = 1234567
	// Delreq.OwnerId = Delreq.UserId
	// Delreq.SdkId = 0
	// Delreq.FileId = 2
	// Delreq.FileId = 53

	// Delreq.NewName = "logs/recoverfilename"
	// Delreq.FileType = TYPE_FILE // just file
	// Delreq.OprType = OPR_DEL    // del file type

	// requestJson, _ := ffjson.Marshal(&Delreq)
	// fmt.Println(" del Delreq =", string(requestJson))

	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)

	var strHttpHead string
	strHttpHead += "POST /cloudfile/v1/handlefile HTTP/1.1\r\n"
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"

	strHttpHead += "skdId: 0\r\n"
	// 删除用户1234567的48目录
	// strHttpHead += "userId: 1234567\r\n"
	// strHttpHead += "fileId: 48\r\n"
	// strHttpHead += "delFileIds: 48,\r\n"

	// 删除用户9999的71目录
	strHttpHead += "userId: 9999\r\n"
	strHttpHead += "delFolderIds: 71\r\n"

	strHttpHead += "oprType: 2\r\n" // 1：重命名；2：删除；
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err = conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	// send json file
	// nSendLen, err = conn.Write(Utilgo.Str2bytes(string(requestJson)))
	// checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	fmt.Println("delfile resp beg----------------------")
	fmt.Println("delfile  respBuf=", string(respBuf))
	fmt.Println("delfile resp end----------------------")

	return
	fmt.Println(nSendLen)
}

func PostRenameFolder() {
	// var req ReDelRequest
	// req.UserId = 1234567
	// req.OwnerId = req.UserId
	// req.SdkId = 0
	// req.FileId = 2
	// req.FileId = 53
	// req.NewName = "logs/FolderAAA"
	// req.FileType = TYPE_FOLDER
	// req.OprType = OPR_RENAME
	// requestJson, _ := ffjson.Marshal(&req)
	// fmt.Println(" rename req =", string(requestJson))
	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)

	var strHttpHead string
	strHttpHead += "POST /cloudfile/v1/handlefile HTTP/1.1\r\n"
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"

	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "fileId: 24\r\n"
	strHttpHead += fmt.Sprintf("newName: %s\r\n", "myspace/path3--newname")
	strHttpHead += "fileType: 2\r\n"
	strHttpHead += "oprType: 1\r\n"
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err = conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	// send json file
	// nSendLen, err = conn.Write(Utilgo.Str2bytes(string(requestJson)))
	// checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	fmt.Println("resp beg----------------------")
	fmt.Println("rename file respBuf=", string(respBuf))
	fmt.Println("resp end----------------------")

	return
	fmt.Println(nSendLen)
}

func GetAllFile(pathname string) error {
	rd, err := ioutil.ReadDir(pathname)
	for _, fi := range rd {
		if fi.IsDir() {
			fmt.Printf("%s\n", pathname+"\\"+fi.Name())
			GetAllFile(pathname + fi.Name() + "\\")
		} else {
			fmt.Println(fi.Name())
		}
	}
	return err
}
func GetAllFileAndDirEx(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}
	for _, fi := range rd {
		if fi.IsDir() {
			fullDir := pathname + "/" + fi.Name()
			s = append(s, fullDir)

			s, err = GetAllFileAndDirEx(fullDir, s)
			if err != nil {
				fmt.Println("read dir fail:", err)
				return s, err
			}
		} else {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
			// fmt.Println("fullName=", fullName)
		}
	}
	return s, nil
}

func ExtractPath() []string {
	var s []string
	s, _ = GetAllFileAndDirEx("myspace", s)
	// fmt.Printf("slice: %+v", s)
	// for _, fil := range s {
	// 	fmt.Println("fil=", fil)
	// }
	return s
}

func PostAllPathFile() {
	// defer fmt.Println("all upload thread done")

	// retAllFiles := ExtractPath()
	// fileCount := len(retAllFiles)
	// fmt.Println("size=", len(retAllFiles))

	// var wg sync.WaitGroup
	// wg.Add(fileCount)

	// for _, file := range retAllFiles {
	// 	PostUploadFile(0, file, &wg)
	// }
	// wg.Wait()
}

func PostCheckFolderExist(userId int64, folderId int64, PId int64) int {
	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)

	var strHttpHead string
	strHttpHead += "POST /cloudfile/v1/encryptupfile HTTP/1.1\r\n"
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += fmt.Sprintf("userId: %d\r\n", userId)

	{ // Folder type
		strHttpHead += "fileType: 2\r\n"

		strCurId := i64Tostring(folderId)
		strHttpHead += fmt.Sprintf("fileId: %s\r\n", strCurId)

		strPId := i64Tostring(PId)
		strHttpHead += fmt.Sprintf("PID: %s\r\n", strPId)
	}
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err = conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	fmt.Println("resp beg----------------------")
	fmt.Println("CheckFolderExist folder respBuf=", string(respBuf))
	fmt.Println("resp end----------------------")

	findPos := strings.Index(string(respBuf), "\"retCode\":-3,")
	if findPos > 0 {
		return -3 // find exist
	}

	findPos = strings.Index(string(respBuf), "\"retCode\":-4,")
	if findPos > 0 {
		return -4 // can't found
	} else {
		return 0 // all ok
	}
	fmt.Println(nSendLen)
	return 0
}

func PostCheckFileExist(md5, sha2 string, fileSize int64) int {
	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)

	var strHttpHead string
	strHttpHead += "POST /cloudfile/v1/preuploadFile HTTP/1.1\r\n"
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"
	strHttpHead += "Content-Type: application/json\r\n"

	strHttpHead += "fileType: 1\r\n"
	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "sdkId: 0\r\n"
	strHttpHead += fmt.Sprintf("md5: %s\r\n", md5)
	strHttpHead += fmt.Sprintf("sha2: %s\r\n", sha2)
	strHttpHead += fmt.Sprintf("fileSize: %d\r\n", fileSize)
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err = conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	// send json file
	// nSendLen, err = conn.Write(Utilgo.Str2bytes(string(requestJson)))
	// checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	fmt.Println("resp beg----------------------")
	fmt.Println("CheckFileExist file respBuf=", string(respBuf))
	fmt.Println("resp end----------------------")

	findPos := strings.Index(string(respBuf), "\"retCode\":-3,")
	if findPos > 0 {
		fmt.Println("CheckFileExist find exist!!!")
		return -3 // find exist
	} else {
		findPos = strings.Index(string(respBuf), "\"retCode\":-4,")
		if findPos > 0 {
			fmt.Println("CheckFileExist no found!!!")
		}
		return 0 // all ok
	}
	fmt.Println(nSendLen)
	return 0
}

type ReDelRequest struct {
	UserId int64 `json:"userId"`
	SdkId  int64 `json:"sdkId"`

	FileId  int64  `json:"fileId"`
	OwnerId int64  `json:"ownerId"`
	OldName string `json:"oldName"`
	NewName string `json:"newName"`

	FileType int `json:"fileType"` // 1: file; 2: path;
	OprType  int `json:"oprType"`  // 1: rename; 2: delete;
}

func PostDelFolder() {
	conn, err := net.Dial("tcp", srvAddr)
	checkError(err)

	var strHttpHead string
	strHttpHead += "POST /cloudfile/v1/handlefile HTTP/1.1\r\n"
	strHttpHead += fmt.Sprintf("Host: %s\r\n", strHost)
	strHttpHead += "User-Agent: Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:23.0) Gecko/20100101 Firefox/23.0\r\n"
	strHttpHead += "Content-Type: application/json\r\n"
	strHttpHead += "Connection: Keep-Alive\r\n"
	strHttpHead += "Content-Type: application/json\r\n"

	strHttpHead += "userId: 1234567\r\n"
	strHttpHead += "sdkId: 0\r\n"
	strHttpHead += "oldName: myspace/path1\r\n"
	strHttpHead += "fileType: 2\r\n" // 1: 文件；2：文件夹；
	strHttpHead += "oprType: 2\r\n"  // 1：重命名；2：删除；

	// strHttpHead += "\r\n"
	strHttpHead += "\r\n"

	// send http head
	var nSendLen int
	nSendLen, err = conn.Write(Utilgo.Str2bytes(strHttpHead))
	checkError(err)

	respBuf := make([]byte, 20*2048)
	conn.Read(respBuf)
	fmt.Println("")
	fmt.Println("resp beg----------------------")
	fmt.Println("del Folder respBuf=", string(respBuf))
	fmt.Println("resp end----------------------")
	return
	fmt.Println(nSendLen)
}

// Creates a new file upload http request with optional extra params
func newfileUploadRequest(uri string, params map[string]string, paramName, filePath string) (*http.Request, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fmt.Println("filepath.Base(filePath)=", filepath.Base(filePath))

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	// part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	// if err != nil {
	// 	fmt.Println("create form file failed, err=", err)
	// 	return nil, err
	// }
	// _, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	part, err := writer.CreateFormFile(paramName, filepath.Base(filePath))
	if err != nil {
		fmt.Println("create form file failed, err=", err)
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		fmt.Println("writer close failed, err=", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	fmt.Println("writer.FormDataContentType()=", writer.FormDataContentType())
	return req, err
}

func postmulitifileImpl() {
	filePath, _ := os.Getwd()
	// filePath += "/myspace/path3/file333"
	filePath = "myspace/1.go"
	fmt.Println("up new file=", filePath)

	fileSize := GetFileSize(filePath)
	// fmt.Println(fileSize)

	extraParams := map[string]string{}

	justName := ""
	posSlash := strings.LastIndex(filePath, "/")
	if posSlash > 0 {
		justName = filePath[posSlash+1:]
	}
	extraParams["fileName"] = fmt.Sprintf("%s", justName)

	extraParams["fileSize"] = fmt.Sprintf("%d", fileSize)
	extraParams["md5"] = fmt.Sprintf("%s", "11asdfamd5")
	extraParams["sha2"] = fmt.Sprintf("%s", "11asdfasha2")
	extraParams["userId"] = fmt.Sprintf("%d", 1234567)
	extraParams["sdkId"] = fmt.Sprintf("%d", 0)
	extraParams["PID"] = fmt.Sprintf("%d", 0)
	fmt.Println("extraParams=", extraParams)

	url := fmt.Sprintf("http://%s/cloudfile/v1/encryptupfile", srvAddr)
	request, err := newfileUploadRequest(url, extraParams, "file", filePath)
	if err != nil {
		fmt.Println("post file failed! err=", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		// log.Fatal(err)
		fmt.Println("http client do failed, err=", err)
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			fmt.Println("read from failed, err=", err)
			// log.Fatal(err)
		}
		resp.Body.Close()
		fmt.Println("StatusCode=", resp.StatusCode)
		fmt.Println("Header=", resp.Header)
		fmt.Println("body=", body)
	}

}

func PostMultiFile() {
	postmulitifileImpl()
}

func postOneFolder() {
	var wg sync.WaitGroup
	wg.Add(4)

	// 23
	retFileId := PostUploadFile(0, "myspace", &wg)
	fmt.Println("retFileId = ", retFileId)

	// myspace/1.go
	retFileId = PostUploadFile(retFileId, "myspace/1.go", &wg)
	fmt.Println("retFileId = ", retFileId)
	// return

	// 24
	retFileId = PostUploadFile(retFileId, "myspace/path3", &wg)
	fmt.Println("retFileId = ", retFileId)

	// 27
	retFileId = PostUploadFile(retFileId, "myspace/path3/file333", &wg)
	fmt.Println("retFileId = ", retFileId)

	wg.Wait()
}
func CheckPIdPath() {
	// ret := PostCheckFolderExist(1234567, upFileName)
	// if ret == -3 {
	// 	fmt.Println("post dir: ", upFileName, " exist!!!")
	// 	return
	// }
	// if ret == -4 {
	// 	// fmt.Println("post dir: ", upFileName, " no found!!!")
	// }

	PostCheckFolderExist(1234567, 23, 0)
	PostCheckFolderExist(1234567, 24, 23)
	PostCheckFolderExist(1234567, 27, 24)
}
func TstYunPanClient() {
	PostMvFile()
	return

	PostRenameFile()
	return

	postOneFolder()
	return

	TstSyncRequest()
	return

	PostDelFile()
	return

	PostRenameFolder()
	return

	TstPullFileRequest()
	return

	// CheckPIdPath()
	// return

	PostAllPathFile()
	return

	PostMultiFile()
	return

	PostDelFolder()
	return

	PostUploadEmptyPath()
	return

	// PostCheckFolderExist(1234567, "tstspace/9981")
	return

	PostMoidfyFile()
	return

	ExtractPath()
	return

}

func ReadAll(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func GetFileContent(fileName string) ([]byte, error) {
	fp, err := os.Open(fileName) // 获取文件指针
	if err != nil {
		fmt.Println("open file faild!!! filename=", fileName, ", err=", err)
		return nil, err
	}
	defer fp.Close()

	readBuffer := make([]byte, ONE_PIECE_SIZE)
	for {
		// 注意这里要取bytesRead, 否则有问题
		bytesRead, err := fp.Read(readBuffer) // 文件内容读取到buffer中
		fmt.Println("bytesRead=", bytesRead, ", err=", err)

		// nSendLen, err = conn.Write(buffer[:bytesRead])
		checkError(err)
		// fmt.Println("nSendLen=", nSendLen, ", err=", err)

		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				fmt.Println("some error happened")
				return nil, err
			}
		}

	}
	return readBuffer, nil
}

func GetFileEEEE() {
	upFileName := "logs/debug.log.2019-07-28"
	// fileSize := GetFileSize(upFileName)

	fp, err := os.Open(upFileName) // 获取文件指针
	if err != nil {
		fmt.Println("open file faild!!! filename=", upFileName, ", err=", err)
		return
		// return nil, err
	}
	defer fp.Close()

	readBuffer := make([]byte, ONE_PIECE_SIZE)

	for {
		// 注意这里要取bytesRead, 否则有问题
		bytesRead, err := fp.Read(readBuffer) // 文件内容读取到buffer中
		fmt.Println("bytesRead=", bytesRead, ", err=", err)

		// nSendLen, err = conn.Write(buffer[:bytesRead])
		checkError(err)
		// fmt.Println("nSendLen=", nSendLen, ", err=", err)

		if err != nil {
			if err == io.EOF {
				err = nil
				break
			} else {
				fmt.Println("some error happened")
				return
			}
		}

	}

}

// fileName:文件名字(带全路径)
// content: 写入的内容
func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("write file create failed. err: " + err.Error())
	} else {
		// // 查找文件末尾的偏移量
		// n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		// _, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
}

// func appendToFile(fileName string, content []byte) error {
// 	// 以只写的模式，打开文件
// 	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_APPEND, 0666)
// 	if err != nil {
// 		fmt.Println("write file create failed. err: " + err.Error())
// 	} else {
// 		// 从末尾的偏移量开始写入内容
// 		_, err = f.Write(content)
// 	}
// 	defer f.Close()
// 	return err
// }

func AddFile() {

	appendToFile("save123", "12345")
	appendToFile("save123", "678910")
	appendToFile("save123", "abcdefgh")

}

type ShareItem struct {
	ShareFolderIds []int64 `json:"shareFolderIds"`
	WUserIds       []int64 `json:"wUserIds"`
	RUserIds       []int64 `json:"rUserIds"`
	WOrgIds        []int64 `json:"wOrgIds"`
	ROrgIds        []int64 `json:"rOrgIds"`
}
type ShareFolder struct {
	UserId     int64       `json:"userId"`
	SdkId      int64       `json:"sdkId"`
	ShareCount int         `json:"shareCount"`
	ShareGroup []ShareItem `json:"shareGroup"`
}

func tst_json() {
	/*
		{
			"userId": 2890,
			"sdkId": 0,
			"shareGroup":
			[
				{ "shareFolderIds":  [123, 124], "wOrgIds": [100, 101], "rOrgIds": [778, 779], "wUserIds": [9981, 9982],  "rUserIds": [9983, 9985] },
				{ "shareFolderIds":  [211, 276], "wOrgIds": [600, 601], "rOrgIds": [878, 879], "wUserIds": [1011, 1200],  "rUserIds": [761, 480] }
			]
		}
	*/

	var shareObj ShareFolder
	fileBuff, err := ReadAll("1.json")

	err = ffjson.Unmarshal(fileBuff, &shareObj)
	if err != nil {
		fmt.Println("unmarshal failed! err=", err)
	}
	fmt.Printf(" upload FileJson=%+v\n", shareObj)

	for i, v := range shareObj.ShareGroup {
		fmt.Printf("i=%d v=%+v\n", i, v)
	}

}

func TstTimeInterval() {

	tStart := time.Now()
	time.Sleep(time.Second * 2)
	fmt.Println("client.Do timeinterval=", time.Now().Sub(tStart))

}
func DialCustom(network, address string, timeout time.Duration, localIP []byte, localPort int) (net.Conn, error) {
	netAddr := &net.TCPAddr{Port: localPort}

	if len(localIP) != 0 {
		netAddr.IP = localIP
	}

	fmt.Println("netAddr:", netAddr)

	d := net.Dialer{Timeout: timeout, LocalAddr: netAddr}
	return d.Dial(network, address)
}

func Tstlocalipclient() {
	//url 要请求的URL
	serverAddr := "192.168.242.157:9981"

	// 172.28.0.180
	//localIP := []byte{0xAC, 0x1C, 0, 0xB4}  // 指定IP
	// localIP := []byte(string("127.0.0.1")) //  any IP，不指定IP
	localIP := []byte(string(""))
	localPort := 9001 // 指定端口
	conn, err := DialCustom("tcp", serverAddr, time.Second*2, localIP, localPort)
	if err != nil {
		fmt.Println("dial failed:", err)
		os.Exit(1)
	}
	defer conn.Close()

	for {
		nSendLen, err := conn.Write([]byte("asdfadsfadfasdfaaaaaaaaaaaaaaaaaaaa"))
		// checkError(err)
		fmt.Println("nsendlen=", nSendLen, ", err=", err)

		// buffer := make([]byte, 512)
		// reader := bufio.NewReader(conn)

		// n, err2 := reader.Read(buffer)
		// if err2 != nil {
		// 	fmt.Println("Read failed:", err2)
		// 	return
		// }

		// fmt.Println("count:", n, "msg:", string(buffer))

		time.Sleep(time.Second)
	}
	select {}
}

func tst_map_fun() {
	persons := make(map[string]int)
	persons["张三"] = 19

	mp := &persons

	fmt.Printf("原始map的内存地址是：%p\n", mp)
	modifyex(persons)
	fmt.Println("map值被修改了，新值为:", persons)
}

func modifyex(p map[string]int) {
	fmt.Printf("函数里接收到map的内存地址是：%p\n", &p)
	p["张三"] = 20
}

func JustTest() {
	stuff := []interface{}{"this", "that", "otherthing"}
	fmt.Println("repeat content=", strings.Repeat(",?", len(stuff)-1))

	sql := "select * from foo where id=? and name in (?" + strings.Repeat(",?", len(stuff)-1) + ")"
	fmt.Println("SQL:", sql)

	args := []interface{}{10}
	args = append(args, stuff...)
	fakeExec(args...)
	// This also works, but I think it's harder for folks to read
	//fakeExec(append([]interface{}{10},stuff...)...)
}

func fakeExec(args ...interface{}) {
	fmt.Println("Got:", args)
}

func BitOpr() {
	var databyte byte = 0x7f
	isLongHeader := databyte&0x80 > 0
	fmt.Println(isLongHeader)
}

func write(ch chan int) {
	fmt.Println("write beg -------")
	defer fmt.Println("write end -------")
	for i := 0; i < 5; i++ {
		ch <- i
		fmt.Println("successfully wrote", i, "to ch")
	}
	close(ch)
}
func Tt1() {
	ch := make(chan int, 2)

	go write(ch)
	time.Sleep(2 * time.Second)

	for v := range ch {
		fmt.Println("read value", v, "from ch")
		time.Sleep(2 * time.Second)
	}
}

var ch1 chan int
var ch2 chan int
var chs = []chan int{ch1, ch2}
var numbers = []int{1, 2, 3, 4, 5}

func TttChan() {
	persionChan := make(chan Person, 1)

	p1 := Person{"Harry", 32, Addr{"Shanxi", "Xian"}}
	fmt.Printf("P1 (1): %v\n", p1)

	persionChan <- p1

	p1.Address.district = "shijingshan"
	fmt.Printf("P2 (2): %v\n", p1)

	p1_copy := <-persionChan
	fmt.Printf("p1_copy: %v\n", p1_copy)
}

func T_timeout() {
	t := time.NewTimer(2 * time.Second)

	now := time.Now()
	fmt.Printf("Now time : %v.\n", now)

	expire := <-t.C
	fmt.Printf("Expiration time: %v.\n", expire)
}

func processex() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	resultChan := make(chan Result, 1)

	req, err := http.NewRequest("GET", "http://www.google.com", nil)
	if err != nil {
		fmt.Println("http request failed, err:", err)
		return
	}

	go func() {
		resp, err := client.Do(req)

		pack := Result{r: resp, err: err}
		//将返回信息写入管道(正确或者错误的)
		resultChan <- pack
	}()

	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		er := <-resultChan
		fmt.Println("Timeout! err=", er.err)
	case res := <-resultChan:
		defer res.r.Body.Close()
		out, _ := ioutil.ReadAll(res.r.Body)
		fmt.Printf("Server Response: %d\n", len(out))
	}
}

func IoTest() {
	reader := strings.NewReader("Clear is better than clever")
	p := make([]byte, 4)

	for {
		n, err := reader.Read(p)
		if err != nil {
			if err == io.EOF {
				fmt.Println("EOF:", n)
				break
			}
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(n, string(p[:n]))
	}

}

func copyFile2(srcFile, destFile string) (int64, error) {
	file1, err := os.Open(srcFile)
	if err != nil {
		fmt.Println("srcfile failed! err=", err)
		return 0, err
	}
	file2, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		fmt.Println("dstfile failed! err=", err)
		return 0, err
	}
	defer file1.Close()
	defer file2.Close()

	return io.Copy(file2, file1)
}

func CpyFileFun() {
	copyFile2("tstfun", "cpyfile")
}

type Origin struct {
	a uint64
	b uint64
}
type WithPadding struct {
	a uint64
	_ [56]byte
	b uint64
	_ [56]byte
}

var num = 1000 * 1000

func OriginParallel() {
	var v Origin

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.a, 1)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.b, 1)
		}
		wg.Done()
	}()
	wg.Wait()
	_ = v.a + v.b
}
func WithPaddingParallel() {
	var v WithPadding

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.a, 1)
		}
		wg.Done()
	}()

	go func() {
		for i := 0; i < num; i++ {
			atomic.AddUint64(&v.b, 1)
		}
		wg.Done()
	}()

	wg.Wait()
	_ = v.a + v.b
}

func TstPaddingFun() {
	var b time.Time

	b = time.Now()
	OriginParallel()
	fmt.Printf("OriginParallel. Cost=%+v.\n", time.Now().Sub(b))

	b = time.Now()
	WithPaddingParallel()
	fmt.Printf("WithPaddingParallel. Cost=%+v.\n", time.Now().Sub(b))
}

var bufpool *sync.Pool

func init() {
	bufpool = &sync.Pool{}
	bufpool.New = func() interface{} {
		return make([]byte, 32*1024)
	}
}

func Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	buf := bufpool.Get().([]byte)
	defer bufpool.Put(buf)

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

type ObjCh struct {
	Idx int
}

var readyObjs []*ObjCh

func TstObjCh() {

	fmt.Println("len(objs)=", len(readyObjs))

	for i := 0; i < 10; i++ {
		obj := &ObjCh{
			Idx: i,
		}
		readyObjs = append(readyObjs, obj)
	}

	for _, tmpobj := range readyObjs {
		fmt.Printf("tmpobj=%+v\n", *tmpobj)
	}

	tmpObjs := readyObjs

	fmt.Printf("  tmpObjs=%p\n", tmpObjs)
	fmt.Printf("readyObjs=%p\n", readyObjs)

	n := len(tmpObjs) - 1

	var ch *ObjCh
	ch = readyObjs[n]
	readyObjs[n] = nil
	readyObjs = readyObjs[:n]

	for _, tmpobj := range readyObjs {
		fmt.Printf("tmpobj=%+v\n", *tmpobj)
	}
	fmt.Printf("  ch=%+v\n", *ch)
}

func main() {
	TstObjCh()
	return

	TstPaddingFun()
	return

	CpyFileFun()
	return

	TstReaderEntry()
	return

	processex()
	return
	TttChan()
	return

	TstCtxEntry()
	return

	TstRoundTripEntry()
	return
	Tt1()
	return

	BitOpr()
	return
	TstWorkPoolFun()
	return

	TstGoLog()
	return

	TstPanic()
	return

	JustTest()
	return
	tst_map_fun()
	return

	Tstlocalipclient()
	return

	TstTimeInterval()
	return
	tst_json()
	return

	// TstPostCheck()
	// return
	TstYunPanClient()
	return

	{
		// var slice_ []int = make([]int, 0, 10)
		// fmt.Println(slice_)

		// var slice_1 []int = make([]int, 3)
		// fmt.Println(slice_1)

		// var slice_2 []int = []int{1, 2}
		// fmt.Println(slice_2)

		// return
	}
	// GetFileEEEE()
	// return

	// AddFile()
	// return

	TstSliceEntry()
	return

	TstReflectEntry()
	return

	TstChanEntry()
	return

	// TstChanEntry()
	// return

	TstCtx()
	return

	TstThreadId()
	return

	TstTick()
	return

	TstCloseCh()
	return

	Tst_MRecver_MSender()
	return

	fmt.Println("GetGoroutineIDStr=", GetGoroutineIDStr())
	return
	// Wrrap()
	// return

	Tst_MRecver_MSender()
	return

	Tst1Recver_NSender()
	return

	TstChByte()
	return

	TstDefer()
	return
	// {
	// 	var mapUse sync.Map
	// 	mapUse.Store(1, 1)

	// 	//Load 方法，获得value
	// 	if v, ok := mapUse.Load(1); ok {
	// 		fmt.Println("v=", v)
	// 		mapUse.Store(1, v.(int)+1)

	// 		vt, _ := mapUse.Load(1)
	// 		fmt.Println("vt=", vt)
	// 	} else {

	// 	}
	// 	return
	// }
	TstTcpConnPool()
	return

	TstDinner()
	return

	tst15()
	return

	StartRecvUpload()
	return

	TstChanEntry()
	return

	TstBlg4Fun()

	time.Sleep(2 * time.Second)
	return

	tst_fun_entry()
	return

	var obj T
	p1 := &obj

	// p1.A()
	// p1.B()

	_ = *identity(p1)
	_ = *ref(obj)

}

// {
// var fileBean OneFile
// fileBean.UserId = 1234567
// fileBean.OwnerId = fileBean.UserId

// fileBean.FileName = upFileName
// fileBean.FileType = 1

// fileBean.FileSuffix = "log"
// fileBean.FileSize = fileSize
// fileBean.Md5Hash = "asdfasdfadfasdaaa=="
// fileBean.Sha2Hash = "tsasdf123901234aaaa=="
// fileBean.SecretKey = "1234456"
// fileBean.Md5Hash = "123"
// fileBean.Sha2Hash = "456"

// retBuffer, _ := ffjson.Marshal(&fileBean)
// fmt.Println(" upload FileJson =", string(retBuffer))
// }
