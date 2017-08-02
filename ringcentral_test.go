package ringcentral

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

var (
	ctx   = context.Background()
	token *Token
	a     *API
)

func init() {
	a = New(os.Getenv("RINGCENTRAL_APP_ID"), os.Getenv("RINGCENTRAL_APP_SECRET"), "~")
	a.TestMode = true
}

func TestMakeURL(t *testing.T) {
	urlStr := a.makeURL("/foo/bar", nil)
	assert.Equal(t, "https://platform.devtest.ringcentral.com/foo/bar", urlStr)
}

// TestAuthorizeUsername should run first to get a valid auth token.
func TestAuthorizeUsername(t *testing.T) {

	var err error
	u := os.Getenv("RINGCENTRAL_PHONE_NUMBER")
	ext := os.Getenv("RINGCENTRAL_EXTENSION")
	pwd := os.Getenv("RINGCENTRAL_PASSWORD")

	token, err = a.Authorize(ctx, u, ext, pwd)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, a.Token)
	assert.NotNil(t, token)
}

func TestGetExtensionList(t *testing.T) {

	if !assert.NotNil(t, a.Token) {
		t.Fail()
	}

	l, err := a.GetExtensionList(ctx, nil)
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotNil(t, l) {
		return
	}
	if !assert.NotNil(t, l.Records) {
		return
	}
	assert.True(t, len(l.Records) > 0)
}
