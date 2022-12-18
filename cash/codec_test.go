package cash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodec(t *testing.T) {
	codecs := []Codec{
		new(DeflateCodec),
	}
	for _, codec := range codecs {
		t.Run(fmt.Sprintf("%T", codec), func(t *testing.T) {
			cases := []string{
				"",
				"testing",
			}
			for _, c := range cases {
				b, err := codec.Encode([]byte(c))
				assert.NoError(t, err, "%T.Encode([]byte(%q))", codec, c)

				restored, err := codec.Decode(b)
				assert.NoError(t, err, "%T.Decode(%v)", codec, b)

				assert.Equal(t, []byte(c), restored, "%T failed decoding from its encoding bytes", codec)
			}
		})
	}
}
