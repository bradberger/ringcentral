// +build appengine

package ringcentral

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/appengine/urlfetch"
)

func getClient(ctx context.Context, to time.Duration) *http.Client {
	if to <= 0 {
		to = defaultRequestTimeout
	}
	ctx, _ = context.WithTimeout(ctx, to)
	return urlfetch.Client(ctx)
}
