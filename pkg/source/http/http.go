package http

import (
	"io"
	"net/http"

	"github.com/Water-W/PVP/pkg/metrics"
)

var _ metrics.Source = (*Source)(nil)

type Source struct {
	url string
}

func NewSource(url string) *Source {
	return &Source{
		url: url,
	}
}

func (h *Source) Source() io.Reader {
	resp, err := http.Get(h.url) //using http.default client
	return &ErrReader{
		Reader: resp.Body,
		Err:    err,
	}
}

/*===========================================================================*/
// helper

// ErrReader wraps an error and a reader.
// ErrReader first checks the given error, if the error is not nil, the read method will
// return the given error; otherwise, the read process will be called.
type ErrReader struct {
	io.Reader
	Err error
}

func (r *ErrReader) Read(b []byte) (int, error) {
	if r.Err != nil {
		return 0, r.Err
	}
	return r.Reader.Read(b)
}
