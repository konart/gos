package lib

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const K = 5

func TestGetData(t *testing.T) {
	testCases := []string{
		"",
		"Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.",
	}
	var i int
	var tc string

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if r.Method != "GET" {
			t.Errorf("exptected GET method, got %s instead", r.Method)
		}

		w.Write([]byte(testCases[i]))
	}))
	defer ts.Close()

	for i, tc = range testCases {
		w := NewWorker(K)
		data, err := w.getData(ts.URL)
		if err != nil {
			t.Errorf("Case #%d: %s", i, err)
		}
		str := string(data)
		if str != tc {
			t.Errorf("Case #%d, expected: \"%s\",\n got: \"%s\"", i, tc, str)
		}
	}
}

func TestCounting(t *testing.T) {
	testCases := []struct {
		body   string
		expect int
	}{
		{"", 0},
		{"Go is an open source programming language that makes it easy to build simple, reliable, and efficient software.", 1},
		{"Build version go1.10.3.", 0},
		{"Go, go, Inspector Gadget!", 1},
	}

	w := NewWorker(K)
	for i, tc := range testCases {
		result := w.getCount([]byte(tc.body), "Go")
		if result != tc.expect {
			t.Errorf("Case #%d, expected: %d, got %d", i, tc.expect, result)
		}
	}
}
