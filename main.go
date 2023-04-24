package main

import (
	"strings"

	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm"
	"github.com/tetratelabs/proxy-wasm-go-sdk/proxywasm/types"
	"github.com/tidwall/gjson"
)

func main() {
	//json := `{"basic_auth_rules":[{"credentials":["ok:test","YWRtaW4zOmFkbWluMw=="],"prefix":"/productpage","request_methods":["GET","POST"]}]}`
	//fmt.Print(gjson.Valid(json))

	proxywasm.SetVMContext(&vmContext{})

}

type vmContext struct {
	// Embed the default VM context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultVMContext
}

// Override types.DefaultVMContext.
func (*vmContext) NewPluginContext(contextID uint32) types.PluginContext {
	return &pluginContext{}
}

type pluginContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultPluginContext
	headerName  string
	headerValue string
}

// 获取配置
func (p *pluginContext) OnPluginStart(pluginConfigurationSize int) types.OnPluginStartStatus {
	proxywasm.LogWarn("loading plugin config")
	rawData, err := proxywasm.GetPluginConfiguration()
	if rawData == nil {
		return types.OnPluginStartStatusOK
	}

	if err != nil {
		proxywasm.LogCriticalf("error reading plugin configuration: %v", err)
		return types.OnPluginStartStatusFailed
	}
	data := string(rawData)
	proxywasm.LogCriticalf("%s", data)

	if !gjson.Valid(data) {
		proxywasm.LogCritical(`**********************invalid configuration format; expected {"header": "<header name>", "value": "<header value>"}`)
		return types.OnPluginStartStatusFailed
	}

	p.headerName = strings.TrimSpace(gjson.Get(data, "header").String())
	p.headerValue = strings.TrimSpace(gjson.Get(data, "values").String())

	if p.headerName == "" || p.headerValue == "" {
		proxywasm.LogCritical(`invalid configuration format; expected {"header": "<header name>", "value": "<header value>"}`)
		return types.OnPluginStartStatusFailed
	}

	proxywasm.LogInfof("header from config: %s = %s", p.headerName, p.headerValue)

	return types.OnPluginStartStatusOK
}

type httpContext struct {
	// Embed the default plugin context here,
	// so that we don't need to reimplement all the methods.
	types.DefaultHttpContext
	contextID uint32
}

func (p *pluginContext) NewHttpContext(contextID uint32) types.HttpContext {
	return &httpContext{
		contextID: contextID,
	}
}

func (ctx *httpContext) OnHttpRequestHeaders(numHeaders int, endOfStream bool) types.Action {
	//err := proxywasm.ReplaceHttpRequestHeader("test", "best")
	//if err != nil {
	//	proxywasm.LogCritical("failed to set request header: test")
	//}

	hs, err := proxywasm.GetHttpRequestHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get request headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogCriticalf("request header --> %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

func (ctx *httpContext) OnHttpResponseHeaders(_ int, _ bool) types.Action {

	// Get and log the headers
	hs, err := proxywasm.GetHttpResponseHeaders()
	if err != nil {
		proxywasm.LogCriticalf("failed to get response headers: %v", err)
	}

	for _, h := range hs {
		proxywasm.LogCriticalf("response header <-- %s: %s", h[0], h[1])
	}
	return types.ActionContinue
}

// Override types.DefaultHttpContext.
func (ctx *httpContext) OnHttpStreamDone() {
	proxywasm.LogCriticalf("%d finished", ctx.contextID)
}
