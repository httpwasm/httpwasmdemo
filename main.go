package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/httpwasm/http-wasm-host-go/api"
	apihandler "github.com/httpwasm/http-wasm-host-go/api/handler"
	"github.com/httpwasm/http-wasm-host-go/handler"
)

type host struct {
	apihandler.UnimplementedHost
}

// GetURI implements the same method as documented on handler.Host.
func (h host) GetURI(ctx context.Context) string {
	rs := ctx.Value(requestStateKey{}).(*requestState)
	return rs.r.URL.Path
}

// SetURI implements the same method as documented on handler.Host.
func (h host) SetURI(ctx context.Context, uri string) {
	rs := ctx.Value(requestStateKey{}).(*requestState)
	rs.r.RequestURI = uri
}

// 使用全局Middleware
var mw handler.Middleware

// 初始化
func init() {
	ctx := context.Background()
	wasmfile, err := os.ReadFile("./router/main.wasm")
	if err != nil {
		log.Panicln(err)
	}
	mw, err = handler.NewMiddleware(ctx, wasmfile, host{}, handler.Logger(api.ConsoleLogger{}))
	if err != nil {
		log.Panicln(err)
	}
}

var wasmHello = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rs := &requestState{w: w, r: r}
	rs.enableFeatures(mw.Features())
	ctx = context.WithValue(ctx, requestStateKey{}, rs)
	outCtx, ctxNext, err := mw.HandleRequest(ctx)
	fmt.Println(outCtx, ctxNext, err)
	fmt.Println("go 宿主程序输出修改后的uri:" + r.RequestURI)
})

func main() {

	// Wrap the real handler with an interceptor implemented in WebAssembly.
	//wasmWrapped := mw.NewHandler(ctx, wasmHello)

	http.Handle("/", wasmHello)
	//http.Handle("/wasm", wasmWrapped)
	http.ListenAndServe(":8000", nil)
}
