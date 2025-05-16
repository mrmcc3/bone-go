package bone

import (
	"bytes"
	"testing"
)

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
