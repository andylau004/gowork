package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

func dummy_reader() {
	fmt.Println("")
}

type alphaReader struct {
	// alphaReader 里组合了标准库的 io.Reader
	reader io.Reader
	// 资源
	// src string
	// 当前读取到的位置
	// cur int
}

func newAlphaReader(reader io.Reader) *alphaReader {
	return &alphaReader{reader: reader}
}

func alpha(r byte) byte {
	if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
		return r
	}
	return 0
}

// Read 方法
func (a *alphaReader) Read(p []byte) (int, error) {
	// fmt.Println("len(p)=", len(p))

	// 这行代码调用的就是 io.Reader
	n, err := a.reader.Read(p)
	if err != nil {
		return n, err
	}

	buf := make([]byte, n)
	for i := 0; i < n; i++ {
		if char := alpha(p[i]); char != 0 {
			buf[i] = char
		}
	}

	copy(p, buf)
	return n, nil
}

type chanWriter struct {
	// ch 实际上就是目标资源
	ch chan byte
}

func newChanWriter() *chanWriter {
	return &chanWriter{make(chan byte, 1024)}
}

func (w *chanWriter) Chan() <-chan byte {
	return w.ch
}

func (w *chanWriter) Write(p []byte) (int, error) {
	n := 0
	// 遍历输入数据，按字节写入目标资源
	for _, b := range p {
		w.ch <- b
		n++
	}
	return n, nil
}
func (w *chanWriter) Close() error {
	close(w.ch)
	return nil
}

func ChanWriterFun() {
	writer := newChanWriter()

	go func() {
		defer fmt.Println("func noname done!")
		defer writer.Close()
		writer.Write([]byte("Cyber Punk "))
		writer.Write([]byte("me!"))
	}()

	for c := range writer.Chan() {
		fmt.Printf("%c", c)
	}
	fmt.Println()

}

func TstWriterEntry() {
	{
		proverbs := new(bytes.Buffer)
		proverbs.WriteString("Channels orchestrate mutexes serialize\n")
		proverbs.WriteString("Cgo is not Go\n")
		proverbs.WriteString("Errors are values\n")
		proverbs.WriteString("Don't panic\n")

		file, err := os.Create("./proverbs.txt")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		// io.Copy 完成了从 proverbs 读取数据并写入 file 的流程
		if _, err := io.Copy(file, proverbs); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("file created")
		return
	}
	// {
	// 	ChanWriterFun()
	// 	return
	// }

	proverbs := []string{
		"Channels orchestrate mutexes serialize",
		"Cgo is not Go",
		"Errors are values",
		"Don't panic",
	}

	// var writer bytes.Buffer
	for _, p := range proverbs {
		// n, err := writer.Write([]byte(p))
		// 因为 os.Stdout 也实现了 io.Writer
		n, err := os.Stdout.Write([]byte(p))
		if err != nil {
			fmt.Println("write failed! err=", err)
			os.Exit(1)
		}
		if n != len(p) {
			fmt.Println("failed to write data")
			os.Exit(1)
		}
		// fmt.Println("save string=", writer.String())
		// fmt.Println("")
	}

}

func TstReaderEntry() {
	{
		TstWriterEntry()
		return
	}

	{
		file, err := os.Open("./tstlog.go")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		// 任何实现了 io.Reader 的类型都可以传入 newAlphaReader
		// 至于具体如何读取文件，那是标准库已经实现了的，我们不用再做一遍，达到了重用的目的
		reader := newAlphaReader(file)
		p := make([]byte, 5)
		for {
			n, err := reader.Read(p)
			if err == io.EOF {
				break
			}
			fmt.Println("p=", string(p[:n]))
		}
		fmt.Println()
		return
	}
	{
		//  使用实现了标准库 io.Reader 接口的 strings.Reader 作为实现
		reader := newAlphaReader(strings.NewReader("Hello! It's 9am, where is the sun?"))
		p := make([]byte, 5)
		for {
			n, err := reader.Read(p)
			if err == io.EOF {
				break
			}
			fmt.Println("p=", string(p[:n]))
		}
		fmt.Println()
	}

	// r := strings.NewReader("some io.Reader stream to be read\n")

	// if _, err := io.Copy(os.Stdout, r); err != nil {
	// 	// log.Fatal(err)
	// 	fmt.Println("err=", err)
	// }

}
