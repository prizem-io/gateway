package fasthttp

import (
	"bufio"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/valyala/fasthttp"

	"github.com/prizem-io/gateway/context"
	"github.com/prizem-io/gateway/server"
	"github.com/prizem-io/gateway/utils"
)

type FastHttpRequest struct {
	*fasthttp.RequestCtx
}

type FastHttpResponse struct {
	*fasthttp.RequestCtx
}

type FastHttpContext struct {
	*fasthttp.RequestCtx
	rq FastHttpRequest
	rs FastHttpResponse
	context.Common
}

var _fasthttpContextPool = &sync.Pool{
	New: func() interface{} {
		return &FastHttpContext{}
	},
}

func AcquireFastHttpContext(ctx *fasthttp.RequestCtx, subjectType string) *FastHttpContext {
	rc := _fasthttpContextPool.Get().(*FastHttpContext)
	rc.RequestCtx = ctx
	rc.rq.RequestCtx = ctx
	rc.rs.RequestCtx = ctx
	rc.Common.Initialize(subjectType)

	return rc
}

func ReleaseFastHttpContext(context *FastHttpContext) {
	_fasthttpContextPool.Put(context)
}

func (c *FastHttpContext) Reset() {
	c.Common.Reset()
	c.RequestCtx = nil
	c.rq.RequestCtx = nil
	c.rs.RequestCtx = nil
}

func (c *FastHttpContext) Rq() context.Request {
	return &c.rq
}

func (c *FastHttpContext) Rs() context.Response {
	return &c.rs
}

func (ctx *FastHttpContext) SendEntity(data interface{}) error {
	return server.WriteEntity(ctx, data)
}

func (ctx *FastHttpContext) Locale() []string {
	return []string{"en", "US", ""}
}

func (ctx *FastHttpRequest) Param(key string) string {
	return ctx.RequestCtx.UserValue(key).(string)
}

func (ctx *FastHttpRequest) ParamInt(key string) (int, error) {
	val, err := strconv.Atoi(ctx.RequestCtx.UserValue(key).(string))
	return val, err
}

func (ctx *FastHttpRequest) VisitParams(f func(key, value string)) {
	ctx.RequestCtx.VisitUserValues(func(key []byte, value interface{}) {
		if _value, ok := value.(string); ok {
			f(string(key), _value)
		}
	})
}

// URLParam returns the get parameter from a request , if any
func (ctx *FastHttpRequest) URLParam(key string) string {
	return string(ctx.RequestCtx.Request.URI().QueryArgs().Peek(key))
}

// URLParams returns a map of a list of each url(query) parameter
func (ctx *FastHttpRequest) URLParams(f func(key, value string)) {
	ctx.RequestCtx.Request.URI().QueryArgs().VisitAll(func(key, value []byte) {
		f(string(key), string(value))
	})
}

// URLParamInt returns the get parameter int value from a request , if any
func (ctx *FastHttpRequest) URLParamInt(key string) (int, error) {
	return strconv.Atoi(ctx.URLParam(key))
}

// Method returns the HTTP Method
func (ctx *FastHttpRequest) Method() string {
	return utils.BytesToString(ctx.RequestCtx.Method())
}

// MethodBytes returns the HTTP Method
func (ctx *FastHttpRequest) MethodBytes() []byte {
	return ctx.RequestCtx.Method()
}

// SetMethod returns the HTTP Method
func (ctx *FastHttpRequest) SetMethod(method string) {
	ctx.RequestCtx.Request.Header.SetMethod(method)
}

// SetMethodBytes returns the HTTP Method
func (ctx *FastHttpRequest) SetMethodBytes(method []byte) {
	ctx.RequestCtx.Request.Header.SetMethodBytes(method)
}

// Host returns the Host of the request( the url as string )
func (ctx *FastHttpRequest) Host() string {
	return utils.BytesToString(ctx.RequestCtx.Host())
}

// HostString returns the Host of the request( the url as bytes )
func (ctx *FastHttpRequest) HostBytes() []byte {
	return ctx.RequestCtx.Host()
}

// Path returns the full escaped path as string
// for unescaped use: ctx.RequestCtx.RequestURI() or RequestPath(escape bool)
func (ctx *FastHttpRequest) Path() string {
	return ctx.RequestPath(true)
}

// PathBytes returns the full escaped path as bytes
// for unescaped use: ctx.RequestCtx.RequestURI() or RequestPath(escape bool)
func (ctx *FastHttpRequest) PathBytes() []byte {
	return ctx.RequestPathBytes(true)
}

// RequestPath returns the requested path
func (ctx *FastHttpRequest) RequestPath(escape bool) string {
	if escape {
		return utils.BytesToString(ctx.RequestCtx.Path())
	}
	return utils.BytesToString(ctx.RequestCtx.RequestURI())
}

// RequestPathBytes returns the requested path
func (ctx *FastHttpRequest) RequestPathBytes(escape bool) []byte {
	if escape {
		return ctx.RequestCtx.Path()
	}
	return ctx.RequestCtx.RequestURI()
}

// RequestIP gets just the Remote Address from the client.
func (ctx *FastHttpRequest) RequestIP() string {
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(ctx.RequestCtx.RemoteAddr().String())); err == nil {
		return ip
	}
	return ""
}

// RemoteAddr is like RequestIP but it checks for proxy servers also, tries to get the real client's request IP
func (ctx *FastHttpRequest) RemoteAddr() string {
	header := string(ctx.RequestCtx.Request.Header.Peek("X-Real-Ip"))
	realIP := strings.TrimSpace(header)
	if realIP != "" {
		return realIP
	}
	realIP = string(ctx.RequestCtx.Request.Header.Peek("X-Forwarded-For"))
	idx := strings.IndexByte(realIP, ',')
	if idx >= 0 {
		realIP = realIP[0:idx]
	}
	realIP = strings.TrimSpace(realIP)
	if realIP != "" {
		return realIP
	}
	return ctx.RequestIP()
}

// FormValue returns a single value from post request's data
func (ctx *FastHttpRequest) FormValue(name string) string {
	return utils.BytesToString(ctx.RequestCtx.FormValue(name))
}

func (ctx *FastHttpRequest) FormValueBytes(name string) []byte {
	return ctx.RequestCtx.FormValue(name)
}

// PostFormValue returns a single value from post request's data
func (ctx *FastHttpRequest) PostFormValue(name string) string {
	return string(ctx.RequestCtx.PostArgs().Peek(name))
}

// PostFormMulti returns a slice of string from post request's data
func (ctx *FastHttpRequest) PostFormMulti(name string) []string {
	arrBytes := ctx.PostArgs().PeekMulti(name)
	arrStr := make([]string, len(arrBytes))
	for i, v := range arrBytes {
		arrStr[i] = string(v)
	}
	return arrStr
}

// SetAccepts sets the response writer's header key 'Accepts' to a given value(s)
func (ctx *FastHttpRequest) SetAccepts(s string) {
	ctx.RequestCtx.Request.Header.Set(server.HeaderAccepts, s)
}

// SetAcceptsBytes sets the response writer's header key 'Accepts' to a given value(s)
func (ctx *FastHttpRequest) SetAcceptsBytes(s []byte) {
	ctx.RequestCtx.Request.Header.SetBytesV(server.HeaderAccepts, s)
}

// SetContentType sets the response writer's header key 'Content-Type' to a given value(s)
func (ctx *FastHttpRequest) SetContentType(s string) {
	ctx.RequestCtx.Request.Header.Set(server.HeaderContentType, s)
}

// SetContentTypeBytes sets the response writer's header key 'Content-Type' to a given value(s)
func (ctx *FastHttpRequest) SetContentTypeBytes(s []byte) {
	ctx.RequestCtx.Request.Header.SetBytesV(server.HeaderContentType, s)
}

func (ctx *FastHttpRequest) Headers(f func(key, value string)) {
	ctx.RequestCtx.Request.Header.VisitAll(func(key, value []byte) {
		f(utils.BytesToString(key), utils.BytesToString(value))
	})
}

func (ctx *FastHttpRequest) HeadersBytes(f func(key, value []byte)) {
	ctx.RequestCtx.Request.Header.VisitAll(func(key, value []byte) {
		f(key, value)
	})
}

// Header returns the request header's value
// accepts one parameter, the key of the header (string)
// returns string
func (ctx *FastHttpRequest) Header(k string) string {
	return utils.BytesToString(ctx.RequestCtx.Request.Header.Peek(k))
}

func (ctx *FastHttpRequest) HeaderBytesK(k []byte) string {
	return utils.BytesToString(ctx.RequestCtx.Request.Header.PeekBytes(k))
}

func (ctx *FastHttpRequest) HeaderBytesKV(k []byte) []byte {
	return ctx.RequestCtx.Request.Header.PeekBytes(k)
}

func (ctx *FastHttpRequest) HeaderBytesV(k string) []byte {
	return ctx.RequestCtx.Request.Header.Peek(k)
}

func (ctx *FastHttpRequest) DeleteHeader(k string) {
	ctx.RequestCtx.Request.Header.Del(k)
}

func (ctx *FastHttpRequest) AddHeader(k string, v string) {
	ctx.RequestCtx.Request.Header.Add(k, v)
}

func (ctx *FastHttpRequest) AddHeaderBytesK(k []byte, v string) {
	ctx.RequestCtx.Request.Header.AddBytesK(k, v)
}

func (ctx *FastHttpRequest) AddHeaderBytesKV(k, v []byte) {
	ctx.RequestCtx.Request.Header.AddBytesKV(k, v)
}

func (ctx *FastHttpRequest) AddHeaderBytesV(k string, v []byte) {
	ctx.RequestCtx.Request.Header.AddBytesV(k, v)
}

// SetHeader write to the response writer's header to a given key the given value(s)
//
// Note: If you want to send a multi-line string as header's value use: strings.TrimSpace first.
func (ctx *FastHttpRequest) SetHeader(k string, v string) {
	//v = strings.TrimSpace(v)
	ctx.RequestCtx.Request.Header.Set(k, v)
}

func (ctx *FastHttpRequest) SetHeaderBytesK(k []byte, v string) {
	ctx.RequestCtx.Request.Header.SetBytesK(k, v)
}

func (ctx *FastHttpRequest) SetHeaderBytesKV(k, v []byte) {
	ctx.RequestCtx.Request.Header.SetBytesKV(k, v)
}

func (ctx *FastHttpRequest) SetHeaderBytesV(k string, v []byte) {
	ctx.RequestCtx.Request.Header.SetBytesV(k, v)
}

// GetCookie returns cookie's value by it's name
// returns empty string if nothing was found
func (ctx *FastHttpRequest) GetCookie(name string) (val string) {
	bcookie := ctx.RequestCtx.Request.Header.Cookie(name)
	if bcookie != nil {
		val = string(bcookie)
	}
	return
}

// Body returns the request body
func (ctx *FastHttpRequest) Body() []byte {
	return ctx.Request.Body()
}

func (ctx *FastHttpRequest) BodyGunzip() ([]byte, error) {
	return ctx.Request.BodyGunzip()
}

func (ctx *FastHttpRequest) BodyInflate() ([]byte, error) {
	return ctx.Request.BodyInflate()
}

func (ctx *FastHttpRequest) AppendBody(body string) {
	ctx.Request.AppendBodyString(body)
}
func (ctx *FastHttpRequest) AppendBodyBytes(body []byte) {
	ctx.Request.AppendBody(body)
}

// Body returns the request body
func (ctx *FastHttpRequest) SetBody(body []byte) {
	ctx.Request.SetBody(body)
}

func (ctx *FastHttpRequest) SetBodyString(body string) {
	ctx.Request.SetBodyString(body)
}

func (ctx *FastHttpRequest) SetBodyStream(bodyStream io.Reader, bodySize int) {
	ctx.Request.SetBodyStream(bodyStream, bodySize)
}

func (ctx *FastHttpRequest) SetBodyStreamWriter(sw context.StreamWriter) {
	ctx.Request.SetBodyStreamWriter(fasthttp.StreamWriter(sw))
}

func (ctx *FastHttpRequest) SetConnectionClose() {
	ctx.Request.SetConnectionClose()
}

func (ctx *FastHttpRequest) BodyWriteTo(w io.Writer) error {
	return ctx.Request.BodyWriteTo(w)
}

func (ctx *FastHttpRequest) BodyWriter() io.Writer {
	return ctx.Request.BodyWriter()
}

func (ctx *FastHttpRequest) ConnectionClose() bool {
	return ctx.Request.ConnectionClose()
}

func (ctx *FastHttpRequest) MayContinue() bool {
	return ctx.Request.MayContinue()
}

func (ctx *FastHttpRequest) Read(r *bufio.Reader) error {
	return ctx.Request.Read(r)
}

func (ctx *FastHttpRequest) ReadLimitBody(r *bufio.Reader, maxBodySize int) error {
	return ctx.Request.ReadLimitBody(r, maxBodySize)
}

func (ctx *FastHttpRequest) SetHost(host string) {
	ctx.Request.SetHost(host)
}

func (ctx *FastHttpRequest) SetHostBytes(host []byte) {
	ctx.Request.SetHostBytes(host)
}

func (ctx *FastHttpRequest) SetRequestURI(requestURI string) {
	ctx.Request.SetRequestURI(requestURI)
}

func (ctx *FastHttpRequest) SetRequestURIBytes(requestURI []byte) {
	ctx.Request.SetRequestURIBytes(requestURI)
}

/* URI is fasthttp specific
func (ctx *FastHttpRequest) URI() *URI {
	return ctx.Request.URI()
}
*/

/////////////////
/////////////////
/////////////////

func (ctx *FastHttpResponse) Headers(f func(key, value string)) {
	ctx.RequestCtx.Response.Header.VisitAll(func(key, value []byte) {
		f(utils.BytesToString(key), utils.BytesToString(value))
	})
}

func (ctx *FastHttpResponse) HeadersBytes(f func(key, value []byte)) {
	ctx.RequestCtx.Response.Header.VisitAll(func(key, value []byte) {
		f(key, value)
	})
}

func (ctx *FastHttpResponse) Header(k string) string {
	return utils.BytesToString(ctx.RequestCtx.Response.Header.Peek(k))
}

func (ctx *FastHttpResponse) HeaderBytesK(k []byte) string {
	return utils.BytesToString(ctx.RequestCtx.Response.Header.PeekBytes(k))
}

func (ctx *FastHttpResponse) HeaderBytesKV(k []byte) []byte {
	return ctx.RequestCtx.Response.Header.PeekBytes(k)
}

func (ctx *FastHttpResponse) HeaderBytesV(k string) []byte {
	return ctx.RequestCtx.Response.Header.Peek(k)
}

func (ctx *FastHttpResponse) DeleteHeader(k string) {
	ctx.RequestCtx.Response.Header.Del(k)
}

func (ctx *FastHttpResponse) AddHeader(k string, v string) {
	ctx.RequestCtx.Response.Header.Add(k, v)
}

func (ctx *FastHttpResponse) AddHeaderBytesK(k []byte, v string) {
	ctx.RequestCtx.Response.Header.AddBytesK(k, v)
}

func (ctx *FastHttpResponse) AddHeaderBytesKV(k, v []byte) {
	ctx.RequestCtx.Response.Header.AddBytesKV(k, v)
}

func (ctx *FastHttpResponse) AddHeaderBytesV(k string, v []byte) {
	ctx.RequestCtx.Response.Header.AddBytesV(k, v)
}

// SetHeader write to the response writer's header to a given key the given value(s)
//
// Note: If you want to send a multi-line string as header's value use: strings.TrimSpace first.
func (ctx *FastHttpResponse) SetHeader(k string, v string) {
	ctx.RequestCtx.Response.Header.Set(k, v)
}

func (ctx *FastHttpResponse) SetHeaderBytesK(k []byte, v string) {
	ctx.RequestCtx.Response.Header.SetBytesK(k, v)
}

func (ctx *FastHttpResponse) SetHeaderBytesKV(k, v []byte) {
	ctx.RequestCtx.Response.Header.SetBytesKV(k, v)
}

func (ctx *FastHttpResponse) SetHeaderBytesV(k string, v []byte) {
	ctx.RequestCtx.Response.Header.SetBytesV(k, v)
}

// SetCookie adds a cookie
func (ctx *FastHttpResponse) SetCookie(cookie *fasthttp.Cookie) {
	ctx.RequestCtx.Response.Header.SetCookie(cookie)
}

// SetCookieKV adds a cookie, receives just a key(string) and a value(string)
func (ctx *FastHttpResponse) SetCookieKV(key, value string) {
	c := fasthttp.AcquireCookie() // &fasthttp.Cookie{}
	c.SetKey(key)
	c.SetValue(value)
	c.SetHTTPOnly(true)
	c.SetExpire(time.Now().Add(time.Duration(120) * time.Minute))
	ctx.SetCookie(c)
	fasthttp.ReleaseCookie(c)
}

// RemoveCookie deletes a cookie by it's name/key
func (ctx *FastHttpResponse) RemoveCookie(name string) {
	cookie := fasthttp.AcquireCookie()
	cookie.SetKey(name)
	cookie.SetValue("")
	cookie.SetPath("/")
	cookie.SetHTTPOnly(true)
	exp := time.Now().Add(-time.Duration(1) * time.Minute) //RFC says 1 second, but make sure 1 minute because we are using fasthttp
	cookie.SetExpire(exp)
	ctx.Response.Header.SetCookie(cookie)
	fasthttp.ReleaseCookie(cookie)
}

// Body returns the response body
func (ctx *FastHttpResponse) Body() []byte {
	return ctx.Response.Body()
}

func (ctx *FastHttpResponse) BodyGunzip() ([]byte, error) {
	return ctx.Response.BodyGunzip()
}

func (ctx *FastHttpResponse) BodyInflate() ([]byte, error) {
	return ctx.Response.BodyInflate()
}

func (ctx *FastHttpResponse) IsBodyStream() bool {
	return ctx.Response.IsBodyStream()
}

func (ctx *FastHttpResponse) AppendBody(body string) {
	ctx.Response.AppendBodyString(body)
}

func (ctx *FastHttpResponse) AppendBodyBytes(body []byte) {
	ctx.Response.AppendBody(body)
}

// Body returns the request body
func (ctx *FastHttpResponse) SetBody(body []byte) {
	ctx.Response.SetBody(body)
}

func (ctx *FastHttpResponse) SetBodyStream(bodyStream io.Reader, bodySize int) {
	ctx.Response.SetBodyStream(bodyStream, bodySize)
}

func (ctx *FastHttpResponse) SetBodyString(body string) {
	ctx.Response.SetBodyString(body)
}

func (ctx *FastHttpResponse) ConnectionClose() bool {
	return ctx.Response.ConnectionClose()
}

func (ctx *FastHttpResponse) Reset() {
	ctx.Response.Reset()
}

func (ctx *FastHttpResponse) ResetBody() {
	ctx.Response.ResetBody()
}

func (ctx *FastHttpResponse) SetConnectionClose() {
	ctx.Response.SetConnectionClose()
}

func (ctx *FastHttpResponse) SetStatusCode(status int) {
	ctx.Response.SetStatusCode(status)
}

func (ctx *FastHttpResponse) SetContentLength(contentLength int) {
	ctx.Response.Header.SetContentLength(contentLength)
}

func (ctx *FastHttpResponse) SetContentRange(startPos, endPos, contentLength int) {
	ctx.Response.Header.SetContentRange(startPos, endPos, contentLength)
}

func (ctx *FastHttpResponse) SetLastModified(t time.Time) {
	ctx.Response.Header.SetLastModified(t)
}

func (ctx *FastHttpResponse) SetServer(server string) {
	ctx.Response.Header.SetServer(server)
}

func (ctx *FastHttpResponse) SetServerBytes(server []byte) {
	ctx.Response.Header.SetServerBytes(server)
}

// SetContentType sets the response writer's header key 'Content-Type' to a given value(s)
func (ctx *FastHttpResponse) SetContentType(s string) {
	ctx.Response.Header.Set(server.HeaderContentType, s)
}

// SetContentTypeBytes sets the response writer's header key 'Content-Type' to a given value(s)
func (ctx *FastHttpResponse) SetContentTypeBytes(s []byte) {
	ctx.Response.Header.SetBytesV(server.HeaderContentType, s)
}

func (ctx *FastHttpResponse) Write(w *bufio.Writer) error {
	return ctx.Response.Write(w)
}

func (ctx *FastHttpResponse) WriteDeflate(w *bufio.Writer) error {
	return ctx.Response.WriteDeflate(w)
}

func (ctx *FastHttpResponse) WriteDeflateLevel(w *bufio.Writer, level int) error {
	return ctx.Response.WriteDeflateLevel(w, level)
}

func (ctx *FastHttpResponse) WriteGzip(w *bufio.Writer) error {
	return ctx.Response.WriteGzip(w)
}

func (ctx *FastHttpResponse) WriteGzipLevel(w *bufio.Writer, level int) error {
	return ctx.Response.WriteGzipLevel(w, level)
}

func (ctx *FastHttpResponse) WriteTo(w io.Writer) (int64, error) {
	return ctx.Response.WriteTo(w)
}

/////////////////
/////////////////
/////////////////

/* Storage */

// Get returns the user's value from a key
// if doesn't exists returns nil
func (ctx *FastHttpContext) Get(key string) interface{} {
	return ctx.RequestCtx.UserValue(key)
}

// GetFmt returns a value which has this format: func(format string, args ...interface{}) string
// if doesn't exists returns nil
func (ctx *FastHttpContext) GetFmt(key string) func(format string, args ...interface{}) string {
	if v, ok := ctx.Get(key).(func(format string, args ...interface{}) string); ok {
		return v
	}
	return func(format string, args ...interface{}) string { return "" }

}

// GetString same as Get but returns the value as string
// if nothing founds returns empty string ""
func (ctx *FastHttpContext) GetString(key string) string {
	if v, ok := ctx.Get(key).(string); ok {
		return v
	}

	return ""
}

// GetInt same as Get but returns the value as int
// if nothing founds returns -1
func (ctx *FastHttpContext) GetInt(key string) int {
	if v, ok := ctx.Get(key).(int); ok {
		return v
	}

	return -1
}

// UserData calls visitor for each user data value
func (ctx *FastHttpContext) Values(visitor func(string, interface{})) {
	ctx.RequestCtx.VisitUserValues(func(key []byte, value interface{}) {
		visitor(string(key), value)
	})
}

// UserDataMap returns a map of a user data values
func (ctx *FastHttpContext) ValueMap() map[string]interface{} {
	userdata := make(map[string]interface{})
	ctx.RequestCtx.VisitUserValues(func(key []byte, value interface{}) {
		userdata[string(key)] = value
	})
	return userdata
}

// Set sets a value to a key in the values map
func (ctx *FastHttpContext) Set(key string, value interface{}) {
	ctx.RequestCtx.SetUserValue(key, value)
}

// Log logs to the defined logger
func (ctx *FastHttpContext) Log(format string, a ...interface{}) {
	//ctx.framework.Logger.Printf(format, a...)
}

func (c *FastHttpContext) Execute() error {
	return c.Common.DoExecute(c)
}

func (c *FastHttpContext) Next() error {
	return c.Common.DoNext(c)
}
