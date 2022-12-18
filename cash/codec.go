package cash

import (
	"bytes"
	"compress/flate"
	"io"
	"sync"
)

type DeflateCodec struct {
	CompressionLevel int

	once sync.Once
}

func (codec *DeflateCodec) init() {
	codec.once.Do(func() {
		if codec.CompressionLevel == 0 {
			codec.CompressionLevel = flate.DefaultCompression
		}
	})
}

func (codec *DeflateCodec) Encode(b []byte) ([]byte, error) {
	codec.init()

	buf := new(bytes.Buffer)
	wc, err := flate.NewWriter(buf, codec.CompressionLevel)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(wc, bytes.NewReader(b)); err != nil {
		wc.Close()
		return nil, err
	}
	if err := wc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (codec *DeflateCodec) Decode(b []byte) ([]byte, error) {
	codec.init()

	rc := flate.NewReader(bytes.NewReader(b))
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, rc); err != nil {
		rc.Close()
		return nil, err
	}
	if err := rc.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
