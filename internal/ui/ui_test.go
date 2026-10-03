package ui_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/raffleberry/riptvtime/internal/ui"
)

// positive is true if the test is expected-to-pass
// if false, the test is expected-to-fail
func test_h(t *testing.T, positive bool, h http.Handler, filepath string, expected []string) {
	req := httptest.NewRequest(http.MethodGet, filepath, nil)
	resRec := httptest.NewRecorder()

	h.ServeHTTP(resRec, req)

	res := resRec.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status OK; got %v", res.Status)
	}

	got, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}
	t.Log(string(got))
	for i := range expected {
		if positive {
			if strings.Contains(string(got), expected[i]) {
				// good !
			} else {
				t.Errorf("expected response body to contain '%s'; got %v", expected[i], string(got))
			}
		} else {
			if strings.Contains(string(got), expected[i]) {
				t.Errorf("expected response body to not contain '%s'; got %v", expected[i], string(got))
			} else {
				// good !
			}
		}
	}
}

func Test_SpaHandler_BaseJs(t *testing.T) {
	h1 := ui.NewSpaHandler("static", "")
	test_h(t, true, h1, "/app/base.js", []string{`""`})

	h2 := ui.NewSpaHandler("static", "/rtt")
	test_h(t, true, h2, "/app/base.js", []string{`"/rtt"`})

}

func Test_SpaHandler_IndexHtml(t *testing.T) {
	h1 := ui.NewSpaHandler("static", "")
	test_h(t, true, h1, "/", []string{
		"html",
		"head",
		"title",
		"script",
		"body",
		"bootstrap",
		"vue",
		"main.js",

		`href="/libs`,
		`src="/libs`,
		`src="/app`,
	})

	h2 := ui.NewSpaHandler("static", "/rtt")
	test_h(t, true, h2, "/", []string{
		"html",
		"head",
		"title",
		"script",
		"body",
		"bootstrap",
		"vue",
		"main.js",

		`href="/rtt/libs`,
		`src="/rtt/libs`,
		`src="/rtt/app`,
	})
	test_h(t, false, h2, "/", []string{
		`href="/libs`,
		`src="/libs`,
		`src="/app`,
	})

}
