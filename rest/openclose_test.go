package rest_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	gen "github.com/massive-com/client-go/v3/rest/gen"
)

// Regression test for https://github.com/massive-com/client-go/issues/571.
// The daily open/close endpoint returns decimal preMarket and afterHours
// values, so they must unmarshal from JSON numbers with fractional parts.
func TestGetStocksOpenCloseDecimalFields(t *testing.T) {
	body := `{"status":"OK","from":"2024-01-02","symbol":"AAPL","open":1.5,"high":2.5,"low":1.0,"close":2.0,"volume":3,"preMarket":406.03,"afterHours":187.55}`
	rsp := &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}

	parsed, err := gen.ParseGetStocksOpenCloseResponse(rsp)
	if err != nil {
		t.Fatalf("ParseGetStocksOpenCloseResponse: %v", err)
	}
	if parsed.JSON200 == nil {
		t.Fatal("expected JSON200 to be populated")
	}
	if parsed.JSON200.PreMarket == nil || *parsed.JSON200.PreMarket != 406.03 {
		t.Errorf("preMarket: got %v, want 406.03", parsed.JSON200.PreMarket)
	}
	if parsed.JSON200.AfterHours == nil || *parsed.JSON200.AfterHours != 187.55 {
		t.Errorf("afterHours: got %v, want 187.55", parsed.JSON200.AfterHours)
	}
}
