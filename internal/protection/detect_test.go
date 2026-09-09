package protection

import (
	"crypto/rand"
	"testing"
)

func TestShannonRange(t *testing.T) {
	jsonVal := []byte(`{"user":"alice","status":"active","city":"berlin","ok":true}`)
	jh := shannonEntropy(jsonVal)
	enc := make([]byte, 256)
	if _, err := rand.Read(enc); err != nil {
		t.Fatal(err)
	}
	eh := shannonEntropy(enc)
	t.Logf("json=%v enc=%v", jh, eh)
	if eh < 7.0 {
		t.Fatalf("random entropy too low: %v", eh)
	}
	if jh > 6.5 {
		t.Fatalf("json entropy too high: %v", jh)
	}
}
