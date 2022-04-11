package reference

import (
	"context"
	"net/http"

	"github.com/polygon-io/client-go/rest/iter"
	"github.com/polygon-io/client-go/rest/models"
)

// ListConditionsIter is an iterator for the ListConditions method.
type ListConditionsIter struct {
	iter.Iter
}

// Condition returns the current result that the iterator points to.
func (it *ListConditionsIter) Condition() models.Condition {
	if it.Item() != nil {
		return it.Item().(models.Condition)
	}
	return models.Condition{}
}

// ListConditions retrieves reference conditions.
// For more details see https://polygon.io/docs/stocks/get_v3_reference_conditions.
// This method returns an iterator that should be used to access the results via this pattern:
//   iter, err := c.ListConditions(context.TODO(), params, opts...)
//   for iter.Next() {
//       // do something with the current value
//       log.Print(iter.Condition())
//   }
//   if iter.Err() != nil {
//       return err
//   }
func (c *Client) ListConditions(ctx context.Context, params *models.ListConditionsParams, options ...models.RequestOption) *ListConditionsIter {
	return &ListConditionsIter{
		Iter: iter.NewIter(ctx, models.ListConditionsPath, params, func(uri string) (iter.ListResponse, []interface{}, error) {
			res := &models.ListConditionsResponse{}
			err := c.CallURL(ctx, http.MethodGet, uri, res, options...)

			results := make([]interface{}, len(res.Results))
			for i, v := range res.Results {
				results[i] = v
			}

			return res, results, err
		}),
	}
}
