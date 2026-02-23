package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

type Iterator struct {
	client  *Client
	page    []map[string]any
	idx     int
	err     error
	nextURL *string
}

func NewIterator(c *Client, firstPage []map[string]any, nextURL *string) *Iterator {
	return &Iterator{
		client:  c,
		page:    firstPage,
		idx:     0,
		nextURL: nextURL,
	}
}

func (it *Iterator) Next() bool {
	if it.err != nil {
		return false
	}

	if it.idx < len(it.page) {
		it.idx++
		return true
	}

	if !it.client.pagination || it.nextURL == nil || *it.nextURL == "" {
		return false
	}

	it.page, it.nextURL, it.err = it.fetchNextPage(*it.nextURL)
	it.idx = 0

	if len(it.page) > 0 && it.err == nil {
		it.idx = 1
		return true
	}
	return false
}

func (it *Iterator) Item() map[string]any { return it.page[it.idx-1] }
func (it *Iterator) Err() error           { return it.err }

func (it *Iterator) fetchNextPage(urlStr string) ([]map[string]any, *string, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}
	if err := it.client.addHeaders(context.Background(), req); err != nil {
		return nil, nil, err
	}

	resp, err := it.client.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("http status %d", resp.StatusCode)
	}

	type paginated struct {
		Results []map[string]any `json:"results"`
		NextURL *string          `json:"next_url,omitempty"`
	}

	var p paginated
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, nil, err
	}
	return p.Results, p.NextURL, nil
}

// Fixed NewIteratorFromResponse — now safely handles BOTH []T and *[]T for Results
func NewIteratorFromResponse(c *Client, resp any) *Iterator {
	if resp == nil {
		return NewIterator(c, nil, nil)
	}

	rv := reflect.ValueOf(resp)
	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return NewIterator(c, nil, nil)
		}
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return NewIterator(c, nil, nil)
	}

	json200Field := rv.FieldByName("JSON200")
	if !json200Field.IsValid() || json200Field.IsNil() {
		return NewIterator(c, nil, nil)
	}
	body := json200Field.Elem()

	// Handle both Results []T and Results *[]T safely
	var page []map[string]any
	if resultsField := body.FieldByName("Results"); resultsField.IsValid() && !resultsField.IsNil() {
		sliceVal := resultsField
		if sliceVal.Kind() == reflect.Pointer {
			if !sliceVal.IsNil() {
				sliceVal = sliceVal.Elem()
			} else {
				sliceVal = reflect.ValueOf([]interface{}{})
			}
		}
		if sliceVal.Kind() == reflect.Slice {
			n := sliceVal.Len()
			page = make([]map[string]any, n)
			for i := 0; i < n; i++ {
				item := sliceVal.Index(i).Interface()
				var m map[string]any
				if b, err := json.Marshal(item); err == nil {
					json.Unmarshal(b, &m)
				}
				page[i] = m
			}
		}
	}

	// Extract next_url (handles both common casing)
	var nextURL *string
	if f := body.FieldByName("NextUrl"); f.IsValid() && f.Kind() == reflect.Pointer && !f.IsNil() {
		nextURL = f.Interface().(*string)
	} else if f := body.FieldByName("NextURL"); f.IsValid() && f.Kind() == reflect.Pointer && !f.IsNil() {
		nextURL = f.Interface().(*string)
	}

	if !c.pagination {
		nextURL = nil
	}

	return NewIterator(c, page, nextURL)
}
