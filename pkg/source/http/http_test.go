package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrReader(t *testing.T) {
	expBin := []byte("ok")
	r := bytes.NewReader(expBin)

	t.Run("not ok", func(t *testing.T) {
		expErr := fmt.Errorf("not ok")
		er := &ErrReader{
			Reader: r,
			Err:    expErr,
		}
		bin, actErr := ioutil.ReadAll(er)
		assert.Equal(t, expErr, actErr)
		assert.Empty(t, bin)
	})
	t.Run("ok", func(t *testing.T) {
		er := &ErrReader{
			Reader: r,
			Err:    nil,
		}
		actBin, actErr := ioutil.ReadAll(er)
		assert.NoError(t, actErr)
		assert.Equal(t, expBin, actBin)
	})

}
