package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/httpwasm/http-wasm-host-go/api"
	"github.com/httpwasm/http-wasm-host-go/handler"
	wasmhandler "github.com/httpwasm/http-wasm-host-go/handler/nethttp"
)

var wasmHello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Println("go 宿主程序输出修改后的uri:" + r.RequestURI)
})

func main() {
	ctx := context.Background()
	wasmfile, err := os.ReadFile("./router/main.wasm")
	if err != nil {
		log.Panicln(err)
	}
	mw, err := wasmhandler.NewMiddleware(ctx, wasmfile, handler.Logger(api.ConsoleLogger{}))
	if err != nil {
		log.Panicln(err)
	}
	defer mw.Close(ctx)

	// Wrap the real handler with an interceptor implemented in WebAssembly.
	wasmWrapped := mw.NewHandler(ctx, wasmHello)

	//http.Handle("/hello", wasmWrapped)
	http.Handle("/wasm", wasmWrapped)
	http.ListenAndServe(":8000", nil)
}
