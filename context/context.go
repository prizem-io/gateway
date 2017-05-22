package context

import (
	"bufio"
	"io"
	"mime/multipart"
	"net"
	"time"

	"github.com/prizem-io/gateway/config"
	"github.com/prizem-io/gateway/identity"
	//"golang.org/x/text/language"
)

type (
	MiddlewareHandler interface {
		Execute(ctx Context) error
		Next(ctx Context) error
		Stop()
		IsStopped() bool
	}

	StreamWriter func(w *bufio.Writer)

	Request interface {
		Param(string) string
		ParamInt(string) (int, error)
		VisitParams(func(key, value string))
		URLParam(string) string
		URLParamInt(string) (int, error)
		URLParams(func(key, value string))
		Method() string
		MethodBytes() []byte
		SetMethod(string)
		SetMethodBytes([]byte)
		Host() string
		HostBytes() []byte
		Path() string
		PathBytes() []byte
		RequestPath(bool) string
		RequestIP() string
		RemoteAddr() string
		SetAccepts(string)
		SetAcceptsBytes([]byte)
		SetContentType(string)
		SetContentTypeBytes(contentType []byte)
		FormValue(string) string
		FormValueBytes(string) []byte
		PostFormValue(string) string
		PostFormMulti(string) []string
		MultipartForm() (*multipart.Form, error)
		IsDelete() bool
		IsGet() bool
		IsHead() bool
		IsPost() bool
		IsPut() bool
		IsTLS() bool
		LocalAddr() net.Addr

		Body() []byte
		BodyGunzip() ([]byte, error)
		BodyInflate() ([]byte, error)
		AppendBody(string)
		AppendBodyBytes([]byte)
		SetBody([]byte)
		SetBodyString(string)
		SetBodyStream(bodyStream io.Reader, bodySize int)
		SetBodyStreamWriter(sw StreamWriter)
		SetConnectionClose()

		BodyWriteTo(w io.Writer) error
		BodyWriter() io.Writer
		ConnectionClose() bool
		MayContinue() bool
		Read(r *bufio.Reader) error
		ReadLimitBody(r *bufio.Reader, maxBodySize int) error

		SetHost(string)
		SetHostBytes([]byte)
		SetRequestURI(requestURI string)
		SetRequestURIBytes(requestURI []byte)
		// URI() *URI

		Headers(func(key, value string))
		HeadersBytes(func(key, value []byte))
		Header(string) string
		HeaderBytesK(k []byte) string
		HeaderBytesKV(k []byte) []byte
		HeaderBytesV(k string) []byte
		DeleteHeader(string)
		AddHeader(string, string)
		AddHeaderBytesK([]byte, string)
		AddHeaderBytesKV([]byte, []byte)
		AddHeaderBytesV(string, []byte)
		SetHeader(string, string)
		SetHeaderBytesK([]byte, string)
		SetHeaderBytesKV([]byte, []byte)
		SetHeaderBytesV(string, []byte)
	}

	Response interface {
		Body() []byte
		BodyGunzip() ([]byte, error)
		BodyInflate() ([]byte, error)
		IsBodyStream() bool
		AppendBody(string)
		AppendBodyBytes([]byte)
		SetBody(body []byte)
		SetBodyStream(bodyStream io.Reader, bodySize int)
		SetBodyString(body string)
		ConnectionClose() bool
		Reset()
		ResetBody()
		SetConnectionClose()
		SetStatusCode(int)
		SetContentLength(contentLength int)
		SetContentRange(startPos, endPos, contentLength int)
		SetLastModified(t time.Time)
		SetServer(server string)
		SetServerBytes(server []byte)
		SetContentType(string)
		SetContentTypeBytes([]byte)
		Write(w *bufio.Writer) error
		WriteDeflate(w *bufio.Writer) error
		WriteDeflateLevel(w *bufio.Writer, level int) error
		WriteGzip(w *bufio.Writer) error
		WriteGzipLevel(w *bufio.Writer, level int) error
		WriteTo(w io.Writer) (int64, error)

		Headers(func(key, value string))
		HeadersBytes(func(key, value []byte))
		Header(string) string
		HeaderBytesK(k []byte) string
		HeaderBytesKV(k []byte) []byte
		HeaderBytesV(k string) []byte
		DeleteHeader(string)
		AddHeader(string, string)
		AddHeaderBytesK([]byte, string)
		AddHeaderBytesKV([]byte, []byte)
		AddHeaderBytesV(string, []byte)
		SetHeader(string, string)
		SetHeaderBytesK([]byte, string)
		SetHeaderBytesKV([]byte, []byte)
		SetHeaderBytesV(string, []byte)
	}

	DataAccessor interface {
		GetPlugin(name string) (*config.Plugin, error)
		GetConsumer(id string) (*config.Consumer, error)
		GetCredential(id string) (interface{}, error)
		FindCredential(credentialType, clientID string) (interface{}, error)
		GetPlan(id string) (*config.Plan, error)
		GetPermission(id string) (*config.Permission, error)
	}

	GatewayContext interface {
		Locale() []string //language.Tag
		RequestID() string
		SubjectType() string
		GetDataAccessor() DataAccessor
		SetDataAccessor(DataAccessor)
		Credential() *config.Credential
		SetCredential(*config.Credential)
		Identity() identity.Identity
		SetIdentity(identity.Identity)
		Consumer() *config.Consumer
		SetConsumer(*config.Consumer)
		Plan() *config.Plan
		SetPlan(*config.Plan)
		Service() *config.Service
		SetService(*config.Service)
		Operation() *config.Operation
		SetOperation(*config.Operation)
		Version() string
		SetVersion(string)
		Claims() identity.Claims
		Get(string) interface{}
		GetString(string) string
		GetInt(string) int
		Values(func(string, interface{}))
		ValueMap() map[string]interface{}
		Set(string, interface{})
		Log(string, ...interface{})
	}

	Framework interface {
		SetMiddlewareHandler(handler MiddlewareHandler)
		Execute() error
		Next() error
		Stop()
		SetError(error)
		GetError() error
		IsStopped() bool
	}

	Context interface {
		Rq() Request
		Rs() Response
		DataAccessor
		GatewayContext
		Framework
		SendEntity(data interface{}) error
	}
)
