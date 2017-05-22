package http

import (
	"fmt"
	"strings"
	"time"

	"github.com/prizem-io/gateway/context"
	"github.com/prizem-io/gateway/utils"
	"github.com/valyala/fasthttp"
)

type HTTP struct{}

var client = &fasthttp.Client{}

var filteredRequestHeaders = map[string]struct{}{
	"content-length": {},
}

var filteredResponseHeaders = map[string]struct{}{
	"content-length": {},
}
var singleOccuranceHeaders = map[string]struct{}{
	"host":            {},
	"server":          {},
	"accept":          {},
	"content-type":    {},
	"user-agent":      {},
	"connection":      {},
	"cache-control":   {},
	"origin":          {},
	"accept-encoding": {},
	"accept-language": {},
}

func New() *HTTP {
	return &HTTP{}
}

func (r *HTTP) Name() string {
	return "http"
}

func (r *HTTP) String() string {
	return r.Name()
}

func (r *HTTP) Handle(ctx context.Context) error {
	rq := ctx.Rq()
	rs := ctx.Rs()
	s := ctx.Service()

	// TODO: Use discovery for hostname

	target := fmt.Sprintf("%s://%s%s", utils.StringDefault(s.Scheme, "http"), s.Hostnames[0], rq.Path())
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethodBytes(rq.MethodBytes())
	req.SetRequestURI(target)
	rq.Headers(func(key, value string) {
		keyLower := strings.ToLower(key)
		if _, ok := filteredRequestHeaders[keyLower]; !ok {
			if _, ok := singleOccuranceHeaders[keyLower]; ok {
				req.Header.Set(key, value)
			} else {
				req.Header.Add(key, value)
			}
		}
	})
	body := rq.Body()
	req.SetBody(body)

	err := client.DoTimeout(req, resp, time.Second)
	if err != nil {
		return err
	}

	rs.SetStatusCode(resp.StatusCode())
	rs.SetHeader("Server", "Prizem Gateway")
	resp.Header.VisitAll(func(key, value []byte) {
		keyStr := utils.BytesToString(key)
		keyLower := strings.ToLower(keyStr)
		if _, ok := filteredResponseHeaders[keyLower]; !ok {
			valueStr := utils.BytesToString(value)
			if _, ok := singleOccuranceHeaders[keyLower]; ok {
				rs.SetHeader(keyStr, valueStr)
			} else {
				rs.AddHeader(keyStr, valueStr)
			}
		}
	})
	rs.SetBody(resp.Body())

	return nil
}
