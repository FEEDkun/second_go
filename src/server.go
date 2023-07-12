package main

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
)

type code struct {
	Message string `json:"message"`
}

type hello struct {
	Message1 string `json:"newtext"`
	Message2 string `json:"newtext2"`
}

const keyServerAddr = "serverAddr"

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Printf("%s: got / request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "This is my website!\n")
}

func getHello(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	len := r.ContentLength //不确定长度，直接用length方法获取
	by := make([]byte, len)

	fmt.Println(r.URL)
	fmt.Println(r.URL.Path, " -----1----- ", r.URL.RawQuery, "  ------2------  ", r.URL.RawFragment, " -------3------- ")
	fmt.Println(r.RequestURI)
	fmt.Println(r.Method)
	fmt.Println("Header:", r.Header)

	for s, strings := range r.Header { //获取请求头各个键值对
		fmt.Println("%s:%s", s, strings)
	}
	r.Body.Read(by)
	fmt.Println("请求body是：", string(by))

	//---------------------------------响应-----------------------------------------
	fmt.Printf("%s: got /hello request\n", ctx.Value(keyServerAddr))
	io.WriteString(w, "Hello, HTTP!\n")

	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	r.Body.Read(by)   //读取请求体内容到by字节数组
	str := string(by) //转字符串

	key := hello{ //json赋值
		Message1: str,
		Message2: str + str}
	val, _ := json.Marshal(key) //转json
	w.Write(val)

	fmt.Println()

	fmt.Println("响应头的Accept为：", w.Header().Get("Accept"))
	fmt.Println("新增响应项目为", w.Header().Get("Content-Type"))

}

func getfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	openFile, err := os.OpenFile("D:/Goproject/SAG/src/server.go", os.O_RDWR|os.O_CREATE, 777)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		io.WriteString(w, "文件读取失败！\n")
		return
	}

	reader := bufio.NewReader(openFile)
	for {
		s, e := reader.ReadString('\n')

		//fmt.Println(s)
		//io.WriteString(w, s)
		w.Header().Set("Content-Type", "application/json")

		key := code{s} //转json
		val, _ := json.Marshal(key)
		w.Write(val)

		if e == io.EOF {
			break
		}
		if e != nil {
			fmt.Println(e)
		}
	}
	openFile.Close()

	fmt.Printf("%s: got /fileout request\n", ctx.Value(keyServerAddr))

}

func main() {
	//name := [5]int{1,2,3,4,5}
	//str := make(chan string, 10)

	mux1 := http.NewServeMux()
	mux1.HandleFunc("/", getRoot)
	mux1.HandleFunc("/hello", getHello)
	mux1.HandleFunc("/fileout", getfile)

	//err := http.ListenAndServe(":3333", mux)

	ctx, cancelCtx := context.WithCancel(context.Background())
	serverOne := &http.Server{
		Addr:    ":8000",
		Handler: mux1,
		BaseContext: func(l net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, l.Addr().String())
			return ctx
		},
	}
	go func() {
		err := serverOne.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server one closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server one: %s\n", err)
		}
		cancelCtx()
	}()
	<-ctx.Done()

	servertwo := &http.Server{
		Addr:    ":3333",
		Handler: mux1,
		BaseContext: func(listener net.Listener) context.Context {
			ctx = context.WithValue(ctx, keyServerAddr, listener.Addr().String())
			return ctx
		},
	}
	go func() {
		err := servertwo.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server two closed\n")
		} else if err != nil {
			fmt.Printf("error listening for server two: %s\n", err)
		}
		cancelCtx()
	}()

}
