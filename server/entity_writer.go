package server

import (
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/golang/gddo/httputil/header"
	"github.com/prizem-io/gateway/context"
)

// MIME types
const (
	JSON = "application/json"
	XML  = "text/xml"
)

type EntityWriter func(ctx context.Context, data interface{}) error

// DataWriters lists all supported content types and the corresponding data writers.
// By default, JSON, XML, and HTML are supported. You may modify this variable before calling TypeNegotiator
// to customize supported data writers.
var EntityWriters = map[string]EntityWriter{
	JSON: JSONEntityWriter,
	XML:  XMLEntityWriter,
}

var EntityWriterFormats = []string{JSON, XML}

func WriteEntity(ctx context.Context, data interface{}) error {
	format := NegotiateContentType(ctx, EntityWriterFormats, EntityWriterFormats[0])
	return EntityWriters[format](ctx, data)
}

func NegotiateContentType(ctx context.Context, offers []string, defaultOffer string) string {
	bestOffer := defaultOffer
	bestQ := -1.0
	bestWild := 3
	accept := ctx.Rq().Header("Accept")
	headers := http.Header{
		"Accept": []string{accept},
	}
	specs := header.ParseAccept(headers, "Accept")
	for _, offer := range offers {
		for _, spec := range specs {
			switch {
			case spec.Q == 0.0:
				// ignore
			case spec.Q < bestQ:
				// better match found
			case spec.Value == "*/*":
				if spec.Q > bestQ || bestWild > 2 {
					bestQ = spec.Q
					bestWild = 2
					bestOffer = offer
				}
			case strings.HasSuffix(spec.Value, "/*"):
				if strings.HasPrefix(offer, spec.Value[:len(spec.Value)-1]) &&
					(spec.Q > bestQ || bestWild > 1) {
					bestQ = spec.Q
					bestWild = 1
					bestOffer = offer
				}
			default:
				if spec.Value == offer &&
					(spec.Q > bestQ || bestWild > 0) {
					bestQ = spec.Q
					bestWild = 0
					bestOffer = offer
				}
			}
		}
	}
	return bestOffer
}

func JSONEntityWriter(ctx context.Context, data interface{}) error {
	var bytes []byte
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	rs := ctx.Rs()
	rs.SetContentType("application/json")
	rs.SetBody(bytes)
	return nil
}

func XMLEntityWriter(ctx context.Context, data interface{}) error {
	var bytes []byte
	bytes, err := xml.Marshal(data)
	if err != nil {
		return err
	}

	rs := ctx.Rs()
	rs.SetContentType("application/json")
	rs.SetBody(bytes)
	return nil
}
