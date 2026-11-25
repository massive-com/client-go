package rest_test

import (
	"reflect"
	"testing"
	"time"

	gen "github.com/massive-com/client-go/v3/rest/gen"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Regression test for https://github.com/massive-com/client-go/issues/579.
// The published_utc comparison params on ListNews were generated as anonymous
// `*struct{ union json.RawMessage }` fields, which cannot be constructed from
// outside the gen package. They must be a named type instead.
func TestListNewsPublishedUtcParamsAreNamedTypes(t *testing.T) {
	params := reflect.TypeOf(gen.ListNewsParams{})
	for _, fieldName := range []string{"PublishedUtc", "PublishedUtcGte", "PublishedUtcGt", "PublishedUtcLte", "PublishedUtcLt"} {
		field, ok := params.FieldByName(fieldName)
		if !ok {
			t.Fatalf("ListNewsParams.%s is missing", fieldName)
		}
		typ := field.Type.Elem()
		if typ.Name() == "" {
			t.Errorf("ListNewsParams.%s has anonymous type %s; callers cannot construct it", fieldName, field.Type)
		}
	}
}

// The named union type must accept both a time.Time (date-time) and an
// openapi_types.Date (date), and the request builder must serialize it into
// the expected published_utc.* query parameter.
func TestListNewsPublishedUtcFiltersRoundTrip(t *testing.T) {
	params := &gen.ListNewsParams{}

	dateTime := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	gte := gen.DateTimeOrDate{}
	if err := gte.FromDateTimeOrDate0(dateTime); err != nil {
		t.Fatalf("FromDateTimeOrDate0: %v", err)
	}
	params.PublishedUtcGte = &gte

	date := openapi_types.Date{Time: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC)}
	lt := gen.DateTimeOrDate{}
	if err := lt.FromDateTimeOrDate1(date); err != nil {
		t.Fatalf("FromDateTimeOrDate1: %v", err)
	}
	params.PublishedUtcLt = &lt

	req, err := gen.NewListNewsRequest("https://api.massive.com", params)
	if err != nil {
		t.Fatalf("NewListNewsRequest: %v", err)
	}

	query := req.URL.Query()
	if got := query.Get("published_utc.gte"); got == "" {
		t.Error("published_utc.gte query param is empty")
	}
	if got := query.Get("published_utc.lt"); got == "" {
		t.Error("published_utc.lt query param is empty")
	}
}
