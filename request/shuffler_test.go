package request

import (
	"net/http"
	"testing"
)

type testGenerator struct{}

func (t *testGenerator) Generate(limit int) string {
	return "1234"
}

func Test_newShuffler(t *testing.T) {
	r, _ := http.NewRequest("GET", "https://test.com/{{1001:}}/abc", nil)
	s := newShuffler(r, func(match string) generator {
		return &testGenerator{}
	})

	s.Shuffle(r, 10)

	expectedURL := "https://test.com/1011/abc"
	if r.URL.String() != expectedURL {
		t.Errorf("Generated request URL does not match: %s, expected: %s", r.URL, expectedURL)
	}
}
