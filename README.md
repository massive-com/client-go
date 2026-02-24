# Massive (formerly Polygon.io) Go Client - WebSocket & RESTful APIs

The official Go client library for the [Massive](https://massive.com/) REST and WebSocket API. This client makes use of Go generics and thus requires Go 1.18. See the [docs](https://massive.com/docs/stocks/getting-started) for more details on our API.

**Note:** Polygon.io has rebranded as [Massive.com](https://massive.com) on Oct 30, 2025. Existing API keys, accounts, and integrations continue to work exactly as before. The only change in this SDK is that it now defaults to the new API base at `api.massive.com`, while `api.polygon.io` remains supported for an extended period.

For details, see our [rebrand announcement blog post](https://massive.com/blog/polygon-is-now-massive/) or open an issue / contact [support@massive.com](mailto:support@massive.com) if you have questions.

## Getting Started

This section guides you through setting up a simple project with massive.com/client-go.

First, make a new directory for your project and navigate into it:
```bash
mkdir myproject && cd myproject
```

Next, initialize a new module for dependency management. This creates a `go.mod` file to track your dependencies:
```bash
go mod init example
```

Then, create a `main.go` file. For quick start, you can find over 100+ [example code snippets](https://github.com/massive-com/client-go/tree/master/rest/example) that demonstrate connecting to both the REST and WebSocket APIs.

Here's a working example that fetches daily aggregates for AAPL (with full pagination and trace support):

```bash
cat > main.go <<EOF
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/massive-com/client-go/v3/rest"
	"github.com/massive-com/client-go/v3/rest/gen"
)

func main() {
	c := rest.NewWithOptions(os.Getenv("MASSIVE_API_KEY"),
		rest.WithTrace(false),
		rest.WithPagination(true),
	)
	ctx := context.Background()

	resp, err := c.GetStocksAggregatesWithResponse(ctx,
		"AAPL",
		1,
		gen.GetStocksAggregatesParamsTimespan("day"),
		"2025-11-03",
		"2025-11-28",
		nil,   // nil = use API defaults (recommended)
	)
	if err != nil {
		log.Fatal(err)
	}

	// prints status + full error body (e.g. "Unknown API Key")
	if err := rest.CheckResponse(resp); err != nil {
		log.Fatal(err)
	}

	iter := rest.NewIteratorFromResponse(c, resp)
	for iter.Next() {
		item := iter.Item()
		fmt.Printf("%+v\n", item)
	}
	if err := iter.Err(); err != nil {
		log.Fatal(err)
	}
}
EOF
```

Please remember to set your Massive API key, which you can find on the massive.com dashboard, in the environment variable `MASSIVE_API_KEY`. Or, as a less secure option, by hardcoding it in your code. But please note that hardcoding the API key can be risky if your code is shared or exposed. You can configure the environment variable by running:

```
export MASSIVE_API_KEY="<your_api_key>"        <- mac/linux
xset MASSIVE_API_KEY "<your_api_key>"          <- windows
```

Then, run `go mod tidy` to automatically download and install the necessary dependencies. This command ensures your `go.mod` file reflects all dependencies used in your project:
```bash
go mod tidy
```

Finally, execute your application:
```bash
go run main.go
```

## REST API Client

[![rest-docs][rest-doc-img]][rest-doc]

To get started, you'll need to import two main packages.

```golang
import (
	"context"
	"fmt"
	"log"

	"github.com/massive-com/client-go/v3/rest"
	"github.com/massive-com/client-go/v3/rest/gen"
)
```

Next, create a new client with your [API key](https://massive.com/dashboard/signup).

```go
c := rest.NewWithOptions("YOUR_API_KEY",
	rest.WithTrace(false),      // set true for full request/response logging
	rest.WithPagination(true),  // enables automatic pagination via iterator
)
ctx := context.Background()
```

Or create a client with a custom HTTP client implementation.

```golang
hc := http.Client{} // some custom HTTP client
c := massive.NewWithClient("YOUR_API_KEY", hc)
```

### Using the client

After creating the client, making calls to the Massive API is simple. Most endpoints now use the generated `*WithResponse` methods:

```go
// Example with custom params (note the Sort field requirement)
params := &gen.GetStocksAggregatesParams{
	Adjusted: rest.Ptr(true),
	Sort:     "asc",
	Limit:    rest.Ptr(120),
}

resp, err := c.GetStocksAggregatesWithResponse(ctx,
	"AAPL", 1, gen.GetStocksAggregatesParamsTimespan("day"),
	"2025-11-03", "2025-11-28",
	params,
)
```

### Pagination

Our client iterators that handle pagination for you, so when there are multiple pages of results, we'll follow and build the `next_url` page for you and stich the results together.

```go
c := rest.NewWithOptions("YOUR_API_KEY",
	rest.WithTrace(true),
	rest.WithPagination(true), // turn this on
)
```

### Debugging

Debug/trace mode is now enabled at client creation time (much simpler!):

```go
c := rest.NewWithOptions("YOUR_API_KEY",
	rest.WithTrace(true), // turn this on
	rest.WithPagination(true),
)
```

When enabled you will see clean output for every request:

```
Request URL: https://api.massive.com/v2/aggs/ticker/AAPL/range/1/day/2025-11-03/2025-11-28?adjusted=true&limit=120&sort=asc
Request Headers: map[Authorization:[Bearer REDACTED] User-Agent:[massive-go-test]]
Response Headers: map[Content-Type:[application/json] Date:[Mon, 23 Feb 2026 16:03:29 GMT] ...]
```

This is extremely useful when troubleshooting query params, authentication, or rate limits.

## WebSocket Client

[![ws-docs][ws-doc-img]][ws-doc]

Import the WebSocket client and models packages to get started.

```golang
import (
    massivews "github.com/massive-com/client-go/v3/websocket"
    "github.com/massive-com/client-go/v3/websocket/models"
)
```

Next, create a new client with your API key and a couple other config options.

```golang
// create a new client
c, err := massivews.New(massivews.Config{
    APIKey:    "YOUR_API_KEY",
    Feed:      massivews.RealTime,
    Market:    massivews.Stocks,
})
if err != nil {
    log.Fatal(err)
}
defer c.Close() // the user of this client must close it

// connect to the server
if err := c.Connect(); err != nil {
    log.Error(err)
    return
}
```

The client automatically reconnects to the server when the connection is dropped. By default, it will attempt to reconnect indefinitely but the number of retries is configurable. When the client successfully reconnects, it automatically resubscribes to any topics that were set before the disconnect.

### Using the client

After creating a client, subscribe to one or more topics and start accessing data. Currently, all of the data is pushed to a single output channel.

```golang
// passing a topic by itself will subscribe to all tickers
if err := c.Subscribe(massivews.StocksSecAggs); err != nil {
    log.Fatal(err)
}
if err := c.Subscribe(massivews.StocksTrades, "TSLA", "GME"); err != nil {
    log.Fatal(err)
}

for {
    select {
    case err := <-c.Error(): // check for any fatal errors (e.g. auth failed)
        log.Fatal(err)
    case out, more := <-c.Output(): // read the next data message
        if !more {
            return
        }

        switch out.(type) {
        case models.EquityAgg:
            log.Print(out) // do something with the agg
        case models.EquityTrade:
            log.Print(out) // do something with the trade
        }
    }
}
```

See the [full example](./websocket/example/main.go) for more details on how to use this client effectively.

## Contributing

If you found a bug or have an idea for a new feature, please first discuss it with us by [submitting a new issue](https://github.com/massive-com/client-go/issues/new/choose). We will respond to issues within at most 3 weeks. We're also open to volunteers if you want to submit a PR for any open issues but please discuss it with us beforehand. PRs that aren't linked to an existing issue or discussed with us ahead of time will generally be declined.

-------------------------------------------------------------------------------

[doc-img]: https://pkg.go.dev/badge/github.com/massive-com/client-go/v3
[doc]: https://pkg.go.dev/github.com/massive-com/client-go/v3
[rest-doc-img]: https://pkg.go.dev/badge/github.com/massive-com/client-go/v3/rest
[rest-doc]: https://pkg.go.dev/github.com/massive-com/client-go/v3/rest
[ws-doc-img]: https://pkg.go.dev/badge/github.com/massive-com/client-go/v3/websocket
[ws-doc]: https://pkg.go.dev/github.com/massive-com/client-go/v3/websocket
[build-img]: https://github.com/massive-com/client-go/v3/actions/workflows/test.yml/badge.svg
[build]: https://github.com/massive-com/client-go/v3/actions
[report-card-img]: https://goreportcard.com/badge/github.com/massive-com/client-go/v3
[report-card]: https://goreportcard.com/report/github.com/massive-com/client-go/v3
