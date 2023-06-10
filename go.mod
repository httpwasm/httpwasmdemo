module github.com/httpwasm/httpwasmdemo

go 1.20

require (
	github.com/httpwasm/http-wasm-guest-tinygo v0.2.0
	github.com/httpwasm/http-wasm-host-go v0.5.1
)

replace github.com/httpwasm/http-wasm-guest-tinygo => ../http-wasm-guest-tinygo

replace github.com/httpwasm/http-wasm-host-go => ../http-wasm-host-go

require github.com/tetratelabs/wazero v1.2.0 // indirect
