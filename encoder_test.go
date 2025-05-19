package bone

import (
	"bytes"
	"testing"
)

func TestEncodeIllegal(t *testing.T) {
	t.Skip("values that should cause the encoder to error")
}

func FuzzEncode(f *testing.F) {
	for _, data := range DecodableSeedCorpus {
		f.Add(data)
	}
	f.Fuzz(func(t *testing.T, payload []byte) {
		values, err := Decode(payload)
		if err != nil {
			return
		}
		if !bytes.Equal(Encode(values), payload) {
			t.Errorf("re-encoded bytes do not match the original payload")
		}
	})
}
