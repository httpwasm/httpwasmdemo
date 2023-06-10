package main

import (
	"io"

	"github.com/httpwasm/http-wasm-guest-tinygo/handler"
	"github.com/httpwasm/http-wasm-guest-tinygo/handler/api"
)

// tinygo build -o ./main.wasm -scheduler=none --no-debug -target=wasi ./main.go
func main() {
	handler.HandleRequestFn = handleRequest
	//handler.Host.LogEnabled(api.LogLevelDebug)
}

// handleRequest implements a simple HTTP router.
func handleRequest(req api.Request, resp api.Response) (next bool, reqCtx uint32) {
	handler.Host.Log(api.LogLevelError, "-------------wasm输出开始---------")
	handler.Host.Log(api.LogLevelError, "uri:"+req.GetURI())
	handler.Host.Log(api.LogLevelError, "请求方法:"+req.GetMethod())
	handler.Host.EnableFeatures(api.FeatureBufferRequest)
	w := bodyWriter{}
	req.Body().WriteTo(&w)
	handler.Host.Log(api.LogLevelError, "请求的body:"+string(w))
	handler.Host.Log(api.LogLevelError, "修改uri为/abc,让go宿主程序接收:")
	req.SetURI("/abc")

	handler.Host.Log(api.LogLevelError, "---------wasm输出结束---------")
	next = true // proceed to the next handler on the host.
	/*
		// If the URI starts with /host, trim it and dispatch to the next handler.
		if uri := req.GetURI(); strings.HasPrefix(uri, "/wasm") {
			req.SetURI(uri[5:])
			next = true // proceed to the next handler on the host.
			return
		}

		// Serve a static response
		//resp.Headers().Set("Content-Type", "text/plain")
		//resp.Body().WriteString("hello")
		//handler.Host.Log(api.LogLevelError, "wasm 调用 输出")
		return // skip the next handler, as we wrote a response.
	*/
	return
}

// compile-time check to ensure bodyWriter implements io.Writer.
var _ io.Writer = (*bodyWriter)(nil)

type bodyWriter []byte

// Write adds an extra newline prior to printing it to the wasi.
func (b *bodyWriter) Write(p []byte) (n int, err error) {
	*b = append(*b, p...)
	return len(*b), nil
}
