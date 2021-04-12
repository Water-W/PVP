package metrics

import "io"

type Source interface {
	Source() io.Reader
}
