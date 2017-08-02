// +build !appengine

package ringcentral

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

func getClient(ctx context.Context, to time.Duration) *http.Client {
	if to <= 0 {
		to = defaultRequestTimeout
	}
	return &http.Client{Timeout: to}
}
